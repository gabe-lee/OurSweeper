package active_worlds_response

import (
	"github.com/gabe-lee/OurSweeper/data_buffer"
	C "github.com/gabe-lee/OurSweeper/internal/consts"
	"github.com/gabe-lee/OurSweeper/utils"
)

type ActiveWorldsReport struct {
	Worlds    [C.MAX_ACTIVE_WORLDS]ActiveWorldStats
	WorldsLen byte
}

// SizeOnBuffer implements data_buffer.SizeReadWritable.
func (a *ActiveWorldsReport) SizeOnBuffer() int {
	total := 1
	for i := range a.WorldsLen {
		total += a.Worlds[i].SizeOnBuffer()
	}
	return total
}

// ReadFromBuffer implements data_buffer.SizeReadWritable.
func (a *ActiveWorldsReport) ReadFromBuffer(buf *data_buffer.ReadBuffer) error {
	e := utils.FirstError{}
	e.Add(buf.U8(&a.WorldsLen))
	for i := range a.WorldsLen {
		e.Add(buf.Readable(&a.Worlds[i]))
	}
	return e.Err
}

// WriteToBuffer implements data_buffer.SizeReadWritable.
func (a *ActiveWorldsReport) WriteToBuffer(buf *data_buffer.WriteBuffer) {
	buf.U8(a.WorldsLen)
	for i := range a.WorldsLen {
		buf.Writable(&a.Worlds[i])
	}
}

var _ data_buffer.SizeReadWritable = (*ActiveWorldsReport)(nil)

type ActiveWorldStats struct {
	ID              uint64
	Expires         int64
	TotalMines      uint32
	RemainingMines  uint32
	RemainingSpaces uint32
	CurrentUsers    uint16
	Difficulty      uint8
}

// SizeOnBuffer implements data_buffer.SizeReadWritable.
func (a *ActiveWorldStats) SizeOnBuffer() int {
	return 31
}

// ReadFromBuffer implements data_buffer.SizeReadWritable.
func (a *ActiveWorldStats) ReadFromBuffer(buf *data_buffer.ReadBuffer) error {
	e := utils.FirstError{}
	e.Add(buf.U64_LE(&a.ID))
	e.Add(buf.I64_LE(&a.Expires))
	e.Add(buf.U32_LE(&a.TotalMines))
	e.Add(buf.U32_LE(&a.RemainingMines))
	e.Add(buf.U32_LE(&a.RemainingSpaces))
	e.Add(buf.U16_LE(&a.CurrentUsers))
	e.Add(buf.U8(&a.Difficulty))
	return e.Err
}

// WriteToBuffer implements data_buffer.SizeReadWritable.
func (a *ActiveWorldStats) WriteToBuffer(buf *data_buffer.WriteBuffer) {
	buf.U64_LE(a.ID)
	buf.I64_LE(a.Expires)
	buf.U32_LE(a.TotalMines)
	buf.U32_LE(a.RemainingMines)
	buf.U32_LE(a.RemainingSpaces)
	buf.U16_LE(a.CurrentUsers)
	buf.U8(a.Difficulty)
}

var _ data_buffer.SizeReadWritable = (*ActiveWorldStats)(nil)
