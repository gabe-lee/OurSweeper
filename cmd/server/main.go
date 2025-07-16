package main

import (
	"os"

	App "github.com/gabe-lee/OurSweeper/internal"
	_ "modernc.org/sqlite"
)

func main() {
	db := App.NewSweepDB(os.Stdout)
	db.CheckFile()
	db.Open()
	defer db.Close()
	var easyWorld App.ServerWorld
	if !db.GetActiveWorld(&easyWorld, App.DIFFICULTY_EASY) {
		db.CreateNewWorld(&easyWorld, App.DIFFICULTY_EASY)
	} else {
		db.LoadAllChunks(&easyWorld)
	}
}
