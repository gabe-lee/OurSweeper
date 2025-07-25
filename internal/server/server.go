package server

import (
	"context"

	"github.com/gabe-lee/OurSweeper/lock"
)

const (
	SIG_SHOULD_CLOSE byte = iota
	SIG_DID_CLOSE_ERR
)

type (
	Context    = context.Context
	CancelFunc = context.CancelFunc
	MiniWait   = lock.MiniWait
)

const (
	SHUT_WEB_NET  uint32 = 1
	SHUT_GAME_NET uint32 = 2
	SHUT_DATABASE uint32 = 4
	SHUT_TOTAL    uint32 = SHUT_WEB_NET | SHUT_GAME_NET | SHUT_DATABASE
	SHUT_NETWORKS uint32 = SHUT_WEB_NET | SHUT_GAME_NET
)

type Server struct {
	Database       SweepDB
	WebNetwork     WebNetwork
	ClientManager  ClientManager
	WorldManager   WorldManager
	Logger         Logger
	Utils          ServerUtility
}

func (s *Server) Start() {
	s.Database.Start()
	s.WebNetwork.Start()
}

func (s *Server) Stop() {
	s.Logger.Info("Beginning server shutdown...")
	s.Logger.Info("Waiting on network shutdown...")
	s.WebNetwork.Close()
	s.Logger.Info("Waiting on database shutdown...")
	s.Database.Close()
	s.Logger.Close()
}
