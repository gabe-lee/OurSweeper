package internal

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/gabe-lee/OurSweeper/serializer"
	"github.com/gabe-lee/OurSweeper/utils"
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

func (s *ServerErrorResult) Deserialize(r io.Reader, order binary.ByteOrder) error {
	e := utils.ErrorCollector{}
	e.Do(binary.Read(r, order, s.Code))
	e.Do(binary.Read(r, order, s.Extra))
	return e.Err
}

func (s *ServerErrorResult) CodedSerialize(order binary.ByteOrder) (data []byte, err error) {
	e := utils.ErrorCollector{}
	data = make([]byte, 0, 8)
	w := bytes.NewBuffer(data)
	e.Do(binary.Write(w, order, SERVER_ERROR))
	e.Do(binary.Write(w, order, s.Code))
	e.Do(binary.Write(w, order, s.Extra))
	return data, e.Err
}

var _ serializer.CodedSerializer = (*ServerErrorResult)(nil)
var _ serializer.WireReader = (*ServerErrorResult)(nil)

var SERVER_MSG = [msgCountS]string{
	S_ERR_UNKNOWN_MSG_CODE_FROM_CLIENT: "Unknown ",
	S_ERR_INVALID_USERNAME_PASS:        "Invalid username or password",
	S_ERR_TILE_ALREADY_SWEPT:           "The tile is already swept",
}
