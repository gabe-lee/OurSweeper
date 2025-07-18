package logger

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gabe-lee/OurSweeper/ansi"
	"github.com/gabe-lee/OurSweeper/utils"
)

type (
	Mutex       = sync.Mutex
	Writer      = io.Writer
	WriteCloser = io.WriteCloser
	Builder     = strings.Builder
	Buffer      = bytes.Buffer
	Stringer    = fmt.Stringer
	Time        = time.Time
)

const (
	NORM int = iota
	INFO
	NOTE
	WARN
	ERROR
	FATAL
	levelCount
)

const (
	logFilePrefix   string = "Log"
	fileExt         string = ".txt"
	fileMaxSize     int    = 1 << 14
	logError        string = "<!!>"
	idPrefix        string = "0x"
	initCap         int    = 128
	sep             byte   = byte('_')
	path_sep        byte   = byte('/')
	date_sep        byte   = byte('/')
	time_sep        byte   = byte(':')
	open_brack      byte   = byte('[')
	space           byte   = byte(' ')
	close_brack     byte   = byte(']')
	newline         byte   = byte('\n')
	logFileFlags           = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	countFileFlags         = os.O_RDWR | os.O_CREATE
	logFilePerms           = 0777 | os.ModeAppend
	countFilePerms         = 0777
	logFileDirPerms        = 0777
	ansiPrefixLen   int    = len(ansi.FG_BLK)
	ansiSuffixLen   int    = len(ansi.CLEAR) + 1
	maxNameLen      int    = 14
	nameFill        string = ". . . . . . . ."
	counterName     string = ".counter"
)

var exNewline = [1]byte{newline}

var prelen = [levelCount]int{
	NORM:  len(color[NORM]),
	INFO:  len(color[INFO]),
	NOTE:  len(color[NOTE]),
	WARN:  len(color[WARN]),
	ERROR: len(color[ERROR]),
	FATAL: len(color[FATAL]),
}

var prefix = [levelCount]string{
	NORM:  "  ",
	INFO:  "--",
	NOTE:  "!!",
	WARN:  "**",
	ERROR: "XX",
	FATAL: "@@",
}

