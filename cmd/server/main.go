package main

import (
	App "github.com/gabe-lee/OurSweeper/internal"
	_ "modernc.org/sqlite"
)

func main() {
	db := App.SweepDB{}
	db.CheckFile()
	db.Open()
	defer db.Close()
}
