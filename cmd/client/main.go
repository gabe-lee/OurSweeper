package main

import (
	"log"

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
	ebiten.SetWindowSize(game_client.WINDOW_WIDTH, game_client.WINDOW_HEIGHT)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&client); err != nil {
		log.Fatal(err)
	}
}
