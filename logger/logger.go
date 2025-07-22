package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gabe-lee/OurSweeper/ansi"
	"github.com/gabe-lee/OurSweeper/data_buffer"
	"github.com/gabe-lee/OurSweeper/lock"
	"github.com/gabe-lee/OurSweeper/utils"
	"github.com/gabe-lee/OurSweeper/wire"
)

type (
	Mutex           = sync.Mutex
	Writer          = io.Writer
	WriteCloser     = io.WriteCloser
	Builder         = strings.Builder
	Buffer          = bytes.Buffer
	Stringer        = fmt.Stringer
	Time            = time.Time
	StringBuffer    = data_buffer.WriteBuffer
	ErrorBuffer     = utils.ErrorBuffer
	MiniLock        = lock.MiniLock
	ReadWriteSeeker = io.ReadWriteSeeker
	ReadWriteWire   = wire.ReadWriteWire
)

const (
	LOG int = iota
	INFO
	NOTE
	WARN
	ERROR
	FATAL
	levelCount
)

const (
	logFilePrefix string = "Log"
	fileExt       string = ".txt"
	fileMaxSize   int    = 1 << 14
	logError      string = "<!!>"
	// idPrefix        string = "0x"
	idPrefix        byte   = byte('#')
	initCap         int    = 128
	sep             byte   = byte('_')
	path_sep        byte   = byte('/')
	date_sep        byte   = byte('/')
	time_sep        byte   = byte(':')
	open_brack      byte   = byte('[')
	space           byte   = byte(' ')
	colon           byte   = byte(':')
	close_brack     byte   = byte(']')
	newline         byte   = byte('\n')
	logFileFlags           = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	metaFileFlags          = os.O_RDWR | os.O_CREATE
	logFilePerms           = 0777 | os.ModeAppend
	metaFilePerms          = 0777
	logFileDirPerms        = 0777
	ansiPrefixLen   int    = len(ansi.FG_BLK)
	ansiSuffixLen   int    = len(ansi.CLEAR) + 1
	maxNameLen      int    = 14
	nameFill        string = ". . . . . . . ."
	counterName     string = ".metadata"
)

// Metadata file var offests
const (
	offLogId = 0
)

var bin = wire.LE

var exNewline = [1]byte{newline}

var prelen = [levelCount]int{
	LOG:   len(color[LOG]),
	INFO:  len(color[INFO]),
	NOTE:  len(color[NOTE]),
	WARN:  len(color[WARN]),
	ERROR: len(color[ERROR]),
	FATAL: len(color[FATAL]),
}

var prefix = [levelCount]string{
	LOG:   "--",
	INFO:  "~~",
	NOTE:  "==",
	WARN:  "**",
	ERROR: "XX",
	FATAL: "##",
}

