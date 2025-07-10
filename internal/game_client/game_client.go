package game_client

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/gabe-lee/OurSweeper/internal/client_world"
	"github.com/gabe-lee/OurSweeper/internal/server_world"
	"github.com/gabe-lee/OurSweeper/internal/tile"
	"github.com/gabe-lee/OurSweeper/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TILE_SIZE        int = 32
	TILE_SIZE_SCALED int = TILE_SIZE / SCALE_DOWN

	TILE_SHEET_WIDTH  int = 12
	TILE_SHEET_HEIGHT int = 4

	WINDOW_WIDTH     int     = 800
	WINDOW_HEIGHT    int     = 800
	BOARD_WIDTH      int     = TILE_SIZE_SCALED * server_world.WIDTH
	BOARD_HEIGHT     int     = TILE_SIZE_SCALED * server_world.HEIGHT
	BOARD_OVERFLOW_X int     = BOARD_WIDTH - WINDOW_WIDTH
	BOARD_OVERFLOW_Y int     = BOARD_HEIGHT - WINDOW_HEIGHT
	MIN_BOARD_POS_X  float64 = float64(-BOARD_OVERFLOW_X)
	MIN_BOARD_POS_Y  float64 = float64(-BOARD_OVERFLOW_Y)
	MAX_BOARD_POS_X  float64 = 0
	MAX_BOARD_POS_Y  float64 = 0
	SCALE_DOWN       int     = 2
	WHEEL_SPEED      float64 = 6.0
)

var BOARD_TILES = [16][2]int{
	tile.ICON_CODE_0:      {0, 0},
	tile.ICON_CODE_1:      {1 * TILE_SIZE, 0},
	tile.ICON_CODE_2:      {2 * TILE_SIZE, 0},
	tile.ICON_CODE_3:      {3 * TILE_SIZE, 0},
	tile.ICON_CODE_4:      {4 * TILE_SIZE, 0},
	tile.ICON_CODE_5:      {5 * TILE_SIZE, 0},
	tile.ICON_CODE_6:      {6 * TILE_SIZE, 0},
	tile.ICON_CODE_7:      {7 * TILE_SIZE, 0},
	tile.ICON_CODE_8:      {8 * TILE_SIZE, 0},
	tile.ICON_CODE_FLAG:   {9 * TILE_SIZE, 0},
	tile.ICON_CODE_BOMB:   {10 * TILE_SIZE, 0},
	tile.ICON_CODE_OPAQUE: {11 * TILE_SIZE, 0},
}

//go:embed tiles.png
var tilesPng []byte

type GameClient struct {
	World  client_world.ClientWorld
	Atlas  *ebiten.Image
	BoardX float64
	BoardY float64
}

// Draw implements ebiten.Game.
func (g *GameClient) Draw(screen *ebiten.Image) {
	for i := range server_world.TILES {
		x, y := server_world.GetCoords(i)
		px, py := x*TILE_SIZE/SCALE_DOWN, y*TILE_SIZE/SCALE_DOWN
		iconIdx := g.World.Tiles[i].GetIcon()
		iconTopLeft := BOARD_TILES[iconIdx]
		iconBotRight := [2]int{iconTopLeft[0] + TILE_SIZE, iconTopLeft[1] + TILE_SIZE}
		rect := image.Rect(iconTopLeft[0], iconTopLeft[1], iconBotRight[0], iconBotRight[1])
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(float64(px), float64(py))
		op.GeoM.Translate(g.BoardX, g.BoardY)
		screen.DrawImage(g.Atlas.SubImage(rect).(*ebiten.Image), op)
	}
}

// Layout implements ebiten.Game.
func (g *GameClient) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameClient) Init(world *server_world.World) {
	g.World = client_world.ClientWorld{
		Id:            world.Id.Load(),
		Tiles:         world.Tiles,
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
	wx, wy := ebiten.Wheel()
	g.BoardX += wx * WHEEL_SPEED
	g.BoardY += wy * WHEEL_SPEED
	g.BoardX = utils.Clamp(MIN_BOARD_POS_X, g.BoardX, MAX_BOARD_POS_X)
	g.BoardY = utils.Clamp(MIN_BOARD_POS_Y, g.BoardY, MAX_BOARD_POS_Y)
	return nil
}

var _ ebiten.Game = (*GameClient)(nil)
