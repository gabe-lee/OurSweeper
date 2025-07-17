package wire_codes

const (
	PING uint32 = iota
	PONG
	CLIENT_ERROR
	SERVER_ERROR

	ANON_TOKEN_REFRESH
	ANON_TOKEN_NEW
)
