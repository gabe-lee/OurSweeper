package server

import (
	"math"
	"net"
	"time"

	"github.com/gabe-lee/OurSweeper/data_buffer"
	"github.com/gabe-lee/OurSweeper/internal/user_token"
	C "github.com/gabe-lee/OurSweeper/internal/wire_codes"
	"github.com/gabe-lee/OurSweeper/token"
	"github.com/gabe-lee/OurSweeper/wire"
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
	var readWire wire.Incoming
	var writeWire wire.Outgoing
	var err error
	var inHeader ws.Header
	var msgcode uint32
	var buffer *data_buffer.WriteBuffer = n.utils.WriteBufferPool.GetBufferByClass(data_buffer.Class_128)
	defer n.CloseConn(conn, &closeFrame, buffer)
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
		buffer.Reset()
		readWire = wire.NewIncomingAdv(conn, wire.LE, nil, inHeader.Mask[:], int(inHeader.Length))
		writeWire = wire.NewOutgoingAdv(buffer, wire.LE, nil, nil, math.MaxInt)
		readWire.U32(&msgcode)
		if n.HasReadErr(readWire.Err, conn, &closeFrame) {
			return
		}
		switch msgcode {
		case C.CLIENT_ANON_TOKEN_NEW:
			id, err := uuid.NewRandom()
			if err != nil {
				writeWire.U32(C.SERVER_ERROR)
				n.finishWritingMessage(conn, buffer, &closeFrame)
				continue
			}
			userToken := user_token.UserStats{
				UUID:    id,
				Version: user_token.CurrentVer,
			}
			err = token.Create(n.env.tokenSecret, &userToken, buffer)
			if err != nil {
				writeWire.U32(C.SERVER_ERROR)
				n.finishWritingMessage(conn, buffer, &closeFrame)
				continue
			}
			writeWire.U32(C.SERVER_ANON_TOKEN_NEW)
			n.finishWritingMessage(conn, buffer, &closeFrame)
			continue
		case C.CLIENT_GET_ACTIVE_WORLDS:
			//TODO
			//CHECKPOINT
		default:
			writeWire.U32(C.SERVER_INVALID)
			n.finishWritingMessage(conn, buffer, &closeFrame)
			continue
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

func (n *WebNetwork) finishWritingMessage(conn net.Conn, buf *data_buffer.WriteBuffer, closeFrame *[]byte) (success bool) {
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
	if n.HasWriteErr(err, conn, closeFrame) {
		return false
	}
	return true
}