var long = [levelCount]string{
	NORM:  "LOG  ",
	INFO:  "INFO ",
	NOTE:  "NOTE ",
	WARN:  "WARN ",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

var color = [levelCount]string{
	NORM:  "",
	INFO:  ansi.FG_BLU,
	NOTE:  ansi.FG_GRN,
	WARN:  ansi.FG_YEL,
	ERROR: ansi.FG_RED,
	FATAL: ansi.FG_MAG,
}

type Logger struct {
	outDir        string
	masterDir     string
	date          atomic.Uint32
	log_id        atomic.Uint32
	todayFileName StringBuffer
	consoleWriter Writer
	file          *os.File
	counter       *os.File
	fileLock      Mutex
	counterLock   Mutex
	consoleLock   Mutex
	ready_buffers chan StringBuffer
}

func NewLogger(outputDir string, masterDir string, consoleWriter Writer, bufferPool int) Logger {
	cwd, _ := os.Getwd()
	l := Logger{
		outDir:        path.Join(cwd, outputDir),
		masterDir:     masterDir,
		consoleWriter: consoleWriter,
		ready_buffers: make(chan StringBuffer, bufferPool),
		todayFileName: NewStringBuffer(initCap),
	}
	for range bufferPool {
		l.ready_buffers <- NewStringBuffer(initCap)
	}
	fullMaster := path.Join(l.outDir, masterDir)
	err := os.MkdirAll(fullMaster, logFileDirPerms)
	if err != nil {
		fmt.Fprintf(l.consoleWriter, "%sLogger Error: Could not verify or create log directory `%s`: %s\n", logError, fullMaster, err)
	}
	countFile := path.Join(l.outDir, counterName)
	l.counter, err = os.OpenFile(countFile, countFileFlags, countFilePerms)
	if err != nil {
		fmt.Fprintf(l.consoleWriter, "%sLogger Error: Could not open or create counter file `%s`: %s\n", logError, countFile, err)
	}
	var cnt [4]byte
	l.counter.ReadAt(cnt[:], 0)
	log_id := binary.LittleEndian.Uint32(cnt[:])
	l.log_id.Store(log_id)
	now := time.Now()
	date := buildDate(now)
	y1, y2, m, d := unbuildDate(date)
	l.date.Store(date)
	l.makeTodayFileName(&l.todayFileName, masterDir, y1, y2, m, d)
	l.file, err = os.OpenFile(l.todayFileName.String(), logFileFlags, logFilePerms)
	if err != nil {
		fmt.Fprintf(l.consoleWriter, "%sLogger Error: Could not open or create log file `%s`: %s\n", logError, l.todayFileName.String(), err)
	}
	return l
}

func (l *Logger) NewSubLogger(name string) SubLogger {
	n := []byte(nameFill)
	nLen := min(len(name), maxNameLen)
	copy(n, name[:nLen])
	sl := SubLogger{
		logger:        l,
		name:          string(n),
		subDir:        name,
		todayFileName: NewStringBuffer(initCap),
	}
	fullSub := path.Join(l.outDir, sl.subDir)
	err := os.MkdirAll(fullSub, logFileDirPerms)
	if err != nil {
		fmt.Fprintf(l.consoleWriter, "%sLogger Error: Could not verify or create log directory `%s`: %s\n", logError, fullSub, err)
	}
	date := l.date.Load()
	sl.date.Store(date)
	y1, y2, m, d := unbuildDate(date)
	l.makeTodayFileName(&sl.todayFileName, sl.subDir, y1, y2, m, d)
	sl.file, err = os.OpenFile(sl.todayFileName.String(), logFileFlags, logFilePerms)
	if err != nil {
		fmt.Fprintf(l.consoleWriter, "%sLogger Error: Could not open or create log file  `%s`: %s\n", logError, sl.todayFileName.String(), err)
	}
	return sl
}

func (l *Logger) Close() error {
	close(l.ready_buffers)
	return errors.Join(l.file.Close(), l.counter.Close())
}

func (l *Logger) checkTodaysLog(sl *SubLogger, now time.Time) (y1, y2, m, d byte) {
	newDate := buildDate(now)
	date := l.date.Load()
	y1, y2, m, d = unbuildDate(date)
	var err error
	if date != newDate {
		l.file.Close()
		l.date.Store(newDate)
		l.makeTodayFileName(&l.todayFileName, l.masterDir, y1, y2, m, d)
		l.file, err = os.OpenFile(l.todayFileName.String(), logFileFlags, logFilePerms)
		if err != nil {
			fmt.Fprintf(l.consoleWriter, "%sLogger Error: Could not open or create log file  `%s`: %s\n", logError, l.todayFileName.String(), err)
		}
	}
	if sl.date.Load() != newDate {
		sl.file.Close()
		sl.date.Store(newDate)
		l.makeTodayFileName(&sl.todayFileName, sl.subDir, y1, y2, m, d)
		fmt.Printf("SubDir")
		sl.file, err = os.OpenFile(sl.todayFileName.String(), logFileFlags, logFilePerms)
		if err != nil {
			fmt.Fprintf(l.consoleWriter, "%sLogger Error: Could not open or create log file  `%s`: %s\n", logError, sl.todayFileName.String(), err)
		}
	}
	return
}

func buildDate(now time.Time) uint32 {
	y := uint32(now.Year())
	m := uint32(now.Month())
	d := uint32(now.Day())
	y2 := uint32(y % 100)
	y1 := uint32(y / 100)
	var val uint32 = (y1 << 24) | (y2 << 16) | (m << 8) | d
	return val
}

func unbuildDate(date uint32) (y1, y2, m, d byte) {
	d = byte(date)
	date >>= 8
	m = byte(date)
	date >>= 8
	y2 = byte(date)
	date >>= 8
	y1 = byte(date)
	return
}

func (l *Logger) makeTodayFileName(buf *StringBuffer, dir string, y1, y2, m, d byte) {
	buf.Reset()
	buf.WriteString(l.outDir)
	buf.WriteByte(path_sep)
	buf.WriteString(dir)
	buf.WriteByte(path_sep)
	buf.WriteString(logFilePrefix)
	buf.WriteByte(sep)
	buf.WriteString(utils.QuickItoA[y1])
	buf.WriteString(utils.QuickItoA[y2])
	buf.WriteByte(sep)
	buf.WriteString(utils.QuickItoA[m])
	buf.WriteByte(sep)
	buf.WriteString(utils.QuickItoA[d])
	buf.WriteString(fileExt)
}

func (l *Logger) log(sl *SubLogger, mode int, err error, format string, args ...any) {
	now := time.Now()
	y1, y2, m, d := l.checkTodaysLog(sl, now)
	buf := <-l.ready_buffers
	buf.Reset()
	defer func() {
		l.ready_buffers <- buf
	}()
	buf.WriteString(color[mode])
	buf.WriteString(prefix[mode])
	buf.WriteByte(open_brack)
	buf.WriteString(sl.name)
	buf.WriteByte(space)
	buf.WriteString(long[mode])
	buf.WriteByte(space)
	buf.WriteString(utils.QuickItoA[y1])
	buf.WriteString(utils.QuickItoA[y2])
	buf.WriteByte(date_sep)
	buf.WriteString(utils.QuickItoA[m])
	buf.WriteByte(date_sep)
	buf.WriteString(utils.QuickItoA[d])
	buf.WriteByte(space)
	buf.WriteString(utils.QuickItoA[now.Hour()])
	buf.WriteByte(time_sep)
	buf.WriteString(utils.QuickItoA[now.Minute()])
	buf.WriteByte(time_sep)
	buf.WriteString(utils.QuickItoA[now.Second()])
	buf.WriteByte(space)
	buf.WriteString(idPrefix)
	id := l.log_id.Add(1)
	idNew := id
	id8 := byte(id & 0b1111)
	id >>= 4
	id7 := byte(id & 0b1111)
	id >>= 4
	id6 := byte(id & 0b1111)
	id >>= 4
	id5 := byte(id & 0b1111)
	id >>= 4
	id4 := byte(id & 0b1111)
	id >>= 4
	id3 := byte(id & 0b1111)
	id >>= 4
	id2 := byte(id & 0b1111)
	id >>= 4
	id1 := byte(id)
	buf.WriteByte(utils.QuickItoX[id1])
	buf.WriteByte(utils.QuickItoX[id2])
	buf.WriteByte(utils.QuickItoX[id3])
	buf.WriteByte(utils.QuickItoX[id4])
	buf.WriteByte(utils.QuickItoX[id5])
	buf.WriteByte(utils.QuickItoX[id6])
	buf.WriteByte(utils.QuickItoX[id7])
	buf.WriteByte(utils.QuickItoX[id8])
	buf.WriteByte(close_brack)
	buf.WriteByte(space)
	fmt.Fprintf(&buf, format, args...)
	if err != nil {
		buf.WriteString(err.Error())
	}
	buf.WriteString(ansi.CLEAR)
	buf.WriteByte(newline)
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		l.consoleLock.Lock()
		l.consoleWriter.Write(buf.data)
		l.consoleLock.Unlock()
		wg.Done()
	}()
	go func() {
		l.fileLock.Lock()
		l.file.Write(buf.data[prelen[mode] : buf.Len()-ansiSuffixLen])
		l.file.Write(exNewline[:])
		l.fileLock.Unlock()
		wg.Done()
	}()
	go func() {
		sl.fileLock.Lock()
		sl.file.Write(buf.data[prelen[mode] : buf.Len()-ansiSuffixLen])
		sl.file.Write(exNewline[:])
		sl.fileLock.Unlock()
		wg.Done()
	}()
	go func(newId uint32) {
		l.counterLock.Lock()
		var cnt [4]byte
		l.counter.ReadAt(cnt[:], 0)
		oldId := binary.LittleEndian.Uint32(cnt[:])
		if newId > oldId {
			binary.LittleEndian.PutUint32(cnt[:], newId)
			l.counter.WriteAt(cnt[:], 0)
		}
		l.counterLock.Unlock()
		wg.Done()
	}(idNew)
	wg.Wait()
}

