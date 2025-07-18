package common

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/gabe-lee/OurSweeper/coord"
	"github.com/gabe-lee/OurSweeper/internal/common"
	MSG "github.com/gabe-lee/OurSweeper/internal/wire_codes"
	"github.com/gabe-lee/OurSweeper/logger"
	"github.com/gabe-lee/OurSweeper/wire"
	"github.com/gabe-lee/OurSweeper/xmath"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	EbitImage   = ebiten.Image
	ClientWorld = common.ClientWorld
	ServerWorld = common.ServerWorld
	SweepResult = common.SweepResult
	Coord       = coord.Coord[int]
	ByteCoord   = coord.Coord[byte]
)

const (
	CLIENT_RECIEVE_MESSAGE_BUFFER_LEN int = 32
	CLIENT_SEND_MESSAGE_BUFFER_LEN    int = 4

	TY_SHIFT          = common.TY_SHIFT
	TX_MASK           = common.TX_MASK
	WORLD_TILE_COUNT  = common.WORLD_TILE_COUNT
	WORLD_TILE_WIDTH  = common.WORLD_TILE_WIDTH
	WORLD_TILE_HEIGHT = common.WORLD_TILE_HEIGHT

	TILE_SIZE         int = 32
	TILE_SIZE_SCALED  int = TILE_SIZE / DISPLAY_SCALE_DOWN
	TILE_SHEET_WIDTH  int = 12
	TILE_SHEET_HEIGHT int = 4

	WINDOW_WIDTH         int     = 800
	WINDOW_HEIGHT        int     = 800
	BOARD_WIDTH          int     = TILE_SIZE_SCALED * WORLD_TILE_WIDTH
	BOARD_HEIGHT         int     = TILE_SIZE_SCALED * WORLD_TILE_HEIGHT
	BOARD_OVERFLOW_X     int     = BOARD_WIDTH - WINDOW_WIDTH
	BOARD_OVERFLOW_Y     int     = BOARD_HEIGHT - WINDOW_HEIGHT
	MIN_BOARD_POS_X      float64 = float64(-BOARD_OVERFLOW_X)
	MIN_BOARD_POS_Y      float64 = float64(-BOARD_OVERFLOW_Y)
	MAX_BOARD_POS_X      float64 = 0
	MAX_BOARD_POS_Y      float64 = 0
	DISPLAY_SCALE_DOWN   int     = 2
	DISPLAY_SCALE_DOWN_F float64 = float64(DISPLAY_SCALE_DOWN)
	WHEEL_SPEED          float64 = 6.0
)

var (
	BOARD_TILES = [16][2]int{
		common.ICON_CODE_0:      {0, 0},
		common.ICON_CODE_1:      {1 * TILE_SIZE, 0},
		common.ICON_CODE_2:      {2 * TILE_SIZE, 0},
		common.ICON_CODE_3:      {3 * TILE_SIZE, 0},
		common.ICON_CODE_4:      {4 * TILE_SIZE, 0},
		common.ICON_CODE_5:      {5 * TILE_SIZE, 0},
		common.ICON_CODE_6:      {6 * TILE_SIZE, 0},
		common.ICON_CODE_7:      {7 * TILE_SIZE, 0},
		common.ICON_CODE_8:      {8 * TILE_SIZE, 0},
		common.ICON_CODE_FLAG:   {9 * TILE_SIZE, 0},
		common.ICON_CODE_BOMB:   {10 * TILE_SIZE, 0},
		common.ICON_CODE_OPAQUE: {11 * TILE_SIZE, 0},
	}
	BYTE_ORDER = common.BYTE_ORDER
)

//go:embed tiles.png
var tilesPng []byte

type GameClient struct {
	World           ClientWorld
	Atlas           *ebiten.Image
	BoardX          float64
	BoardY          float64
	Input           Input
	Score           uint32
	Frame           uint64
	RecieveMessages <-chan []byte
	SendMessages    chan<- []byte
	Log             logger.Logger
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
	for i := range WORLD_TILE_COUNT {
		tilePos := coord.CoordFromIndex(i, TY_SHIFT, TX_MASK)
		boardPos := tilePos.MultScalar(TILE_SIZE).DivScalar(DISPLAY_SCALE_DOWN)
		iconIdx := g.World.Tiles[i]
		iconTopLeft := BOARD_TILES[iconIdx]
		iconBotRight := [2]int{iconTopLeft[0] + TILE_SIZE, iconTopLeft[1] + TILE_SIZE}
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

func (g *GameClient) Init(world *ServerWorld, clientToServer chan<- []byte, serverToClient <-chan []byte) {
	g.World = ClientWorld{
		Id:            world.Id.Load(),
		TotalMines:    world.TotalMines,
		ExplodedMines: world.ExplodedMines.Load(),
		SweptTiles:    world.SweptTiles.Load(),
		Ended:         world.Ended.Load(),
		Expires:       world.Expires,
	}

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
				w := wire.NewIncomingSlice(msg, wire.LE)
				var msgCode uint32
				w.TryRead_U32(&msgCode)
				switch msgCode {
				case MSG.SERVER_SWEEP:
					var sweep SweepResult
					sweep.WireRead(&w)
					g.Score += uint32(sweep.Score)
					sweep.DoActionOnAllTiles(func(pos Coord, icon byte) {
						idx := pos.ToIndex(TY_SHIFT)
						g.World.Tiles[idx] = icon
					})
				default:
					g.Log.Warn("invalid msg code: %d", msgCode)
				}
			default:
				more_messages = false
			}
		}
	}
	{ // Update State and Send Messages
		g.BoardX += g.Input.ScrollX * WHEEL_SPEED
		g.BoardY += g.Input.ScrollY * WHEEL_SPEED
		g.BoardX = xmath.Clamp(MIN_BOARD_POS_X, g.BoardX, MAX_BOARD_POS_X)
		g.BoardY = xmath.Clamp(MIN_BOARD_POS_Y, g.BoardY, MAX_BOARD_POS_Y)
		if g.Input.MouseLJustPressed {
			tilePos := g.MousePosToTilePos()
			tileIdx := tilePos.ToIndex(TY_SHIFT)
			if g.World.Tiles[tileIdx] > 8 { //FIXME make `ClientTile` type with readable methods (cheking whether tile is not swept here)

				request := common.NewSweepRequest(tilePos)
				buf := bytes.Buffer{}
				buf.Grow(64)
				outWire := wire.NewOutgoing(&buf, wire.LE)
				request.WireWrite(&outWire)
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
	fx /= float64(TILE_SIZE_SCALED)
	fy /= float64(TILE_SIZE_SCALED)
	return Coord{
		X: int(fx),
		Y: int(fy),
	}
}

var _ ebiten.Game = (*GameClient)(nil)
