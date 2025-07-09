package sweep_request

// import (
// 	"encoding/binary"
// 	"io"

// 	msg "github.com/gabe-lee/OurSweeper/internal/messages"
// 	"github.com/gabe-lee/OurSweeper/internal/utils"
// )

// type SweepRequest struct {
// 	Index uint16
// }

// func NewSweepRequest(x int, y int) SweepRequest {
// 	return SweepRequest{
// 		Index: uint16(world.GetIndex(x, y)),
// 	}
// }

// func (s *SweepRequest) WriteWire(w io.Writer, order binary.ByteOrder) error {
// 	var e utils.ErrorChecker
// 	if e.IsErr(binary.Write(w, order, msg.CLIENT_SWEEP)) {
// 		return e.Err
// 	}
// 	if e.IsErr(binary.Write(w, order, s.Index)) {
// 		return e.Err
// 	}
// 	return e.Err
// }
// func (s *SweepRequest) ReadWire(r io.Reader, order binary.ByteOrder) error {
// 	var e utils.ErrorChecker
// 	if e.IsErr(binary.Read(r, order, s.Index)) {
// 		return e.Err
// 	}
// 	return e.Err
// }