type SubLogger struct {
	logger        *Logger
	file          *os.File
	date          atomic.Uint32
	name          string
	subDir        string
	todayFileName StringBuffer
	fileLock      Mutex
}

func (l *SubLogger) Close() error {
	return l.file.Close()
}

func (l *SubLogger) Fatal(format string, args ...any) {
	l.logger.log(l, FATAL, nil, format, args...)
	os.Exit(1)
}

func (l *SubLogger) Error(format string, args ...any) {
	l.logger.log(l, ERROR, nil, format, args...)
}

func (l *SubLogger) Warn(format string, args ...any) {
	l.logger.log(l, WARN, nil, format, args...)
}

func (l *SubLogger) Note(format string, args ...any) {
	l.logger.log(l, NOTE, nil, format, args...)
}

func (l *SubLogger) Info(format string, args ...any) {
	l.logger.log(l, INFO, nil, format, args...)
}

func (l *SubLogger) Norm(format string, args ...any) {
	l.logger.log(l, NORM, nil, format, args...)
}

func (l *SubLogger) FatalIfErr(err error, format string, args ...any) {
	if err != nil {
		l.logger.log(l, FATAL, err, format, args...)
		os.Exit(1)
	}
}

func (l *SubLogger) ErrorIfErr(err error, format string, args ...any) {
	if err != nil {
		l.logger.log(l, WARN, err, format, args...)
	}
}

