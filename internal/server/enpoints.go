package server

import (
	"net/http"

	"github.com/gabe-lee/OurSweeper/logger"
	"github.com/gobwas/ws"
)

func (n *WebNetwork) RegisterEndpoints(mux *ServeMux) {
	mux.HandleFunc("/", n.dummyPath)
	mux.HandleFunc("GET /game", n.dummyPath)
	// mux.HandleFunc("GET /api/anontoken", s.NewAnonToken)
	// mux.HandleFunc("/login")
	// mux.HandleFunc("PUT /api/anontoken", s.RefreshAnonToken)
}

func (n *WebNetwork) ConnectToGame(w http.ResponseWriter, req *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(req, w)
	if n.WriteErrorReaponseAndLogIfErr(err, logger.ERROR, w, http.StatusBadRequest, "bad game connection request") {
		return
	}
	n.websocketWait.Add(1)
	go n.StartGameConn(conn)
}

// func (s *WebNetwork) GetWorldState(w http.ResponseWriter, req *http.Request) {
// 	defer func() {
// 		if req.Body != nil {
// 			req.Body.Close()
// 		}
// 	}()
// 	UUID, err := uuid.NewRandom()
// 	if s.WriteErrorReaponseAndLogIfErr(err, logger.ERROR, w, http.StatusInternalServerError, "could not create new UUID") {
// 		return
// 	}
// 	anonToken := AnonToken{
// 		UUID:    UUID,
// 		Version: user_token.CurrentVer,
// 	}
// 	var tokenRawBytes []byte
// 	tokenRawBytes, err = token.Create(s.env.tokenSecret, &anonToken)
// 	if s.WriteErrorReaponseAndLogIfErr(err, logger.ERROR, w, http.StatusInternalServerError, "could not serialize new anon token") {
// 		return
// 	}
// 	tokenRaw := AnonTokenRaw{
// 		Token: tokenRawBytes,
// 	}
// 	write := wire.NewOutgoing(w, bin)
// 	write.Struct(&tokenRaw)
// 	s.log.ErrorIfErr(write.Err(), "failed to write anon token to http response")
// }

// func (s *WebNetwork) GameLogin(w http.ResponseWriter, req *http.Request) {
// 	conn, _, _, err := ws.UpgradeHTTP(req, w)
// 	if s.WriteErrorReaponseAndLogIfErr(err, logger.ERROR, w, http.StatusInternalServerError, "websocket upgrade failure") {
// 		return
// 	}
// 	go func() {
// 		defer func() {
// 			conn.Close()
// 		}()
// 		var inWireBuf = data_buffer.NewWriteBuffer(256)
// 		var inWire  = wire.NewIncoming()
// 		var binWriter = wsutil.NewWriterSize(conn, ws.StateServerSide, ws.OpBinary, 1024)
// 		var binReader = wsutil.NewReader(conn, ws.StateServerSide)
// 		for {
// 			conn.
// 			header, err := binReader.NextFrame()
// 			if err != nil {
// 				binWriter.
// 			}
// 			switch header.OpCode {
// 			case ws.OpClose:
// 				return
// 			}
// 		}
// 	}()
// }

// func (s *ServerNetworkManager) RefreshAnonToken(w http.ResponseWriter, req *http.Request) {
// 	defer func() {
// 		if req.Body != nil {
// 			req.Body.Close()
// 		}
// 	}()
// 	validToken := true
// 	if logid := s.log.NoteIfTrue(req.Body == nil || req.ContentLength == 0, "no body provided to `RefreshAnonToken()`"); logid != 0 {
// 		validToken = false
// 	}
// 	var anonTokenRaw AnonTokenRaw
// 	if validToken {
// 		read := wire.NewIncoming(req.Body, bin)
// 		read.WireReader(&anonTokenRaw)
// 		if logid := s.log.NoteIfErr(read.Err(), "could not deserialize AnonTokenRaw"); logid != 0 {
// 			validToken = false
// 		}
// 	}
// 	var anonToken AnonToken
// 	if validToken {
// 		validToken, err := token.OpenAndValidate(s.env.tokenSecret, anonTokenRaw.Token, &anonToken)
// 		if logid := s.log.NoteIfErr(err, "could not deserialize AnonToken"); logid != 0 {
// 			validToken = false
// 		} else if logid := s.log.NoteIfTrue(!validToken, "an AnonToken was corrupted and refused"); logid != 0 {
// 			validToken = false
// 		}
// 	}
// 	if !validToken {
// 		token.Create()
// 	}
// 	if validToken {
// 		s.ClientManager.ClientMapLock.Lock()
// 		client, exists := s.ClientManager.ClientMap[anonToken.UUID]
// 		s.ClientManager.ClientMapLock.Unlock()
// 		if s.WriteErrorReaponseAndLogIfTrue(!exists, logger.WARN, w, http.StatusNotFound, "an AnonToken was corrupted and refused") {
// 			return
// 		}
// 		if logid := s.log.NoteIfTrue(!exists, "Anon"); logid != 0 {
// 			validToken = false
// 		}
// 	}

// uuid.EnableRandPool()
// UUID, err := uuid.NewRandom()
// if s.WriteErrorReaponseAndLogIfErr(err, w, http.StatusInternalServerError, error_response.ERR_UNKNOWN, "could not generate uuid for anon token") {
// 	return
// }
// newToken := AnonToken{
// 	UUID:    UUID,
// 	Version: anon_token.CurrentVer,
// }
// tokenRawBytes, err := token.Create(s.env.tokenSecret, &newToken)
// if s.WriteErrorReaponseAndLogIfErr(err, w, http.StatusInternalServerError, error_response.ERR_UNKNOWN, "could not serialize new anon token") {
// 	return
// }
// tokenRaw := AnonTokenRaw{
// 	Token: tokenRawBytes,
// }
// outWire := wire.NewOutgoing(w, wire.LE)
// tokenRaw.WireWrite(&outWire)
// if s.WriteErrorReaponseAndLogIfErr(outWire.Err(), w, http.StatusInternalServerError, error_response.ERR_UNKNOWN, "failed to serialize anon token") {
// 	return
// }
// }

func (n *WebNetwork) WriteErrorReaponseAndLogIfErr(err error, logLevel int, w http.ResponseWriter, code int, serverLog string, serverLogArgs ...any) (fail bool) {
	if err != nil {
		fail = true
		n.log.LogIfErr(logLevel, err, serverLog, serverLogArgs...)
		w.WriteHeader(code)
		var arr [2]byte
		_, err = w.Write(arr[:])
		n.log.ErrorIfErr(err, "failed to write error code to error response")
	}
	return fail
}

func (n *WebNetwork) WriteErrorReaponseAndLogIfTrue(cond bool, logLevel int, w http.ResponseWriter, code int, serverLog string, serverLogArgs ...any) (fail bool) {
	if cond {
		fail = true
		n.log.Log(logLevel, serverLog, serverLogArgs...)
		w.WriteHeader(code)
		var arr [2]byte
		_, err := w.Write(arr[:])
		n.log.ErrorIfErr(err, "failed to write error code to error response")
	}
	return fail
}
