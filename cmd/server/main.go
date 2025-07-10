package main

import (
	"os"

	"github.com/gabe-lee/OurSweeper/internal/server_world"
)

func main() {
	w := server_world.World{}
	w.InitNew(1)
	// w.DrawNearby(os.Stdout)
	// w.DrawState(os.Stdout)
	w.PrintStatus(os.Stdout)
}
