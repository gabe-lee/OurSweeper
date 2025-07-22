package user_token

import (
	"encoding/base64"
	"io"
	"time"

	"github.com/gabe-lee/OurSweeper/utils"
	"github.com/gabe-lee/OurSweeper/wire"
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

// WireRead implements wire.WireReader.
func (a *UserStats) WireRead(wire *wire.Incoming) {
	wire.U16(&a.Version)
	// Do something different on ver mismatch
	wire.U8_Slice(a.UUID[:])
	wire.I64(&a.Playtime)
	wire.U32(&a.TotalScore)
	wire.U32(&a.ScoreSweeps)
	wire.U32(&a.ScoreFlags)
	wire.U32(&a.Sweeps)
	wire.U32(&a.TotalFlags)
	wire.U32(&a.GoodFlags)
	wire.U32(&a.Deaths)
	wire.U8(&a.ScreenNameLen)
	dst := a.ScreenName[:a.ScreenNameLen]
	wire.U8_Slice(dst)
}

// WireWrite implements wire.WireWriter.
func (a *UserStats) WireWrite(wire *wire.Outgoing) {
	wire.U16(a.Version)
	// Do something different on ver mismatch
	wire.U8_Slice(a.UUID[:])
	wire.I64(a.Playtime)
	wire.U32(a.TotalScore)
	wire.U32(a.ScoreSweeps)
	wire.U32(a.ScoreFlags)
	wire.U32(a.Sweeps)
	wire.U32(a.TotalFlags)
	wire.U32(a.GoodFlags)
	wire.U32(a.Deaths)
	wire.U8(a.ScreenNameLen)
	src := a.ScreenName[:a.ScreenNameLen]
	wire.U8_Slice(src)
}

var _ wire.WireWriter = (*UserStats)(nil)
var _ wire.WireReader = (*UserStats)(nil)

type AnonTokenRaw struct {
	Token []byte
}

// WireRead implements wire.WireReader.
func (a *AnonTokenRaw) WireRead(wire *wire.Incoming) {
	var l uint32
	wire.UVar32(&l)
	a.Token = make([]byte, 0, l)
	wire.U8_Slice(a.Token)
}

// WireWrite implements wire.WireWriter.
func (a *AnonTokenRaw) WireWrite(wire *wire.Outgoing) {
	wire.UVar32(uint32(len(a.Token)))
	wire.U8_Slice(a.Token)
}

var _ wire.WireWriter = (*AnonTokenRaw)(nil)
var _ wire.WireReader = (*AnonTokenRaw)(nil)

var encoding = base64.RawURLEncoding

func RawTokenEncoder(writer io.Writer) io.WriteCloser {
	return base64.NewEncoder(encoding, writer)
}

func RawTokenDecoder(reader io.Reader) io.Reader {
	return base64.NewDecoder(encoding, reader)
}
