package server

import (
	"context"
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

type WebNetwork struct {
	web_server    http.Server
	env           ServerNetworkEnv
	log           SubLogger
	logWriter     SubLoggerWriter
	clients       *ClientManager
	worlds        *WorldManager
	utils         *ServerUtility
	shutdownSig   chan struct{}
	websocketWait MiniWait
}

func (s *WebNetwork) Init(utils *ServerUtility, masterLogger *Logger, clients *ClientManager, worlds *WorldManager) {
	s.log = masterLogger.NewSubLogger("Network")
	s.logWriter = s.log.NewSubLoggerWriter(logger.ERROR)
	s.utils = utils
	s.clients = clients
	s.worlds = worlds
	s.shutdownSig = make(chan struct{})
	env_loader.LoadAndFill(&s.env, &s.logWriter)
	mux := http.NewServeMux()
	s.web_server = http.Server{
		Addr:              string(s.env.host) + ":" + string(s.env.port),
		Handler:           mux,
		ReadHeaderTimeout: timeoutHeaderRead,
		ErrorLog:          log.New(&s.logWriter, "", 0),
	}
	s.RegisterEndpoints(mux)
}

func (s *WebNetwork) Start() {
	go func() {
		s.log.Info("Web Server listening at https://%s", s.web_server.Addr)
		err := s.web_server.ListenAndServeTLS(s.env.tlsCertFile, s.env.tlsKeyFile)
		if err != nil && err != http.ErrServerClosed {
			s.log.ErrorIfErr(err, "server network error")
		}
	}()
}

func (s *WebNetwork) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	close(s.shutdownSig)
	err := s.web_server.Shutdown(ctx)
	s.log.WarnIfErr(err, "Server Shutdown: http network failed to close gracefully")
	allWebsocketsClosed := s.websocketWait.WaitWithTimeout(time.Microsecond*100, time.Second*10)
	s.log.WarnIfTrue(!allWebsocketsClosed, "Server Shutdown: %d game websockets failed to close gracefully", s.websocketWait.RemainingWaitCount())
	s.log.Close()
}

type ServerNetworkEnv struct {
	port        []byte `env:"SERVER_LISTEN_PORT" default:"8080"`
	host        []byte `env:"SERVER_LISTEN_HOST" default:"127.0.0.1"`
	tlsCertFile string `env:"SERVER_TLS_CERT_FILE" default:"localhost.crt"`
	tlsKeyFile  string `env:"SERVER_TLS_KEY_FILE" default:"localhost.key"`
	tokenSecret []byte `env:"SERVER_TOKEN_SECRET" default:"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
	// gameListenProtocol string `env:"SERVER_GAME_LISTEN_PROTOCOL" default:"tcp"`
	// gameListenAddr     string `env:"SERVER_GAME_LISTEN_ADDR" default:"127.0.0.1"`
	// gameListenPort     string `env:"SERVER_GAME_LISTEN_PORT" default:"8081"`
}

func (s *WebNetwork) dummyPath(w http.ResponseWriter, req *http.Request) {
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
