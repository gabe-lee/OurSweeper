package error_response

type ErrorResponse uint16

const (
	ERR_UNKNOWN uint16 = iota
	ERR_ANON_TOKEN_EXPIRED
	ERR_ANON_TOKEN_CORRUPTED
	count
)

var MSG_TABLE = [count]string{
	ERR_UNKNOWN:              "Unknown error",
	ERR_ANON_TOKEN_CORRUPTED: "Anonymous login token was corrupted, all stats lost",
	ERR_ANON_TOKEN_EXPIRED:   "Anonymous login token was expired, all stats lost",
}
