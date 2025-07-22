package slim_protocol

import "net"

type (
	Conn = net.Conn
)

type SlimConn struct {
	Conn
}
