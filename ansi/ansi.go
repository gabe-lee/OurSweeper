package ansi

import "io"

const (
	CLEAR = "\x1b[0m"

	FG_BLK = "\x1b[30m"
	FG_RED = "\x1b[31m"
	FG_GRN = "\x1b[32m"
	FG_YEL = "\x1b[33m"
	FG_BLU = "\x1b[34m"
	FG_MAG = "\x1b[35m"
	FG_CYA = "\x1b[36m"
	FG_WHT = "\x1b[37m"

	BG_BLK = "\x1b[40m"
	BG_RED = "\x1b[41m"
	BG_GRN = "\x1b[42m"
	BG_YEL = "\x1b[43m"
	BG_BLU = "\x1b[44m"
	BG_MAG = "\x1b[45m"
	BG_CYA = "\x1b[46m"
	BG_WHT = "\x1b[47m"

	INV_BLK = "\x1b[7;30m"
	INV_RED = "\x1b[7;31m"
	INV_GRN = "\x1b[7;32m"
	INV_YEL = "\x1b[7;33m"
	INV_BLU = "\x1b[7;34m"
	INV_MAG = "\x1b[7;35m"
	INV_CYA = "\x1b[7;36m"
	INV_WHT = "\x1b[7;37m"

	INV = "\x1b[7m"
)

func write(w io.Writer, pre string, str string) (n int, err error) {
	nn, _ := w.Write([]byte(pre))
	n += nn
	nn, err = w.Write([]byte(str))
	n += nn
	nn, _ = w.Write([]byte(CLEAR))
	n += nn
	return n, err
}

func FgBlk(w io.Writer, str string) (n int, err error) {
	return write(w, FG_BLK, str)
}
func FgRed(w io.Writer, str string) (n int, err error) {
	return write(w, FG_RED, str)
}
func FgGrn(w io.Writer, str string) (n int, err error) {
	return write(w, FG_GRN, str)
}
func FgYel(w io.Writer, str string) (n int, err error) {
	return write(w, FG_YEL, str)
}
func FgBlu(w io.Writer, str string) (n int, err error) {
	return write(w, FG_BLU, str)
}
func FgMag(w io.Writer, str string) (n int, err error) {
	return write(w, FG_MAG, str)
}
func FgCya(w io.Writer, str string) (n int, err error) {
	return write(w, FG_CYA, str)
}
func FgWht(w io.Writer, str string) (n int, err error) {
	return write(w, FG_WHT, str)
}

func BgBlk(w io.Writer, str string) (n int, err error) {
	return write(w, BG_BLK, str)
}
func BgRed(w io.Writer, str string) (n int, err error) {
	return write(w, BG_RED, str)
}
func BgGrn(w io.Writer, str string) (n int, err error) {
	return write(w, BG_GRN, str)
}
func BgYel(w io.Writer, str string) (n int, err error) {
	return write(w, BG_YEL, str)
}
func BgBlu(w io.Writer, str string) (n int, err error) {
	return write(w, BG_BLU, str)
}
func BgMag(w io.Writer, str string) (n int, err error) {
	return write(w, BG_MAG, str)
}
func BgCya(w io.Writer, str string) (n int, err error) {
	return write(w, BG_CYA, str)
}
func BgWht(w io.Writer, str string) (n int, err error) {
	return write(w, BG_WHT, str)
}

func InvBlk(w io.Writer, str string) (n int, err error) {
	return write(w, INV_BLK, str)
}
func InvRed(w io.Writer, str string) (n int, err error) {
	return write(w, INV_RED, str)
}
func InvGrn(w io.Writer, str string) (n int, err error) {
	return write(w, INV_GRN, str)
}
func InvYel(w io.Writer, str string) (n int, err error) {
	return write(w, INV_YEL, str)
}
func InvBlu(w io.Writer, str string) (n int, err error) {
	return write(w, INV_BLU, str)
}
func InvMag(w io.Writer, str string) (n int, err error) {
	return write(w, INV_MAG, str)
}
func InvCya(w io.Writer, str string) (n int, err error) {
	return write(w, INV_CYA, str)
}
func InvWht(w io.Writer, str string) (n int, err error) {
	return write(w, INV_WHT, str)
}

func Invert(w io.Writer, str string) (n int, err error) {
	return write(w, INV, str)
}