func (l *SubLogger) WarnIfErr(err error, format string, args ...any) {
	if err != nil {
		l.logger.log(l, WARN, err, format, args...)
	}
}

func (l *SubLogger) NoteIfErr(err error, format string, args ...any) {
	if err != nil {
		l.logger.log(l, NOTE, err, format, args...)
	}
}

func (l *SubLogger) InfoIfErr(err error, format string, args ...any) {
	if err != nil {
		l.logger.log(l, INFO, err, format, args...)
	}
}

func (l *SubLogger) NormIfErr(err error, format string, args ...any) {
	if err != nil {
		l.logger.log(l, NORM, err, format, args...)
	}
}

func (l *SubLogger) FatalIfTrue(cond bool, format string, args ...any) {
	if cond {
		l.logger.log(l, FATAL, nil, format, args...)
		os.Exit(1)
	}
}

func (l *SubLogger) ErrorIfTrue(cond bool, format string, args ...any) {
	if cond {
		l.logger.log(l, WARN, nil, format, args...)
	}
}

func (l *SubLogger) WarnIfTrue(cond bool, format string, args ...any) {
	if cond {
		l.logger.log(l, WARN, nil, format, args...)
	}
}

func (l *SubLogger) NoteIfTrue(cond bool, format string, args ...any) {
	if cond {
		l.logger.log(l, NOTE, nil, format, args...)
	}
}

func (l *SubLogger) InfoIfTrue(cond bool, format string, args ...any) {
	if cond {
		l.logger.log(l, INFO, nil, format, args...)
	}
}

func (l *SubLogger) NormIfTrue(cond bool, format string, args ...any) {
	if cond {
		l.logger.log(l, NORM, nil, format, args...)
	}
}
