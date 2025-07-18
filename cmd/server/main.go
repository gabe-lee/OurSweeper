package main

import (
	"os"
	"sync"
	"time"

	"github.com/gabe-lee/OurSweeper/internal/common"
	"github.com/gabe-lee/OurSweeper/logger"
	_ "modernc.org/sqlite"
)

type (
	ServerWorld = common.ServerWorld
)

const (
	DIFFICULTY_EASY = common.DIFFICULTY_EASY
)

func main() {
	master := logger.NewLogger("/logs", "Master", os.Stdout, 4)
	defer master.Close()
	A := master.NewSubLogger("ThingA")
	defer A.Close()
	B := master.NewSubLogger("ObjectB")
	defer B.Close()
	C := master.NewSubLogger("SomethingC")
	defer C.Close()
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		for range 10 {
			A.Norm("the world")
			A.Norm("revolves")
			A.Norm("around")
			A.Norm("me")
			time.Sleep(time.Millisecond * 3)
		}
		wg.Done()
	}()
	go func() {
		for range 10 {
			B.Note("icecream")
			B.Note("is far too")
			B.Note("delicious")
			time.Sleep(time.Millisecond * 2)
		}
		wg.Done()
	}()
	go func() {
		for range 10 {
			C.Warn("wraning:")
			C.Warn("hot singles")
			C.Warn("are NOT in")
			C.Warn("your area")
			time.Sleep(time.Millisecond * 5)
		}
		wg.Done()
	}()
	wg.Wait()
	return
	// db := database.NewSweepDB(os.Stdout)
	// db.CheckFile()
	// db.Open()
	// defer db.Close()
	// var easyWorld ServerWorld
	// if !db.GetActiveWorld(&easyWorld, DIFFICULTY_EASY) {
	// 	db.CreateNewWorld(&easyWorld, DIFFICULTY_EASY)
	// } else {
	// 	db.LoadAllChunks(&easyWorld)
	// }
}
