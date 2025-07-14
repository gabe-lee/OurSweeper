package main

import (
	"os"

	App "github.com/gabe-lee/OurSweeper/internal"
)

func main() {
	w := App.ServerWorld{}
	w.InitNew(1)
	// w.DrawNearby(os.Stdout)
	// w.DrawState(os.Stdout)
	w.PrintStatus(os.Stdout)
}
