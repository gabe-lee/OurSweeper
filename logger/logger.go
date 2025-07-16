package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gabe-lee/OurSweeper/ansi"
)

const (
	WARN  = " WARN "
	FATAL = " FATAL "
	INFO  = " INFO "
)

type Logger struct {
	Name   string
	Writer io.Writer
}

func NewLogger(name string, writer io.Writer) Logger {
	return Logger{
		Name:   name,
		Writer: writer,
	}
}

func (l *Logger) Fatal(format string, args ...any) {
	fmt.Fprintf(l.Writer, ansi.FG_RED+"["+l.Name+FATAL+"%s] ", time.Now().Format(time.DateTime))
	fmt.Fprintf(l.Writer, format+ansi.CLEAR+"\n", args...)
	os.Exit(1)
}

func (l *Logger) Warn(format string, args ...any) {
	fmt.Fprintf(l.Writer, ansi.FG_YEL+"["+l.Name+WARN+"%s] ", time.Now().Format(time.DateTime))
	fmt.Fprintf(l.Writer, format+ansi.CLEAR+"\n", args...)
}

func (l *Logger) Info(format string, args ...any) {
	fmt.Fprintf(l.Writer, ansi.FG_BLU+"["+l.Name+INFO+"%s] ", time.Now().Format(time.DateTime))
	fmt.Fprintf(l.Writer, format+ansi.CLEAR+"\n", args...)
}

func (l *Logger) FatalIfErr(err error, format string, args ...any) {
	if err != nil {
		fmt.Fprintf(l.Writer, ansi.FG_RED+"["+l.Name+WARN+"%s] ", time.Now().Format(time.DateTime))
		fmt.Fprintf(l.Writer, format, args...)
		fmt.Fprintf(l.Writer, ": %s\n"+ansi.CLEAR, err)
		os.Exit(1)
	}
}

func (l *Logger) WarnIfErr(err error, format string, args ...any) {
	if err != nil {
		fmt.Fprintf(l.Writer, ansi.FG_YEL+"["+l.Name+WARN+"%s] ", time.Now().Format(time.DateTime))
		fmt.Fprintf(l.Writer, format, args...)
		fmt.Fprintf(l.Writer, ": %s\n"+ansi.CLEAR, err)
	}
}

func (l *Logger) InfoIfErr(err error, format string, args ...any) {
	if err != nil {
		fmt.Fprintf(l.Writer, ansi.FG_BLU+"["+l.Name+INFO+"%s] ", time.Now().Format(time.DateTime))
		fmt.Fprintf(l.Writer, format, args...)
		fmt.Fprintf(l.Writer, ": %s\n"+ansi.CLEAR, err)
	}
}
