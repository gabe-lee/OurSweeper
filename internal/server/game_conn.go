package server

import (
	"math"
	"net"
	"time"

	"github.com/gabe-lee/OurSweeper/data_buffer"
	"github.com/gabe-lee/OurSweeper/internal/user_token"
	MSG "github.com/gabe-lee/OurSweeper/internal/wire_codes"
	"github.com/gabe-lee/OurSweeper/token"
	"github.com/gobwas/ws"
	"github.com/google/uuid"
)

const (
	CurrentVersion    uint32 = 1
	BadVersionMsg     string = "Unsupported version %d, current version is %d, please update your client to play"
	ConnTimeout              = time.Second * 60
	ConnIdleTime             = time.Second * 60
	ConnProbeInterval        = time.Second * 15
	ConnProbeMax             = 10
	ConnMaxMsg               = 1024
	MaxMsgLen                = math.MaxUint16
)

func (n *WebNetwork) StartGameConn(conn net.Conn) {
	refreshDeadline(conn)
	var closeFrame []byte = ws.CompiledCloseNormalClosure
	var err error
	var inHeader ws.Header
	var msgcode uint32
	var writer *data_buffer.WriteBuffer = n.utils.WriteBufferPool.GetBufferByClass(data_buffer.Class_128)
	var reader data_buffer.ReadBuffer
	defer n.CloseConn(conn, &closeFrame, writer)
	for {
		select {
		case <-n.shutdownSig:
			closeFrame = ws.CompiledCloseGoingAway
			return
		default:
		}
		inHeader, err = ws.ReadHeader(conn)
		refreshDeadline(conn)
		if n.log.WarnIfErr(err, "failed to read websocket header") != 0 {
			closeFrame = ws.CompiledCloseProtocolError
			return
		}
		if inHeader.OpCode == ws.OpText {
			n.log.Warn("recieved text websocket frame: game server messages only accept binary frames")
			closeFrame = ws.CompiledCloseUnsupportedData
			return
		}
		if !inHeader.Fin || inHeader.OpCode == ws.OpContinuation {
			n.log.Warn("recieved non-final websocket frame: game server only accepts single frame messages")
			closeFrame = ws.CompiledCloseProtocolError
			return
		}
		if inHeader.Rsv != 0 {
			n.log.Warn("recieved websocket extension bits: game server does not use ws extensions")
			closeFrame = ws.CompiledCloseProtocolError
			return
		}
		if inHeader.OpCode == ws.OpClose {
			closeFrame = ws.CompiledCloseNormalClosure
			return
		}
		if inHeader.OpCode == ws.OpPing {
			_, err = conn.Write(ws.CompiledPong)
			if n.HasWriteErr(err, conn, &closeFrame) {
				return
			}
			continue
		}
		if inHeader.OpCode == ws.OpPong {
			refreshDeadline(conn)
			continue
		}
		writer.Reset()
		_, err = writer.ReadNFrom(conn, int(inHeader.Length))
		if n.HasWriteErr(err, conn, &closeFrame) {
			return
		}
		ws.Cipher(writer.BytesRef(), inHeader.Mask, 0)
		reader = writer.ReaderRef()
		err = reader.U32_LE(&msgcode)
		if n.HasReadErr(err, conn, &closeFrame) {
			return
		}
		switch msgcode {
		case MSG.CLIENT_ANON_TOKEN_NEW:
			writer.Reset()
			writer.U32_LE(MSG.SERVER_ANON_TOKEN_NEW)
			id, err := uuid.NewRandom()
			if err != nil {
				if n.ReturnErrorMessage(conn, writer, MSG.ERR_S_INTERNAL_SERVER_ERROR, &closeFrame, "failed to create UUID for AnonToken") {
					continue
				}
				return
			}
			userStats := user_token.UserStats{
				UUID:    id,
				Version: user_token.CurrentVer,
			}
			token.Create(n.env.tokenSecret, &userStats, writer)
			n.FinishWritingMessage(conn, writer, &closeFrame)
		case MSG.CLIENT_ANON_TOKEN_LOGIN:
			//TODO
			//CHECKPOINT
		case MSG.CLIENT_GET_ACTIVE_WORLDS:
			writer.Reset()
			writer.U32_LE(MSG.SERVER_SEND_ACTIVE_WORLDS)
			response := n.worlds.GetActiveWorldsResponse()
			writer.Writable(&response)
			n.FinishWritingMessage(conn, writer, &closeFrame)
		default:
			writer.U32_LE(MSG.SERVER_INVALID)
			n.FinishWritingMessage(conn, writer, &closeFrame)
		}
	}
}

func (n *WebNetwork) CloseConn(conn net.Conn, closeFrame *[]byte, buffer *data_buffer.WriteBuffer) {
	_, err := conn.Write(*closeFrame)
	n.log.WarnIfErr(err, "failed to send websocket close frame")
	conn.Close()
	n.utils.WriteBufferPool.ReleaseBuffer(buffer)
	n.websocketWait.Done()
}

func (n *WebNetwork) HasReadErr(err error, conn net.Conn, closeFrame *[]byte) bool {
	if n.log.WarnIfErr(err, "network read error") != 0 {
		*closeFrame = ws.CompiledCloseInvalidFramePayloadData
		return true
	}
	return false
}

func (n *WebNetwork) HasWriteErr(err error, conn net.Conn, closeFrame *[]byte) bool {
	if n.log.WarnIfErr(err, "network write error") != 0 {
		*closeFrame = ws.CompiledCloseInternalServerError
		return true
	}
	return false
}

func refreshDeadline(conn net.Conn) {
	conn.SetDeadline(time.Now().Add(ConnTimeout))
}

func (n *WebNetwork) FinishWritingMessage(conn net.Conn, buf *data_buffer.WriteBuffer, closeFrame *[]byte) (success bool) {
	outHeader := ws.Header{
		OpCode: ws.OpBinary,
		Fin:    true,
		Length: int64(buf.Len()),
	}
	err := ws.WriteHeader(conn, outHeader)
	if n.HasWriteErr(err, conn, closeFrame) {
		return false
	}
	_, err = conn.Write(buf.BytesRef())
	return !n.HasWriteErr(err, conn, closeFrame)
}

var ErrHeader = ws.Header{
	OpCode: ws.OpBinary,
	Fin:    true,
	Length: 16,
}
var ErrHeaderPayload = make([]byte, 16)
var CompiledErrHeader = ws.MustCompileFrame(ws.Frame{Header: ErrHeader, Payload: ErrHeaderPayload})

func (n *WebNetwork) ReturnErrorMessage(conn net.Conn, buf *data_buffer.WriteBuffer, errReasonCode uint32, closeFrame *[]byte, logFormat string, logArgs ...any) (success bool) {
	buf.Reset()
	buf.Write(CompiledErrHeader)
	buf.SetLenRelative(-16)
	buf.U32_LE(MSG.SERVER_ERROR)
	buf.U32_LE(errReasonCode)
	logid := n.log.Error(logFormat, logArgs...)
	buf.U64_LE(logid)
	_, err := conn.Write(buf.BytesRef())
	if n.HasWriteErr(err, conn, closeFrame) {
		return false
	}
	return true
}
