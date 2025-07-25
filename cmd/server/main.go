package main

import (
	"bufio"
	"os"

	"github.com/gabe-lee/OurSweeper/env_loader"
	"github.com/gabe-lee/OurSweeper/internal/common"
	_server "github.com/gabe-lee/OurSweeper/internal/server"
	"github.com/gabe-lee/OurSweeper/internal/server_utility"
	"github.com/gabe-lee/OurSweeper/logger"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type (
	ServerWorld = common.ServerWorld
	Server      = _server.Server
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
	// shutdownSig, shutdown := context.WithCancel(context.Background())
	server := Server{
		Logger: logger.NewLogger(logDir, masterDir, os.Stdout, 4),
		Utils:  server_utility.NewServerUtility(),
	}
	logWriter := server.Logger.NewLoggerWriter(logger.ERROR)
	var env ServerAppEnv
	env_loader.LoadAndFill(&env, &logWriter)
	server.Database.Init(&server.Utils, &server.Logger)
	server.WebNetwork.Init(&server.Utils, &server.Logger, &server.ClientManager, &server.WorldManager)
	//CHECKPOINT
	// var easyWorld ServerWorld
	// if !server.Database.GetActiveWorld(&easyWorld, DIFFICULTY_EASY) {
	// 	server.Database.CreateNewWorld(&easyWorld, DIFFICULTY_EASY)
	// } else {
	// 	server.Database.LoadAllChunks(&easyWorld)
	// }
	// do more things?
	server.Start()
	stdin := bufio.NewReader(os.Stdin)
	var input string
	var err error
	for input != env.shutdownPass {
		input, err = stdin.ReadString('\n')
		server.Logger.WarnIfErr(err, "error reading shutdown password from STDIN")
		server.Logger.NormIfTrue(input != env.shutdownPass, "incorrect shutdown password")
	}
	server.Stop()
}

type ServerAppEnv struct {
	shutdownPass string `env:"SERVER_SHUTDOWN_PASSWORD" default:"shutdown\n"`
}
