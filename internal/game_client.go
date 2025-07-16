package internal

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"log"

	"github.com/gabe-lee/OurSweeper/coord"
	"github.com/gabe-lee/OurSweeper/xmath"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	EbitImage = ebiten.Image
)

const (
	CLIENT_RECIEVE_MESSAGE_BUFFER_LEN int = 32
	CLIENT_SEND_MESSAGE_BUFFER_LEN    int = 4
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
	ErrorWriter     io.Writer
	DebugServer     *ServerWorld //DEBUG
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

	for i := range WORLD_TILE_COUNT { //FIXME
		g.World.Tiles[i] = world.Tiles[i].GetIconForClient()
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
				r := bytes.NewReader(msg)
				var msgCode uint32
				err := binary.Read(r, BYTE_ORDER, &msgCode)
				if err != nil {
					fmt.Fprintf(g.ErrorWriter, "could not read message code: %w", err)
				} else {
					switch msgCode {
					case SERVER_SWEEP:
						var sweep SweepResult
						sweep.Deserialize(r, BYTE_ORDER)
						g.Score += uint32(sweep.Score)
						sweep.DoActionOnAllTiles(func(pos Coord, icon byte) {
							idx := pos.ToIndex(TY_SHIFT)
							g.World.Tiles[idx] = icon
						})
					default:
						fmt.Fprintf(g.ErrorWriter, "invalid msg code: %d", msgCode)
					}
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
			request := NewSweepRequest(tilePos)
			msg, err := request.CodedSerialize(BYTE_ORDER)
			if err != nil {
				fmt.Fprintf(g.ErrorWriter, "failed to write SweepRequest message: %w", err)
			} else {
				g.SendMessages <- msg
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
