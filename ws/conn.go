package ws

import (
	"net"
)

type WsConn struct {
	conn   net.Conn
	server bool
}

func ServerConn(conn net.Conn) WsConn {
	return WsConn{
		conn:   conn,
		server: true,
	}
}

func ClientConn(conn net.Conn) WsConn {
	return WsConn{
		conn:   conn,
		server: true,
	}
}

type OpCode byte

var (
	OpContinue    = OpCode(0)
	OpTextStart   = OpCode(1)
	OpBinaryStart = OpCode(2)
	OpClose       = OpCode(8)
	OpPing        = OpCode(9)
	OpPong        = OpCode(10)
)

type CloseStatus uint16

var (
	CloseNormal           = CloseStatus(1000)
	CloseGoingAway        = CloseStatus(1001)
	CloseError            = CloseStatus(1002)
	CloseUnsupported      = CloseStatus(1003)
	CloseNoStatusProvided = CloseStatus(1005)
	CloseAbnormal         = CloseStatus(1006)
	CloseBadContent       = CloseStatus(1007)
	ClosePolicyViolated   = CloseStatus(1008)
	CloseTooLarge         = CloseStatus(1009)
	CloseExtensionNeeded  = CloseStatus(1010)
	CloseUnexpectedCond   = CloseStatus(1011)
	CloseTLSFail          = CloseStatus(1015)
)

type WsHeader struct {
	Code   OpCode
	Len    uint64
	Masked bool
	Final  bool
	Ext    byte
	Mask   [4]byte
}
