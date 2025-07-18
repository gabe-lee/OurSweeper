package anon_token

import (
	"time"

	"github.com/gabe-lee/OurSweeper/utils"
	"github.com/gabe-lee/OurSweeper/wire"
)

type (
	Duration     = time.Duration
	Time         = time.Time
	ErrorChecker = utils.ErrorChecker
)

type AnonToken struct {
	UUID_A        uint64
	UUID_B        uint64
	Playtime      int64
	Expires       int64
	TotalScore    uint32
	ScoreSweeps   uint32
	ScoreFlags    uint32
	Sweeps        uint32
	TotalFlags    uint32
	GoodFlags     uint32
	Deaths        uint32
	ScreenNameLen byte
	ScreenName    [16]byte
}

// WireRead implements wire.WireReader.
func (a *AnonToken) WireRead(wire *wire.IncomingWire) {
	wire.TryRead_U64(&a.UUID_A)
	wire.TryRead_U64(&a.UUID_B)
	wire.TryRead_I64(&a.Playtime)
	wire.TryRead_I64(&a.Expires)
	wire.TryRead_U32(&a.TotalScore)
	wire.TryRead_U32(&a.ScoreSweeps)
	wire.TryRead_U32(&a.ScoreFlags)
	wire.TryRead_U32(&a.Sweeps)
	wire.TryRead_U32(&a.TotalFlags)
	wire.TryRead_U32(&a.GoodFlags)
	wire.TryRead_U32(&a.Deaths)
	wire.TryRead_U8(&a.ScreenNameLen)
	dst := a.ScreenName[:a.ScreenNameLen]
	wire.TryRead_SliceU8(dst)
}

// WireWrite implements wire.WireWriter.
func (a *AnonToken) WireWrite(wire *wire.OutgoingWire) {
	wire.TryWrite_U64(a.UUID_A)
	wire.TryWrite_U64(a.UUID_B)
	wire.TryWrite_I64(a.Playtime)
	wire.TryWrite_I64(a.Expires)
	wire.TryWrite_U32(a.TotalScore)
	wire.TryWrite_U32(a.ScoreSweeps)
	wire.TryWrite_U32(a.ScoreFlags)
	wire.TryWrite_U32(a.Sweeps)
	wire.TryWrite_U32(a.TotalFlags)
	wire.TryWrite_U32(a.GoodFlags)
	wire.TryWrite_U32(a.Deaths)
	wire.TryWrite_U8(a.ScreenNameLen)
	src := a.ScreenName[:a.ScreenNameLen]
	wire.TryWrite_SliceU8(src)
}

var _ wire.WireWriter = (*AnonToken)(nil)
var _ wire.WireReader = (*AnonToken)(nil)

type AnonTokenRaw struct {
	Token []byte
}

// WireRead implements wire.WireReader.
func (a *AnonTokenRaw) WireRead(wire *wire.IncomingWire) {
	var l uint32
	wire.TryRead_UVar32(&l)
	a.Token = make([]byte, 0, l)
	wire.TryRead_SliceU8(a.Token)
}

// WireWrite implements wire.WireWriter.
func (a *AnonTokenRaw) WireWrite(wire *wire.OutgoingWire) {
	wire.TryWrite_UVar32(uint32(len(a.Token)))
	wire.TryWrite_SliceU8(a.Token)
}

var _ wire.WireWriter = (*AnonTokenRaw)(nil)
var _ wire.WireReader = (*AnonTokenRaw)(nil)
