package logger

import (
	"github.com/gabe-lee/OurSweeper/xmath"
)

type SubLoggerWriter struct {
	sl   *SubLogger
	mode int
}

func (slw *SubLoggerWriter) Write(p []byte) (n int, err error) {
	defer slw.sl.logger.errs.Flush("error SubLoggerWriter.Write(): ")
	slw.sl.logger.log(slw.sl, slw.mode, nil, "%s", string(p))
	return len(p), nil
}

func (lw *SubLoggerWriter) SetMode(mode int) {
	lw.mode = xmath.Clamp(0, mode, levelCount-1)
}

type LoggerWriter struct {
	l    *Logger
	mode int
}

func (lw *LoggerWriter) Write(p []byte) (n int, err error) {
	defer lw.l.errs.Flush("error LoggerWriter.Write(): ")
	lw.l.log(nil, lw.mode, nil, "%s", string(p))
	return len(p), nil
}
func (lw *LoggerWriter) SetMode(mode int) {
	lw.mode = xmath.Clamp(0, mode, levelCount-1)
}
