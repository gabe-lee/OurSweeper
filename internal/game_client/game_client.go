package game_client

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/gabe-lee/OurSweeper/coord"
	"github.com/gabe-lee/OurSweeper/data_buffer"
	"github.com/gabe-lee/OurSweeper/internal/active_worlds_response"
	"github.com/gabe-lee/OurSweeper/internal/common"
	C "github.com/gabe-lee/OurSweeper/internal/consts"
	"github.com/gabe-lee/OurSweeper/internal/user_token"
	MSG "github.com/gabe-lee/OurSweeper/internal/wire_codes"
	"github.com/gabe-lee/OurSweeper/logger"
	"github.com/gabe-lee/OurSweeper/scrap/wire"
	"github.com/gabe-lee/OurSweeper/token"
	"github.com/gabe-lee/OurSweeper/vec2"
	"github.com/gabe-lee/OurSweeper/xmath"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	EbitImage          = ebiten.Image
	ClientWorld        = common.ClientWorld
	ServerWorld        = common.ServerWorld
	SweepResult        = common.SweepResult
	Coord              = coord.Coord[int]
	ByteCoord          = coord.Coord[byte]
	WriteBuffer        = data_buffer.WriteBuffer
	UserStats          = user_token.UserStats
	ActiveWorldsReport = active_worlds_response.ActiveWorldsReport
	Vec2_F64           = vec2.Vec2[float64]
	Vec2_Int           = vec2.Vec2[int]
)

//go:embed tiles.png
var tilesPng []byte

type GameClient struct {
	World              ClientWorld
	Atlas              *ebiten.Image
	BoardX             float64
	BoardY             float64
	Input              Input
	Score              uint32
	Frame              uint64
	RecieveMessages    <-chan *WriteBuffer
	SendMessages       chan<- *WriteBuffer
	Log                logger.Logger
	UserStats          UserStats
	AnonToken          []byte
	ActiveServerWorlds ActiveWorldsReport
}

type Input struct {
	ScrollX           float64
	ScrollY           float64
	MouseX            int
	MouseY            int
	MouseLJustPressed bool
	MouseRJustPressed bool
	MouseLDown        bool
}

// Draw implements ebiten.Game.
func (g *GameClient) Draw(screen *EbitImage) {
	for i := range C.WORLD_TILE_COUNT {
		tilePos := coord.CoordFromIndex(i, C.TY_SHIFT, C.TX_MASK)
		boardPos := tilePos.MultScalar(C.TILE_SIZE).DivScalar(C.DISPLAY_SCALE_DOWN)
		iconIdx := g.World.Tiles[i]
		iconTopLeft := C.BOARD_TILES[iconIdx]
		iconBotRight := [2]int{iconTopLeft[0] + C.TILE_SIZE, iconTopLeft[1] + C.TILE_SIZE}
		rect := image.Rect(iconTopLeft[0], iconTopLeft[1], iconBotRight[0], iconBotRight[1])
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(float64(boardPos.X), float64(boardPos.Y))
		op.GeoM.Translate(g.BoardX, g.BoardY)
		screen.DrawImage(g.Atlas.SubImage(rect).(*EbitImage), op)
	}
}

// Layout implements ebiten.Game.
func (g *GameClient) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameClient) Init(clientToServer chan<- []byte, serverToClient <-chan []byte) {
	img, _, err := image.Decode(bytes.NewReader(tilesPng))
	if err != nil {
		log.Fatal(err)
	}
	g.Atlas = ebiten.NewImageFromImage(img)
}

// Update implements ebiten.Game.
func (g *GameClient) Update() error {
	g.Frame += 1
	{ // Poll input
		wx, wy := ebiten.Wheel()
		g.Input.ScrollX = wx
		g.Input.ScrollY = wy
		cx, cy := ebiten.CursorPosition()
		g.Input.MouseX = cx
		g.Input.MouseY = cy
		g.Input.MouseLJustPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		g.Input.MouseRJustPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	}
	{ // Get Messages
		more_messages := true
		for more_messages {
			select {
			case msg := <-g.RecieveMessages:
				defer msg.Close()
				rdr := msg.ReaderRef()
				var msgcode uint32
				err := rdr.U32_LE(&msgcode)
				if g.Log.ErrorIfErr(err, "failed to read incoming message code") != 0 {
					continue
				}
				switch msgcode {
				case MSG.SERVER_SWEEP:
					var sweep SweepResult
					err := rdr.Readable(&sweep)
					if g.Log.ErrorIfErr(err, "failed to read incoming sweep result") != 0 {
						continue
					}
					g.Score += uint32(sweep.Score)
					sweep.DoActionOnAllTiles(func(pos Coord, icon byte) {
						idx := pos.ToIndex(C.TY_SHIFT)
						g.World.Tiles[idx] = icon
					})
				case MSG.SERVER_ANON_TOKEN_NEW:
					g.AnonToken = g.AnonToken[:0]
					g.AnonToken = append(g.AnonToken, rdr.UnreadBytesRef()...)
					var userStats user_token.UserStats
					_, _, err = token.Open(&rdr, &userStats)
					if g.Log.ErrorIfErr(err, "failed to open incoming anon token") != 0 {
						continue
					}
					g.UserStats = userStats
				case MSG.SERVER_SEND_ACTIVE_WORLDS:
					err := rdr.Readable(&g.ActiveServerWorlds)
					if g.Log.ErrorIfErr(err, "failed to read active worlds response") != 0 {
						continue
					}
				default:
					g.Log.Warn("invalid msg code: %d", msgcode)
				}
			default:
				more_messages = false
			}
		}
	}
	{ // Update State and Send Messages
		g.BoardX += g.Input.ScrollX * C.WHEEL_SPEED
		g.BoardY += g.Input.ScrollY * C.WHEEL_SPEED
		g.BoardX = xmath.Clamp(C.MIN_BOARD_POS_X, g.BoardX, C.MAX_BOARD_POS_X)
		g.BoardY = xmath.Clamp(C.MIN_BOARD_POS_Y, g.BoardY, C.MAX_BOARD_POS_Y)
		if g.Input.MouseLJustPressed {
			tilePos := g.MousePosToTilePos()
			tileIdx := tilePos.ToIndex(C.TY_SHIFT)
			if g.World.Tiles[tileIdx] > 8 { //FIXME make `ClientTile` type with readable methods (cheking whether tile is not swept here)

				request := common.NewSweepRequest(tilePos)
				buf := bytes.Buffer{}
				buf.Grow(64)
				outWire := wire.NewOutgoing(&buf, wire.LE) //FIXME
				request.WireWrite(&outWire)                //FIXME
				if outWire.HasErr() {
					g.Log.Warn("failed to write SweepRequest message: %s", outWire.Err())
				} else {
					g.SendMessages <- buf.Bytes()
				}
			}
		}
	}

	return nil
}

func (g *GameClient) MousePosToBoardPos() (x, y float64) {
	x, y = float64(g.Input.MouseX), float64(g.Input.MouseY)
	x += -g.BoardX
	y += -g.BoardY
	return x, y
}

func (g *GameClient) MousePosToTilePos() (pos Coord) {
	fx, fy := g.MousePosToBoardPos()
	fx /= float64(C.TILE_SIZE_SCALED)
	fy /= float64(C.TILE_SIZE_SCALED)
	return Coord{
		X: int(fx),
		Y: int(fy),
	}
}

var _ ebiten.Game = (*GameClient)(nil)
