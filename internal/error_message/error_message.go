package error_message

import (
	"github.com/gabe-lee/OurSweeper/wire"
)

const (
	S_ERR_UNKNOWN_MSG_CODE_FROM_CLIENT uint32 = iota
	S_ERR_INVALID_USERNAME_PASS
	S_ERR_TILE_ALREADY_SWEPT
	msgCountS
)

const (
	C_ERR_UNKNOWN_MSG_CODE_FROM_SERVER uint32 = iota
	C_ERR_INVALID_USERNAME_PASS
	C_ERR_TILE_ALREADY_SWEPT
)

type ServerErrorResult struct {
	Code  uint32
	Extra uint32
}

func (s *ServerErrorResult) WireRead(w *wire.IncomingWire) {
	w.TryRead_U32(&s.Code)
	w.TryRead_U32(&s.Extra)
}

func (s *ServerErrorResult) WireWrite(w *wire.OutgoingWire) {
	w.TryWrite_U32(s.Code)
	w.TryWrite_U32(s.Extra)
}

var _ wire.WireReader = (*ServerErrorResult)(nil)
var _ wire.WireWriter = (*ServerErrorResult)(nil)

var SERVER_MSG = [msgCountS]string{
	S_ERR_UNKNOWN_MSG_CODE_FROM_CLIENT: "Unknown ",
	S_ERR_INVALID_USERNAME_PASS:        "Invalid username or password",
	S_ERR_TILE_ALREADY_SWEPT:           "The tile is already swept",
}
