package server_network

import (
	"context"
	"net"
	"time"

	"github.com/gabe-lee/OurSweeper/internal/user_token"
	C "github.com/gabe-lee/OurSweeper/internal/wire_codes"
	_wire "github.com/gabe-lee/OurSweeper/wire"
	"github.com/google/uuid"
)

type (
	Listener = net.Listener
	Conn     = net.Conn
)

type GameNetwork struct {
	Listener
}

const (
	CurrentVersion    uint32 = 1
	BadVersionMsg     string = "Unsupported version %d, current version is %d, please update your client to play"
	ConnTimeout              = time.Second * 60
	ConnIdleTime             = time.Second * 60
	ConnProbeInterval        = time.Second * 15
	ConnProbeMax             = 10
	ConnMaxMsg               = 1024
)

const (
	TOKEN_KIND_NONE byte = iota
	TOKEN_KIND_ANON
	TOKEN_KIND_LOGIN
)

var listenConfig = net.ListenConfig{
	KeepAliveConfig: net.KeepAliveConfig{
		Enable:   true,
		Idle:     ConnIdleTime,
		Interval: ConnProbeMax,
		Count:    ConnProbeMax,
	},
}

func (g *GameNetwork) Start(server *ServerNetworkManager, shutdownContext context.Context) (err error) {
	g.Listener, err = listenConfig.Listen(shutdownContext, server.env.gameListenProtocol, server.env.gameListenAddr+":"+server.env.gameListenPort)
	if err != nil {
		return err
	}
	defer g.Close()
	for {
		newConn, err := g.Accept()
		if err != nil {
			return err
		}
		go ServerConn(newConn, server, shutdownContext)
	}
}

func ServerConn(conn Conn, s *ServerNetworkManager, shutdownContext context.Context) {
	wire := _wire.NewBidirection(conn, conn, _wire.LE)
	defer Cleanup(s, conn, wire)
	conn.SetDeadline(time.Now().Add(ConnTimeout))
	var version uint32
	wire.Read.U32(&version)
	if version != CurrentVersion {
		wire.Write.U32(C.SERVER_UNSUPORTED_VERSION)
		HasWriteErr(s, conn, wire)
		return
	}
	var code uint32
	for {
		select {
		case <-shutdownContext.Done():
			return
		default:
		}
		wire.Read.U32(&code)
		if HasReadErr(s, conn, wire) {
			return
		}
		switch code {
		case C.CLIENT_PING:
			wire.Write.U32(C.SERVER_PONG)
			if HasWriteErr(s, conn, wire) {
				return
			}
		case C.CLIENT_ANON_TOKEN_NEW:
			id, err := uuid.NewRandom()
			if err != nil {
				wire.Write.U32(C.SERVER_ERROR)
			} else {
				userToken := user_token.UserStats{
					UUID:    id,
					Version: user_token.CurrentVer,
				}
				wire.Write.U32(C.SERVER_ANON_TOKEN_NEW)
				wire.Write.Struct(&userToken)
			}
			if HasWriteErr(s, conn, wire) {
				return
			}
		case C.CLIENT_GET_ACTIVE_WORLDS:
			//TODO
			//CHECKPOINT
		default:
			wire.Write.U32(C.SERVER_INVALID)
			if HasWriteErr(s, conn, wire) {
				return
			}
		}
	}
}

func Cleanup(s *ServerNetworkManager, conn Conn, wire _wire.Bidirection) {
	conn.Close()
	return
}

func HasReadErr(s *ServerNetworkManager, conn Conn, wire _wire.Bidirection) bool {
	if logid := s.log.NoteIfErr(wire.Read.Err(), "network read error"); logid != 0 {
		return true
	} else {
		conn.SetDeadline(time.Now().Add(ConnTimeout))
	}
	return false
}

func HasWriteErr(s *ServerNetworkManager, conn Conn, wire _wire.Bidirection) bool {
	if logid := s.log.NoteIfErr(wire.Write.Err(), "network write error"); logid != 0 {
		return true
	} else {
		conn.SetDeadline(time.Now().Add(ConnTimeout))
	}
	return false
}
