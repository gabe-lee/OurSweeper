package user_token

import (
	"time"

	"github.com/gabe-lee/OurSweeper/data_buffer"
	"github.com/gabe-lee/OurSweeper/utils"
	"github.com/google/uuid"
)

type (
	Duration     = time.Duration
	Time         = time.Time
	ErrorChecker = utils.ErrorChecker
	UUID         = uuid.UUID
)

const (
	CurrentVer uint16 = 1
)

const (
	EXTRA_NONE byte = iota
	EXTRA_OLD_TOKEN_WAS_CORRUPTED
)

type UserStats struct {
	UUID          UUID
	Playtime      int64
	TotalScore    uint32
	ScoreSweeps   uint32
	ScoreFlags    uint32
	Sweeps        uint32
	TotalFlags    uint32
	GoodFlags     uint32
	Deaths        uint32
	Version       uint16
	ScreenNameLen byte
	ScreenName    [16]byte
}

func (a *UserStats) ReadFromBuffer(buf *data_buffer.ReadBuffer) error {
	e := utils.FirstError{}
	e.Add(buf.U16_LE(&a.Version))
	// Do something different on ver mismatch
	e.Add(buf.U8_Slice(a.UUID[:]))
	e.Add(buf.I64_LE(&a.Playtime))
	e.Add(buf.U32_LE(&a.TotalScore))
	e.Add(buf.U32_LE(&a.ScoreSweeps))
	e.Add(buf.U32_LE(&a.ScoreFlags))
	e.Add(buf.U32_LE(&a.Sweeps))
	e.Add(buf.U32_LE(&a.TotalFlags))
	e.Add(buf.U32_LE(&a.GoodFlags))
	e.Add(buf.U32_LE(&a.Deaths))
	e.Add(buf.U8(&a.ScreenNameLen))
	e.Add(buf.U8_Slice(a.ScreenName[:a.ScreenNameLen]))
	return e.Err
}

func (a *UserStats) WriteToBuffer(buf *data_buffer.WriteBuffer) {
	buf.U16_LE(a.Version)
	// Do something different on ver mismatch
	buf.U8_Slice(a.UUID[:])
	buf.I64_LE(a.Playtime)
	buf.U32_LE(a.TotalScore)
	buf.U32_LE(a.ScoreSweeps)
	buf.U32_LE(a.ScoreFlags)
	buf.U32_LE(a.Sweeps)
	buf.U32_LE(a.TotalFlags)
	buf.U32_LE(a.GoodFlags)
	buf.U32_LE(a.Deaths)
	buf.U8(a.ScreenNameLen)
	buf.U8_Slice(a.ScreenName[:a.ScreenNameLen])
}

func (a *UserStats) SizeOnBuffer() int {
	return 56 + int(a.ScreenNameLen)
}

var _ data_buffer.SizeReadWritable = (*UserStats)(nil)

// type AnonTokenRaw struct {
// 	Token []byte
// }

// // WireRead implements wire.WireReader.
// func (a *AnonTokenRaw) WireRead(wire *wire.Incoming) {
// 	var l uint32
// 	wire.UVar32(&l)
// 	a.Token = make([]byte, 0, l)
// 	wire.U8_Slice(a.Token)
// }

// // WireWrite implements wire.WireWriter.
// func (a *AnonTokenRaw) WireWrite(wire *wire.Outgoing) {
// 	wire.UVar32(uint32(len(a.Token)))
// 	wire.U8_Slice(a.Token)
// }

// var _ wire.WireWriter = (*AnonTokenRaw)(nil)
// var _ wire.WireReader = (*AnonTokenRaw)(nil)

// var encoding = base64.RawURLEncoding

// func RawTokenEncoder(writer io.Writer) io.WriteCloser {
// 	return base64.NewEncoder(encoding, writer)
// }

// func RawTokenDecoder(reader io.Reader) io.Reader {
// 	return base64.NewDecoder(encoding, reader)
// }
