package ws

import (
	"github.com/gabe-lee/OurSweeper/data_buffer"
)

const (
	DefaultServerReadBufferSize  = 4096
	DefaultServerWriteBufferSize = 512
)

// const (
// 	headerHost          = "Host"
// 	headerUpgrade       = "Upgrade"
// 	headerConnection    = "Connection"
// 	headerSecVersion    = "Sec-WebSocket-Version"
// 	headerSecProtocol   = "Sec-WebSocket-Protocol"
// 	headerSecExtensions = "Sec-WebSocket-Extensions"
// 	headerSecKey        = "Sec-WebSocket-Key"
// 	headerSecAccept     = "Sec-WebSocket-Accept"

// 	headerHostCanonical          = headerHost
// 	headerUpgradeCanonical       = headerUpgrade
// 	headerConnectionCanonical    = headerConnection
// 	headerSecVersionCanonical    = "Sec-Websocket-Version"
// 	headerSecProtocolCanonical   = "Sec-Websocket-Protocol"
// 	headerSecExtensionsCanonical = "Sec-Websocket-Extensions"
// 	headerSecKeyCanonical        = "Sec-Websocket-Key"
// 	headerSecAcceptCanonical     = "Sec-Websocket-Accept"
// )

// var (
// 	specHeaderValueUpgrade         = []byte("websocket")
// 	specHeaderValueConnection      = []byte("Upgrade")
// 	specHeaderValueConnectionLower = []byte("upgrade")
// 	specHeaderValueSecVersion      = []byte("13")
// )

// var (
// 	httpVersion1_0    = []byte("HTTP/1.0")
// 	httpVersion1_1    = []byte("HTTP/1.1")
// 	httpVersionPrefix = []byte("HTTP/")
// )

// type httpRequestLine struct {
// 	method, uri  []byte
// 	major, minor int
// }

// type httpResponseLine struct {
// 	major, minor int
// 	status       int
// 	reason       []byte
// }

type (
	WriteBufferPool = data_buffer.WriteBufferPool
	WriteBuffer     = data_buffer.WriteBuffer
)

type WsUpgrader struct {
	WriteBufferPool *WriteBufferPool
	ReadBufferSize  int
	WriteBufferSize int
}

// func (w *WsUpgrader) Upgrade(conn io.ReadWriter) {
// 	ws.DefaultUpgrader.Upgrade(conn)
// 	wBufSize := w.WriteBufferSize
// 	if wBufSize == 0 {
// 		wBufSize = DefaultServerWriteBufferSize
// 	}
// 	rBufSize := w.ReadBufferSize
// 	if rBufSize == 0 {
// 		rBufSize = DefaultServerWriteBufferSize
// 	}
// 	writeBuf := w.WriteBufferPool.GetBuffer(wBufSize)
// 	readBuf := w.WriteBufferPool.GetBuffer(rBufSize)
// 	defer func() {
// 		w.WriteBufferPool.ReleaseBuffer(writeBuf)
// 		w.WriteBufferPool.ReleaseBuffer(readBuf)
// 	}()
// 	_, err := readBuf.ReadFrom(conn)
// }
