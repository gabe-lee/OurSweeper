package server_network

import (
	"log"
	"net/http"
	"time"

	"github.com/gabe-lee/OurSweeper/env_loader"
	"github.com/gabe-lee/OurSweeper/internal/database"
	"github.com/gabe-lee/OurSweeper/internal/server_client_manager"
	"github.com/gabe-lee/OurSweeper/internal/server_utility"
	"github.com/gabe-lee/OurSweeper/internal/server_world_manager"
	"github.com/gabe-lee/OurSweeper/internal/user_token"
	"github.com/gabe-lee/OurSweeper/logger"
	"github.com/gabe-lee/OurSweeper/wire"
	"github.com/gobwas/ws"
)

const (
	timeoutHeaderRead = time.Second * 5
)

type (
	Logger          = logger.Logger
	SubLogger       = logger.SubLogger
	SubLoggerWriter = logger.SubLoggerWriter
	SweepDB         = database.SweepDB
	AnonToken       = user_token.UserStats
	AnonTokenRaw    = user_token.AnonTokenRaw
	ServeMux        = http.ServeMux
	ClientManager   = server_client_manager.ClientManager
	ServerUtility   = server_utility.ServerUtility
	WorldManager    = server_world_manager.WorldManager
)

var bin = wire.LE

type ServerNetworkManager struct {
	ClientManager ClientManager
	WorldManager  WorldManager
	server        http.Server
	env           ServerNetworkEnv
	log           SubLogger
	logWriter     SubLoggerWriter
	wsUpgrader    ws.Upgrader
	utils         *ServerUtility
}

func NewServerNetworkManager(masterLogger *Logger, database *SweepDB, utils *ServerUtility) ServerNetworkManager {
	mux := http.NewServeMux()
	subLogger := masterLogger.NewSubLogger("Network")
	subLoggerWriter := subLogger.NewSubLoggerWriter(logger.ERROR)
	s := ServerNetworkManager{
		log:       subLogger,
		logWriter: subLoggerWriter,
		env:       ServerNetworkEnv{},
	}
	env_loader.LoadAndFill(&s.env, &subLoggerWriter)
	s.server = http.Server{
		Addr:              string(s.env.host) + ":" + string(s.env.port),
		Handler:           mux,
		ReadHeaderTimeout: timeoutHeaderRead,
		ErrorLog:          log.New(&subLoggerWriter, "", 0),
	}
	s.RegisterEndpoints(mux)
	return s
}

func (s *ServerNetworkManager) Start() {
	s.log.Info("Server listening at https://%s", s.server.Addr)
	s.log.ErrorIfErr(s.server.ListenAndServeTLS(s.env.tlsCertFile, s.env.tlsKeyFile), "server network error")
}

type ServerNetworkEnv struct {
	port               []byte `env:"SERVER_LISTEN_PORT" default:"8080"`
	host               []byte `env:"SERVER_LISTEN_HOST" default:"127.0.0.1"`
	tlsCertFile        string `env:"SERVER_TLS_CERT_FILE" default:"localhost.crt"`
	tlsKeyFile         string `env:"SERVER_TLS_KEY_FILE" default:"localhost.key"`
	tokenSecret        []byte `env:"SERVER_TOKEN_SECRET" default:"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
	gameListenProtocol string `env:"SERVER_GAME_LISTEN_PROTOCOL" default:"tcp"`
	gameListenAddr     string `env:"SERVER_GAME_LISTEN_ADDR" default:"127.0.0.1"`
	gameListenPort     string `env:"SERVER_GAME_LISTEN_PORT" default:"8081"`
}

func (s *ServerNetworkManager) dummyPath(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("<html><body><h1>COOL</h1></body></html>"))
	s.log.ErrorIfErr(err, "could not serve cool page")
}

// s.wsUpgrader = ws.Upgrader{
// 		OnHeader: func(h_key, h_value []byte) error {
// 			switch {
// 			case bytes.Equal(h_key, H.COOKIE):
// 				hasAnyValidToken := false
// 				httphead.ScanCookie(h_value, func(c_key, c_value []byte) bool {
// 					switch {
// 					case bytes.Equal(c_key, H.COOKIE_ANON_TOKEN):
// 						var userStats user_token.UserStats
// 						valid, err := token.OpenAndValidate(s.env.tokenSecret, c_value, &userStats)
// 						if !valid {
// 							s.log.Warn("corrupted AnonToken provided")
// 							return false
// 						}
// 						if err != nil {
// 							s.log.WarnIfErr(err, "unable to decode AnonToken")
// 							return false
// 						}
// 						alreadyLogged := s.ClientManager.UserIsLoggedOn(userStats.UUID)
// 						if alreadyLogged {
// 							return false
// 						}
// 						hasAnyValidToken = true
// 						return false
// 					default:
// 						return true
// 					}
// 				})
// 				if hasAnyValidToken {
// 					return nil
// 				}
// 				return ws.RejectConnectionError(
// 					ws.RejectionStatus(http.StatusForbidden),
// 					ws.RejectionReason("No Valid")
// 				)
// 			default:
// 				return nil
// 			}
// 		},
// 		OnHost: func(host []byte) error {
// 			if bytes.Equal(s.env.host, host) {
// 				return nil
// 			}
// 			s.log.Warn("invalid host requested `%s`", host)
// 			buf := s.utils.WriteBufferPool.QuickWriteBuffer(H.X_WANT_HOST, ": ", s.env.host, "\r\n")
// 			defer s.utils.WriteBufferPool.ReleaseBuffer(buf)
// 			return ws.RejectConnectionError(
// 				ws.RejectionStatus(http.StatusForbidden),
// 				ws.RejectionHeader(ws.HandshakeHeader(buf)),
// 			)
// 		},
// 		OnRequest: func(uri []byte) error {
// 			return nil
// 		},
// 	}