var long = [levelCount]string{
	LOG:   "LOG  ",
	INFO:  "INFO ",
	NOTE:  "NOTE ",
	WARN:  "WARN ",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

var color = [levelCount]string{
	LOG:   "",
	INFO:  ansi.FG_BLU,
	NOTE:  ansi.FG_GRN,
	WARN:  ansi.FG_YEL,
	ERROR: ansi.FG_RED,
	FATAL: ansi.FG_MAG,
}

type Logger struct {
	outDir        string
	masterDir     string
	masterLine    string
	date          atomic.Uint32
	log_id        uint64 // locked by metadataLock
	todayFileName StringBuffer
	consoleWriter Writer
	file          *os.File
	metaFile      *os.File
	metadataLock  MiniLock
	fileLock      MiniLock
	consoleLock   MiniLock
	ready_buffers chan StringBuffer
	errs          ErrorBuffer
}

var metadataInit = [...]byte{0, 0, 0, 0, 0, 0, 0, 0}

func NewLogger(outputDir string, masterDir string, consoleWriter Writer, bufferPool int) Logger {
	cwd, _ := os.Getwd()
	n := []byte(nameFill)
	nLen := min(len(masterDir), maxNameLen)
	copy(n, masterDir[:nLen])
	l := Logger{
		outDir:        path.Join(cwd, outputDir),
		masterDir:     masterDir,
		masterLine:    string(n),
		consoleWriter: consoleWriter,
		ready_buffers: make(chan StringBuffer, bufferPool),
		todayFileName: data_buffer.NewWriteBuffer(initCap),
		errs:          utils.NewErrorBuffer(consoleWriter, initCap),
	}
	defer l.errs.Flush("logger error .NewLogger(): ")
	for range bufferPool {
		l.ready_buffers <- data_buffer.NewWriteBuffer(initCap)
	}
	fullMaster := path.Join(l.outDir, masterDir)
	err := os.MkdirAll(fullMaster, logFileDirPerms)
	l.errs.IfErrAddErrWithStr(err, "could not verify or create log directory `%s`", fullMaster)
	metaFile := path.Join(l.outDir, counterName)
	l.metaFile, err = os.OpenFile(metaFile, metaFileFlags, metaFilePerms)
	l.errs.IfErrAddErrWithStr(err, "could not open or create metadata file `%s`", metaFile)
	metaStat, _ := l.metaFile.Stat()
	if metaStat.Size() == 0 {
		_, err = l.metaFile.WriteAt(metadataInit[:], 0)
		l.errs.IfErrAddErrWithStr(err, "unable to initialize metadata file")
	}
	var arr [8]byte
	_, err = l.metaFile.ReadAt(arr[:], offLogId)
	l.errs.IfErrAddErrWithStr(err, "unable to read log id counter from `%s` file", metaFile)
	bin.ReadU64(arr, &l.log_id)
	now := time.Now()
	date := buildDate(now)
	y1, y2, m, d := unbuildDate(date)
	l.date.Store(date)
	l.makeTodayFileName(&l.todayFileName, masterDir, y1, y2, m, d)
	l.file, err = os.OpenFile(l.todayFileName.StringRef(), logFileFlags, logFilePerms)
	l.errs.IfErrAddErrWithStr(err, "could not open or create log file `%s`", l.todayFileName.StringRef())
	return l
}

func (l *Logger) NewSubLogger(name string) SubLogger {
	defer l.errs.Flush("logger error .NewSubLogger(): ")
	n := []byte(nameFill)
	nLen := min(len(name), maxNameLen)
	copy(n, name[:nLen])
	sl := SubLogger{
		logger:        l,
		name:          string(n),
		subDir:        name,
		todayFileName: data_buffer.NewWriteBuffer(initCap),
	}
	fullSub := path.Join(l.outDir, sl.subDir)
	err := os.MkdirAll(fullSub, logFileDirPerms)
	l.errs.IfErrAddErrWithStr(err, "could not verify or create log directory `%s`", fullSub)
	date := l.date.Load()
	sl.date.Store(date)
	y1, y2, m, d := unbuildDate(date)
	l.makeTodayFileName(&sl.todayFileName, sl.subDir, y1, y2, m, d)
	sl.file, err = os.OpenFile(sl.todayFileName.StringRef(), logFileFlags, logFilePerms)
	l.errs.IfErrAddErrWithStr(err, "could not open or create log file  `%s`", sl.todayFileName.StringRef())
	return sl
}

func (l *Logger) Close() {
	defer l.errs.Flush("error Logger.Close(): ")
	l.errs.IfErrAddErr(l.file.Close())
	l.errs.IfErrAddErr(l.metaFile.Close())
	close(l.ready_buffers)
}

func (l *Logger) Fatal(format string, args ...any) {
	defer l.errs.Flush("error Logger.Fatal(): ")
	l.log(nil, FATAL, nil, format, args...)
	os.Exit(1)
}

func (l *Logger) Error(format string, args ...any) {
	defer l.errs.Flush("error Logger.Error(): ")
	l.log(nil, ERROR, nil, format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	defer l.errs.Flush("error Logger.Warn(): ")
	l.log(nil, WARN, nil, format, args...)
}

func (l *Logger) Note(format string, args ...any) {
	defer l.errs.Flush("error Logger.Note(): ")
	l.log(nil, NOTE, nil, format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	defer l.errs.Flush("error Logger.Info(): ")
	l.log(nil, INFO, nil, format, args...)
}

func (l *Logger) Norm(format string, args ...any) {
	defer l.errs.Flush("error Logger.Norm(): ")
	l.log(nil, LOG, nil, format, args...)
}

func (l *Logger) FatalIfErr(err error, format string, args ...any) {
	defer l.errs.Flush("error Logger.FatalIfErr(): ")
	if err != nil {
		l.log(nil, FATAL, err, format, args...)
		os.Exit(1)
	}
}

func (l *Logger) ErrorIfErr(err error, format string, args ...any) {
	defer l.errs.Flush("error Logger.ErrorIfErr(): ")
	if err != nil {
		l.log(nil, WARN, err, format, args...)
	}
}

func (l *Logger) WarnIfErr(err error, format string, args ...any) {
	defer l.errs.Flush("error Logger.WarnIfErr(): ")
	if err != nil {
		l.log(nil, WARN, err, format, args...)
	}
}

func (l *Logger) NoteIfErr(err error, format string, args ...any) {
	defer l.errs.Flush("error Logger.NoteIfErr(): ")
	if err != nil {
		l.log(nil, NOTE, err, format, args...)
	}
}

func (l *Logger) InfoIfErr(err error, format string, args ...any) {
	defer l.errs.Flush("error Logger.InfoIfErr(): ")
	if err != nil {
		l.log(nil, INFO, err, format, args...)
	}
}

func (l *Logger) NormIfErr(err error, format string, args ...any) {
	defer l.errs.Flush("error Logger.NormIfErr(): ")
	if err != nil {
		l.log(nil, LOG, err, format, args...)
	}
}

func (l *Logger) FatalIfTrue(cond bool, format string, args ...any) {
	defer l.errs.Flush("error Logger.FatalIfTrue(): ")
	if cond {
		l.log(nil, FATAL, nil, format, args...)
		os.Exit(1)
	}
}

func (l *Logger) ErrorIfTrue(cond bool, format string, args ...any) {
	defer l.errs.Flush("error Logger.ErrorIfTrue(): ")
	if cond {
		l.log(nil, WARN, nil, format, args...)
	}
}

func (l *Logger) WarnIfTrue(cond bool, format string, args ...any) {
	defer l.errs.Flush("error Logger.WarnIfTrue(): ")
	if cond {
		l.log(nil, WARN, nil, format, args...)
	}
}

func (l *Logger) NoteIfTrue(cond bool, format string, args ...any) {
	defer l.errs.Flush("error Logger.NoteIfTrue(): ")
	if cond {
		l.log(nil, NOTE, nil, format, args...)
	}
}

func (l *Logger) InfoIfTrue(cond bool, format string, args ...any) {
	defer l.errs.Flush("error Logger.InfoIfTrue(): ")
	if cond {
		l.log(nil, INFO, nil, format, args...)
	}
}

func (l *Logger) NormIfTrue(cond bool, format string, args ...any) {
	defer l.errs.Flush("error Logger.NormIfTrue(): ")
	if cond {
		l.log(nil, LOG, nil, format, args...)
	}
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
		l.file, err = os.OpenFile(l.todayFileName.StringRef(), logFileFlags, logFilePerms)
		l.errs.IfErrAddErrWithStr(err, "could not open or create log file  `%s`", l.todayFileName.StringRef())
	}
	if sl != nil {
		if sl.date.Load() != newDate {
			sl.file.Close()
			sl.date.Store(newDate)
			l.makeTodayFileName(&sl.todayFileName, sl.subDir, y1, y2, m, d)
			sl.file, err = os.OpenFile(sl.todayFileName.StringRef(), logFileFlags, logFilePerms)
			l.errs.IfErrAddErrWithStr(err, "could not open or create log file  `%s`", sl.todayFileName.StringRef())
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

func (l *Logger) log(sl *SubLogger, mode int, err error, format string, args ...any) (logId uint64) {
	now := time.Now()
	y1, y2, m, d := l.checkTodaysLog(sl, now)
	buf := <-l.ready_buffers
	buf.Reset()
	defer func() {
		l.ready_buffers <- buf
	}()
	name := l.masterLine
	if sl != nil {
		name = sl.name
	}
	buf.WriteString(color[mode])
	buf.WriteString(prefix[mode])
	buf.WriteByte(open_brack)
	buf.WriteString(name)
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
	buf.WriteByte(idPrefix)
	id := atomic.AddUint64(&l.log_id, 1)
	idBytes := utils.QuickIntToHexString(id)
	buf.WriteBytes(idBytes[:]...)
	buf.WriteByte(close_brack)
	buf.WriteByte(space)
	fmt.Fprintf(&buf, format, args...)
	if err != nil {
		buf.WriteBytes(colon, space)
		buf.WriteString(err.Error())
	}
	buf.TrimEndWhitespace()
	buf.WriteString(ansi.CLEAR)
	buf.WriteByte(newline)
	wg := sync.WaitGroup{}
	wg.Add(3)
	if sl != nil {
		wg.Add(1)
		go func() {
			sl.fileLock.Lock()
			_, err := sl.file.Write(buf.BytesRef()[prelen[mode] : buf.Len()-ansiSuffixLen])
			l.WarnIfErr(err, "failed to write %s log", sl.name)
			_, err = sl.file.Write(exNewline[:])
			l.WarnIfErr(err, "failed to write %s log newline", sl.name)
			sl.fileLock.Unlock()
			wg.Done()
		}()
	}
	go func() {
		l.consoleLock.Lock()
		l.consoleWriter.Write(buf.BytesRef())
		l.consoleLock.Unlock()
		wg.Done()
	}()
	go func() {
		l.fileLock.Lock()
		_, err := l.file.Write(buf.BytesRef()[prelen[mode] : buf.Len()-ansiSuffixLen])
		l.WarnIfErr(err, "failed to write %s log", l.masterDir)
		_, err = l.file.Write(exNewline[:])
		l.WarnIfErr(err, "failed to write %s log newline", l.masterDir)
		l.fileLock.Unlock()
		wg.Done()
	}()
	go func(newId uint64) {
		l.metadataLock.Lock()
		var arr [8]byte
		_, err = l.metaFile.ReadAt(arr[:], offLogId)
		var oldId uint64
		bin.ReadU64(arr, &oldId)
		l.WarnIfErr(err, "failed to read metadata file value `id` at offset %d", offLogId)
		if newId > oldId {
			bin.WriteU64(newId, &arr)
			_, err = l.metaFile.WriteAt(arr[:], offLogId)
			l.WarnIfErr(err, "failed to write metadata file value `id` at offset %d", offLogId)
		}
		l.metadataLock.Unlock()
		wg.Done()
	}(id)
	wg.Wait()
	return id
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

func (sl *SubLogger) Close() {
	defer sl.logger.errs.Flush("error SubLogger.Close(): ")
	sl.logger.errs.IfErrAddErr(sl.file.Close())
}

func (sl *SubLogger) NewSubLoggerWriter(mode int) SubLoggerWriter {
	slw := SubLoggerWriter{
		sl: sl,
	}
	slw.SetMode(mode)
	return slw
}

func (sl *SubLogger) Fatal(format string, args ...any) {
	defer sl.logger.errs.Flush("error SubLogger.Fatal(): ")
	sl.logger.log(sl, FATAL, nil, format, args...)
	os.Exit(1)
}

func (sl *SubLogger) Error(format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.Error(): ")
	sl.logger.log(sl, ERROR, nil, format, args...)
	return
}

func (sl *SubLogger) Warn(format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.Warn(): ")
	sl.logger.log(sl, WARN, nil, format, args...)
	return
}

func (sl *SubLogger) Note(format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.Note(): ")
	sl.logger.log(sl, NOTE, nil, format, args...)
	return
}

func (sl *SubLogger) Info(format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.Info(): ")
	sl.logger.log(sl, INFO, nil, format, args...)
	return
}

func (sl *SubLogger) Norm(format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.Norm(): ")
	sl.logger.log(sl, LOG, nil, format, args...)
	return
}

func (sl *SubLogger) Log(mode int, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.Norm(): ")
	sl.logger.log(sl, mode, nil, format, args...)
	if mode == FATAL {
		os.Exit(1)
	}
	return
}

func (sl *SubLogger) FatalIfErr(err error, format string, args ...any) {
	defer sl.logger.errs.Flush("error SubLogger.FatalIfErr(): ")
	if err != nil {
		sl.logger.log(sl, FATAL, err, format, args...)
		os.Exit(1)
	}
}

func (sl *SubLogger) ErrorIfErr(err error, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.ErrorIfErr(): ")
	if err != nil {
		sl.logger.log(sl, WARN, err, format, args...)
	}
	return
}

func (sl *SubLogger) WarnIfErr(err error, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.WarnIfErr(): ")
	if err != nil {
		sl.logger.log(sl, WARN, err, format, args...)
	}
	return
}

func (sl *SubLogger) NoteIfErr(err error, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.NoteIfErr(): ")
	if err != nil {
		sl.logger.log(sl, NOTE, err, format, args...)
	}
	return
}

func (sl *SubLogger) InfoIfErr(err error, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.InfoIfErr(): ")
	if err != nil {
		sl.logger.log(sl, INFO, err, format, args...)
	}
	return
}

func (sl *SubLogger) NormIfErr(err error, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.NormIfErr(): ")
	if err != nil {
		sl.logger.log(sl, LOG, err, format, args...)
	}
	return
}

func (sl *SubLogger) LogIfErr(mode int, err error, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.NormIfErr(): ")
	if err != nil {
		sl.logger.log(sl, mode, err, format, args...)
	}
	return
}

func (sl *SubLogger) FatalIfTrue(cond bool, format string, args ...any) {
	defer sl.logger.errs.Flush("error SubLogger.FatalIfTrue(): ")
	if cond {
		sl.logger.log(sl, FATAL, nil, format, args...)
		os.Exit(1)
	}
}

func (sl *SubLogger) ErrorIfTrue(cond bool, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.ErrorIfTrue(): ")
	if cond {
		sl.logger.log(sl, WARN, nil, format, args...)
	}
	return
}

func (sl *SubLogger) WarnIfTrue(cond bool, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.WarnIfTrue(): ")
	if cond {
		sl.logger.log(sl, WARN, nil, format, args...)
	}
	return
}

func (sl *SubLogger) NoteIfTrue(cond bool, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.NoteIfTrue(): ")
	if cond {
		sl.logger.log(sl, NOTE, nil, format, args...)
	}
	return
}

func (sl *SubLogger) InfoIfTrue(cond bool, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.InfoIfTrue(): ")
	if cond {
		sl.logger.log(sl, INFO, nil, format, args...)
	}
	return
}

func (sl *SubLogger) NormIfTrue(cond bool, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.NormIfTrue(): ")
	if cond {
		sl.logger.log(sl, LOG, nil, format, args...)
	}
	return
}

func (sl *SubLogger) LogIfTrue(mode int, cond bool, format string, args ...any) (logId uint64) {
	defer sl.logger.errs.Flush("error SubLogger.NormIfTrue(): ")
	if cond {
		sl.logger.log(sl, mode, nil, format, args...)
	}
	return
}
