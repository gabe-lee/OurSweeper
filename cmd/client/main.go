package main

import (
	"log"

	C "github.com/gabe-lee/OurSweeper/internal/common"
	"github.com/gabe-lee/OurSweeper/internal/game_client"
	"github.com/gabe-lee/OurSweeper/internal/server_world"
	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	sWorld := server_world.World{}
	sWorld.InitNew(1)
	client := game_client.GameClient{}
	client.Init(&sWorld)
	ebiten.SetWindowSize(C.WINDOW_WIDTH, C.WINDOW_HEIGHT)
	ebiten.SetWindowTitle("OurSweeper")
	if err := ebiten.RunGame(&client); err != nil {
		log.Fatal(err)
	}
}
