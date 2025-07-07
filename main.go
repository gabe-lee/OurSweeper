package main

import (
	"os"

	"github.com/gabe-lee/OurSweeper/internal/world"
)

func main() {
	w := world.World{}
	w.InitNew(1)
	w.DrawNearby(os.Stdout)
	w.DrawState(os.Stdout)
}
