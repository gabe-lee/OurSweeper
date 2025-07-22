package slim_protocol

import (
	"github.com/gabe-lee/OurSweeper/wire"
)

const (
	// Data Ops

	opDataCont   byte = 0
	opDataFinal  byte = 1
	invalidStart byte = 2

	// Control Ops
	invalidEnd byte = 251
	opPing     byte = 252
	opPong     byte = 253
	opClose    byte = 254
	opInvalid  byte = 255
)

const (
	maxLen = 0b11111111_11111111_11111111
)

type OpCode byte

var (
	OpDataCont  = OpCode(opDataCont)
	OpDataFinal = OpCode(opDataFinal)

	OpPing    = OpCode(opPing)
	OpPong    = OpCode(opPong)
	OpClose   = OpCode(opClose)
	OpInvalid = OpCode(opInvalid)
)

type Header struct {
	op  byte
	len uint32
}

// Set the OpCode for the message
//
// Invalid codes will be transformed to `OpInvalid`
func (h *Header) SetOp(op OpCode) {
	h.op = byte(op)
	if h.op >= invalidStart && h.op <= invalidEnd {
		h.op = opInvalid
	}
}

func (h *Header) GetOp() OpCode {
	return OpCode(h.op)
}

func (h *Header) GetLen() uint32 {
	return h.len
}

func (h *Header) SetLen(val uint32) {
	if val > maxLen {
		panic("len > 16777215")
	}
	h.len = val
}

// WireWrite implements wire.WireWriter.
func (h *Header) WireWrite(write *wire.Outgoing) {
	write.U8(h.op)
	write.U24(h.len)
}

// WireRead implements wire.WireReader.
func (h *Header) WireRead(read *wire.Incoming) {
	var o byte
	read.U8(&o)
	h.SetOp(OpCode(o))
	read.U24(&h.len)
}

var _ wire.WireReader = (*Header)(nil)
var _ wire.WireWriter = (*Header)(nil)
