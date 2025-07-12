package game_client

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/gabe-lee/OurSweeper/internal/client_world"
	C "github.com/gabe-lee/OurSweeper/internal/common"
	"github.com/gabe-lee/OurSweeper/internal/coord"
	"github.com/gabe-lee/OurSweeper/internal/server_world"
	"github.com/gabe-lee/OurSweeper/internal/sweep_request"
	"github.com/gabe-lee/OurSweeper/internal/tile"
	"github.com/gabe-lee/OurSweeper/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	Coord       = coord.Coord
	ServerWorld = server_world.World
	ClientWorld = client_world.ClientWorld
	EbitImage   = ebiten.Image
)

//go:embed tiles.png
var tilesPng []byte

type GameClient struct {
	Server *ServerWorld //Temporary until network implemented
	World  ClientWorld
	Atlas  *ebiten.Image
	BoardX float64
	BoardY float64
	Input  Input
	Score  uint32
	Frame  uint64
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
		tilePos := coord.FromIndex(i, C.TY_SHIFT, C.TX_MASK)

		boardPos := tilePos.MultScalar(C.TILE_SIZE).DivScalar(C.DISPLAY_SCALE_DOWN)
		iconIdx := g.World.Tiles[i].GetIconClient()
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

func (g *GameClient) Init(world *ServerWorld) {
	g.World = client_world.ClientWorld{
		Id:            world.Id.Load(),
		TotalMines:    world.TotalMines,
		ExplodedMines: world.ExplodedMines.Load(),
		SweptTiles:    world.SweptTiles.Load(),
		Ended:         world.Ended.Load(),
		Expires:       world.Expires,
	}
	g.Server = world // Temp
	for i := range C.WORLD_TILE_COUNT {
		g.World.Tiles[i] = tile.Tile(g.Server.Tiles[i].GetIconServer())
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
	// Update State
	g.BoardX += g.Input.ScrollX * C.WHEEL_SPEED
	g.BoardY += g.Input.ScrollY * C.WHEEL_SPEED
	g.BoardX = utils.Clamp(C.MIN_BOARD_POS_X, g.BoardX, C.MAX_BOARD_POS_X)
	g.BoardY = utils.Clamp(C.MIN_BOARD_POS_Y, g.BoardY, C.MAX_BOARD_POS_Y)
	if g.Input.MouseLJustPressed {
		tilePos := g.MousePosToTilePos()
		request := sweep_request.NewSweepRequest(tilePos)
		//TODO send network request instead
		result := g.Server.SweepTile(request.Pos.ToCoord())
		//TODO listen for response
		for i := range result.Len {
			newPos := result.Coords[i].ToCoord()
			icon := result.Icons[i]
			idx := newPos.ToIndex(C.TY_SHIFT)
			g.World.Tiles[idx] = tile.Tile(icon)
			g.Score += uint32(result.Score)
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
