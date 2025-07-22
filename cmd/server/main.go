package main

import (
	"os"

	"github.com/gabe-lee/OurSweeper/internal/common"
	"github.com/gabe-lee/OurSweeper/internal/database"
	"github.com/gabe-lee/OurSweeper/internal/server_network"
	"github.com/gabe-lee/OurSweeper/logger"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type (
	ServerWorld = common.ServerWorld
)

const (
	DIFFICULTY_EASY = common.DIFFICULTY_EASY

	uuidRandPool = true
)

const (
	logDir    = "/logs"
	masterDir = "Master"
)

func main() {
	if uuidRandPool {
		uuid.EnableRandPool()
	}
	masterLogger := logger.NewLogger(logDir, masterDir, os.Stdout, 4)
	defer masterLogger.Close()
	db := database.NewSweepDB(&masterLogger)
	db.CheckFile()
	db.Open()
	defer db.Close()
	var easyWorld ServerWorld
	if !db.GetActiveWorld(&easyWorld, DIFFICULTY_EASY) {
		db.CreateNewWorld(&easyWorld, DIFFICULTY_EASY)
	} else {
		db.LoadAllChunks(&easyWorld)
	}
	network := server_network.NewServerNetworkManager(&masterLogger, &db)
	network.Start()
}
