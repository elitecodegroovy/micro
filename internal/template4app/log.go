package template4app

var (
	LogFile =  `
// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/inconshreveable/log15"
)

// FileLogWriter implements LoggerInterface.
// It writes messages by lines limit, file size limit, or time frequency.
type FileLogWriter struct {
	mw *MuxWriter

	Format            log15.Format
	Filename          string
	Maxlines          int
	maxlines_curlines int

	// Rotate at size
	Maxsize         int
	maxsize_cursize int

	// Rotate daily
	Daily          bool
	Maxdays        int64
	daily_opendate int

	Rotate    bool
	startLock sync.Mutex
}

// an *os.File writer with locker.
type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

// write to os.File.
func (l *MuxWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.fd.Write(b)
}

// set os.File in writer.
func (l *MuxWriter) SetFd(fd *os.File) {
	if l.fd != nil {
		l.fd.Close()
	}
	l.fd = fd
}

// create a FileLogWriter returning as LoggerInterface.
func NewFileWriter() *FileLogWriter {
	w := &FileLogWriter{
		Filename: "",
		Format:   log15.LogfmtFormat(),
		Maxlines: 1000000,
		Maxsize:  1 << 28, //256 MB
		Daily:    true,
		Maxdays:  7,
		Rotate:   true,
	}
	// use MuxWriter instead direct use os.File for lock write when rotate
	w.mw = new(MuxWriter)
	return w
}

func (w *FileLogWriter) Log(r *log15.Record) error {
	data := w.Format.Format(r)
	w.docheck(len(data))
	_, err := w.mw.Write(data)
	return err
}

func (w *FileLogWriter) Init() error {
	if len(w.Filename) == 0 {
		return errors.New("config must have filename")
	}
	return w.StartLogger()
}

// start file logger. create log file and set to locker-inside file writer.
func (w *FileLogWriter) StartLogger() error {
	fd, err := w.createLogFile()
	if err != nil {
		return err
	}
	w.mw.SetFd(fd)
	return w.initFd()
}

func (w *FileLogWriter) docheck(size int) {
	w.startLock.Lock()
	defer w.startLock.Unlock()
	if w.Rotate && ((w.Maxlines > 0 && w.maxlines_curlines >= w.Maxlines) ||
		(w.Maxsize > 0 && w.maxsize_cursize >= w.Maxsize) ||
		(w.Daily && time.Now().Day() != w.daily_opendate)) {
		if err := w.DoRotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
			return
		}
	}
	w.maxlines_curlines++
	w.maxsize_cursize += size
}

func (w *FileLogWriter) createLogFile() (*os.File, error) {
	// Open the log file
	return os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

func (w *FileLogWriter) lineCounter() (int, error) {
	r, err := os.OpenFile(w.Filename, os.O_RDONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("lineCounter Open File : %s", err)
	}
	buf := make([]byte, 32*1024)
	count := 0

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], []byte{'\n'})
		switch {
		case err == io.EOF:
			if err := r.Close(); err != nil {
				return count, err
			}
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func (w *FileLogWriter) initFd() error {
	fd := w.mw.fd
	finfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat: %s", err)
	}
	w.maxsize_cursize = int(finfo.Size())
	w.daily_opendate = time.Now().Day()
	if finfo.Size() > 0 {
		count, err := w.lineCounter()
		if err != nil {
			return err
		}
		w.maxlines_curlines = count
	} else {
		w.maxlines_curlines = 0
	}
	return nil
}

// DoRotate means it need to write file in new file.
// new file name like xx.log.2013-01-01.2
func (w *FileLogWriter) DoRotate() error {
	_, err := os.Lstat(w.Filename)
	if err == nil { // file exists
		// Find the next available number
		num := 1
		fname := ""
		for ; err == nil && num <= 999; num++ {
			fname = w.Filename + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), num)
			_, err = os.Lstat(fname)
		}
		// return error if the last file checked still existed
		if err == nil {
			return fmt.Errorf("rotate: cannot find free log number to rename %s", w.Filename)
		}

		// block Logger's io.Writer
		w.mw.Lock()
		defer w.mw.Unlock()

		fd := w.mw.fd
		fd.Close()

		// close fd before rename
		// Rename the file to its newfound home
		if err = os.Rename(w.Filename, fname); err != nil {
			return fmt.Errorf("Rotate: %s", err)
		}

		// re-start logger
		if err = w.StartLogger(); err != nil {
			return fmt.Errorf("Rotate StartLogger: %s", err)
		}

		go w.deleteOldLog()
	}

	return nil
}

func (w *FileLogWriter) deleteOldLog() {
	dir := filepath.Dir(w.Filename)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = fmt.Errorf("Unable to delete old log '%s', error: %+v", path, r)
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.Maxdays) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(w.Filename)) {
				os.Remove(path)
			}
		}
		return returnErr
	})
}

// destroy file logger, close file writer.
func (w *FileLogWriter) Close() {
	w.mw.fd.Close()
}

// flush file logger.
// there are no buffering messages in file logger in memory.
// flush file means sync file from disk.
func (w *FileLogWriter) Flush() {
	w.mw.fd.Sync()
}

// Reload file logger
func (w *FileLogWriter) Reload() {
	// block Logger's io.Writer
	w.mw.Lock()
	defer w.mw.Unlock()

	// Close
	fd := w.mw.fd
	fd.Close()

	// Open again
	err := w.StartLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Reload StartLogger: %s\n", err)
	}
}

`
	LogFileTest = `
package log

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func (w *FileLogWriter) WriteLine(line string) error {
	n, err := w.mw.Write([]byte(line))
	if err != nil {
		return err
	}
	w.docheck(n)
	return nil
}

func TestLogFile(t *testing.T) {

	Convey("When logging to file", t, func() {
		fileLogWrite := NewFileWriter()
		So(fileLogWrite, ShouldNotBeNil)

		fileLogWrite.Filename = "grafana_test.log"
		err := fileLogWrite.Init()
		So(err, ShouldBeNil)

		Convey("Log file is empty", func() {
			So(fileLogWrite.maxlines_curlines, ShouldEqual, 0)
		})

		Convey("Logging should add lines", func() {
			err := fileLogWrite.WriteLine("test1\n")
			So(err, ShouldBeNil)
			err = fileLogWrite.WriteLine("test2\n")
			So(err, ShouldBeNil)
			err = fileLogWrite.WriteLine("test3\n")
			So(err, ShouldBeNil)
			So(fileLogWrite.maxlines_curlines, ShouldEqual, 3)
		})

		fileLogWrite.Close()
		err = os.Remove(fileLogWrite.Filename)
		So(err, ShouldBeNil)
	})
}

`
	LogHandler = `
package log

type DisposableHandler interface {
	Close()
}

type ReloadableHandler interface {
	Reload()
}

`
	LogInterface = `
package log

import "github.com/inconshreveable/log15"

type Lvl int

const (
	LvlCrit Lvl = iota
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
)

type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(ctx ...interface{}) log15.Logger

	// GetHandler gets the handler associated with the logger.
	GetHandler() log15.Handler

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h log15.Handler)

	// Log a message at the given level with context key/value pairs
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})
}

`
	LogLog = `
// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elitecodegroovy/util"
	"github.com/go-stack/stack"
	"github.com/inconshreveable/log15"
	"github.com/mattn/go-isatty"
	"gopkg.in/ini.v1"
)

var Root log15.Logger
var loggersToClose []DisposableHandler
var loggersToReload []ReloadableHandler
var filters map[string]log15.Lvl

func init() {
	loggersToClose = make([]DisposableHandler, 0)
	loggersToReload = make([]ReloadableHandler, 0)
	filters = map[string]log15.Lvl{}
	Root = log15.Root()
	Root.SetHandler(log15.DiscardHandler())
}

func New(logger string, ctx ...interface{}) Logger {
	params := append([]interface{}{"logger", logger}, ctx...)
	return Root.New(params...)
}

func Trace(format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}

	Root.Debug(message)
}

func Debug(format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}

	Root.Debug(message)
}

func Info(format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}

	Root.Info(message)
}

func Warn(format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}

	Root.Warn(message)
}

func Error(skip int, format string, v ...interface{}) {
	Root.Error(fmt.Sprintf(format, v...))
}

func Critical(skip int, format string, v ...interface{}) {
	Root.Crit(fmt.Sprintf(format, v...))
}

func Fatal(skip int, format string, v ...interface{}) {
	Root.Crit(fmt.Sprintf(format, v...))
	Close()
	os.Exit(1)
}

func Close() {
	for _, logger := range loggersToClose {
		logger.Close()
	}
	loggersToClose = make([]DisposableHandler, 0)
}

func Reload() {
	for _, logger := range loggersToReload {
		logger.Reload()
	}
}

func GetLogLevelFor(name string) Lvl {
	if level, ok := filters[name]; ok {
		switch level {
		case log15.LvlWarn:
			return LvlWarn
		case log15.LvlInfo:
			return LvlInfo
		case log15.LvlError:
			return LvlError
		case log15.LvlCrit:
			return LvlCrit
		default:
			return LvlDebug
		}
	}

	return LvlInfo
}

var logLevels = map[string]log15.Lvl{
	"trace":    log15.LvlDebug,
	"debug":    log15.LvlDebug,
	"info":     log15.LvlInfo,
	"warn":     log15.LvlWarn,
	"error":    log15.LvlError,
	"critical": log15.LvlCrit,
}

func getLogLevelFromConfig(key string, defaultName string, cfg *ini.File) (string, log15.Lvl) {
	levelName := cfg.Section(key).Key("level").MustString(defaultName)
	levelName = strings.ToLower(levelName)
	level := getLogLevelFromString(levelName)
	return levelName, level
}

func getLogLevelFromString(levelName string) log15.Lvl {
	level, ok := logLevels[levelName]

	if !ok {
		Root.Error("Unknown log level", "level", levelName)
		return log15.LvlError
	}

	return level
}

func getFilters(filterStrArray []string) map[string]log15.Lvl {
	filterMap := make(map[string]log15.Lvl)

	for _, filterStr := range filterStrArray {
		parts := strings.Split(filterStr, ":")
		if len(parts) > 1 {
			filterMap[parts[0]] = getLogLevelFromString(parts[1])
		}
	}

	return filterMap
}

func getLogFormat(format string) log15.Format {
	switch format {
	case "console":
		if isatty.IsTerminal(os.Stdout.Fd()) {
			return log15.TerminalFormat()
		}
		return log15.LogfmtFormat()
	case "text":
		return log15.LogfmtFormat()
	case "json":
		return log15.JsonFormat()
	default:
		return log15.LogfmtFormat()
	}
}

func ReadLoggingConfig(modes []string, logsPath string, cfg *ini.File) {
	Close()

	defaultLevelName, _ := getLogLevelFromConfig("log", "info", cfg)
	defaultFilters := getFilters(util.SplitString(cfg.Section("log").Key("filters").String()))

	handlers := make([]log15.Handler, 0)

	for _, mode := range modes {
		mode = strings.TrimSpace(mode)
		sec, err := cfg.GetSection("log." + mode)
		if err != nil {
			Root.Error("Unknown log mode", "mode", mode)
		}

		// Log level.
		_, level := getLogLevelFromConfig("log."+mode, defaultLevelName, cfg)
		modeFilters := getFilters(util.SplitString(sec.Key("filters").String()))
		format := getLogFormat(sec.Key("format").MustString(""))

		var handler log15.Handler

		// Generate log configuration.
		switch mode {
		case "console":
			handler = log15.StreamHandler(os.Stdout, format)
		case "file":
			fileName := sec.Key("file_name").MustString(filepath.Join(logsPath, "gnetwork.log"))
			os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
			fileHandler := NewFileWriter()
			fileHandler.Filename = fileName
			fileHandler.Format = format
			fileHandler.Rotate = sec.Key("log_rotate").MustBool(true)
			fileHandler.Maxlines = sec.Key("max_lines").MustInt(1000000)
			fileHandler.Maxsize = 1 << uint(sec.Key("max_size_shift").MustInt(28))
			fileHandler.Daily = sec.Key("daily_rotate").MustBool(true)
			fileHandler.Maxdays = sec.Key("max_days").MustInt64(7)
			fileHandler.Init()

			loggersToClose = append(loggersToClose, fileHandler)
			loggersToReload = append(loggersToReload, fileHandler)
			handler = fileHandler
		case "syslog":
			sysLogHandler := NewSyslog(sec, format)

			loggersToClose = append(loggersToClose, sysLogHandler)
			handler = sysLogHandler
		}

		for key, value := range defaultFilters {
			if _, exist := modeFilters[key]; !exist {
				modeFilters[key] = value
			}
		}

		for key, value := range modeFilters {
			if _, exist := filters[key]; !exist {
				filters[key] = value
			}
		}

		handler = LogFilterHandler(level, modeFilters, handler)
		handlers = append(handlers, handler)
	}

	Root.SetHandler(log15.MultiHandler(handlers...))
}

func LogFilterHandler(maxLevel log15.Lvl, filters map[string]log15.Lvl, h log15.Handler) log15.Handler {
	return log15.FilterHandler(func(r *log15.Record) (pass bool) {

		if len(filters) > 0 {
			for i := 0; i < len(r.Ctx); i += 2 {
				key, ok := r.Ctx[i].(string)
				if ok && key == "logger" {
					loggerName, strOk := r.Ctx[i+1].(string)
					if strOk {
						if filterLevel, ok := filters[loggerName]; ok {
							return r.Lvl <= filterLevel
						}
					}
				}
			}
		}

		return r.Lvl <= maxLevel
	}, h)
}

func Stack(skip int) string {
	call := stack.Caller(skip)
	s := stack.Trace().TrimBelow(call).TrimRuntime()
	return s.String()
}

`
	LogLogWriter = `
package log

import (
	"io"
	"strings"
)

type logWriterImpl struct {
	log    Logger
	level  Lvl
	prefix string
}

func NewLogWriter(log Logger, level Lvl, prefix string) io.Writer {
	return &logWriterImpl{
		log:    log,
		level:  level,
		prefix: prefix,
	}
}

func (l *logWriterImpl) Write(p []byte) (n int, err error) {
	message := l.prefix + strings.TrimSpace(string(p))

	switch l.level {
	case LvlCrit:
		l.log.Crit(message)
	case LvlError:
		l.log.Error(message)
	case LvlWarn:
		l.log.Warn(message)
	case LvlInfo:
		l.log.Info(message)
	default:
		l.log.Debug(message)
	}

	return len(p), nil
}
`
	LogLogWriterTest = `
package log

import (
	"testing"

	"github.com/inconshreveable/log15"
	. "github.com/smartystreets/goconvey/convey"
)

type FakeLogger struct {
	debug string
	info  string
	warn  string
	err   string
	crit  string
}

func (f *FakeLogger) New(ctx ...interface{}) log15.Logger {
	return nil
}

func (f *FakeLogger) Debug(msg string, ctx ...interface{}) {
	f.debug = msg
}

func (f *FakeLogger) Info(msg string, ctx ...interface{}) {
	f.info = msg
}

func (f *FakeLogger) Warn(msg string, ctx ...interface{}) {
	f.warn = msg
}

func (f *FakeLogger) Error(msg string, ctx ...interface{}) {
	f.err = msg
}

func (f *FakeLogger) Crit(msg string, ctx ...interface{}) {
	f.crit = msg
}

func (f *FakeLogger) GetHandler() log15.Handler {
	return nil
}

func (f *FakeLogger) SetHandler(l log15.Handler) {}

func TestLogWriter(t *testing.T) {
	Convey("When writing to a LogWriter", t, func() {
		Convey("Should write using the correct level [crit]", func() {
			fake := &FakeLogger{}

			crit := NewLogWriter(fake, LvlCrit, "")
			n, err := crit.Write([]byte("crit"))

			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(fake.crit, ShouldEqual, "crit")
		})

		Convey("Should write using the correct level [error]", func() {
			fake := &FakeLogger{}

			crit := NewLogWriter(fake, LvlError, "")
			n, err := crit.Write([]byte("error"))

			So(n, ShouldEqual, 5)
			So(err, ShouldBeNil)
			So(fake.err, ShouldEqual, "error")
		})

		Convey("Should write using the correct level [warn]", func() {
			fake := &FakeLogger{}

			crit := NewLogWriter(fake, LvlWarn, "")
			n, err := crit.Write([]byte("warn"))

			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(fake.warn, ShouldEqual, "warn")
		})

		Convey("Should write using the correct level [info]", func() {
			fake := &FakeLogger{}

			crit := NewLogWriter(fake, LvlInfo, "")
			n, err := crit.Write([]byte("info"))

			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(fake.info, ShouldEqual, "info")
		})

		Convey("Should write using the correct level [debug]", func() {
			fake := &FakeLogger{}

			crit := NewLogWriter(fake, LvlDebug, "")
			n, err := crit.Write([]byte("debug"))

			So(n, ShouldEqual, 5)
			So(err, ShouldBeNil)
			So(fake.debug, ShouldEqual, "debug")
		})

		Convey("Should prefix the output with the prefix", func() {
			fake := &FakeLogger{}

			crit := NewLogWriter(fake, LvlDebug, "prefix")
			n, err := crit.Write([]byte("debug"))

			So(n, ShouldEqual, 5) // n is how much of input consumed
			So(err, ShouldBeNil)
			So(fake.debug, ShouldEqual, "prefixdebug")
		})
	})
}

`
	LogSyslog = `
//+build !windows,!nacl,!plan9

package log

import (
	"errors"
	"log/syslog"
	"os"

	"github.com/inconshreveable/log15"
	"gopkg.in/ini.v1"
)

type SysLogHandler struct {
	syslog   *syslog.Writer
	Network  string
	Address  string
	Facility string
	Tag      string
	Format   log15.Format
}

func NewSyslog(sec *ini.Section, format log15.Format) *SysLogHandler {
	handler := &SysLogHandler{
		Format: log15.LogfmtFormat(),
	}

	handler.Format = format
	handler.Network = sec.Key("network").MustString("")
	handler.Address = sec.Key("address").MustString("")
	handler.Facility = sec.Key("facility").MustString("local7")
	handler.Tag = sec.Key("tag").MustString("")

	if err := handler.Init(); err != nil {
		Root.Error("Failed to init syslog log handler", "error", err)
		os.Exit(1)
	}

	return handler
}

func (sw *SysLogHandler) Init() error {
	prio, err := parseFacility(sw.Facility)
	if err != nil {
		return err
	}

	w, err := syslog.Dial(sw.Network, sw.Address, prio, sw.Tag)
	if err != nil {
		return err
	}

	sw.syslog = w
	return nil
}

func (sw *SysLogHandler) Log(r *log15.Record) error {
	var err error

	msg := string(sw.Format.Format(r))

	switch r.Lvl {
	case log15.LvlDebug:
		err = sw.syslog.Debug(msg)
	case log15.LvlInfo:
		err = sw.syslog.Info(msg)
	case log15.LvlWarn:
		err = sw.syslog.Warning(msg)
	case log15.LvlError:
		err = sw.syslog.Err(msg)
	case log15.LvlCrit:
		err = sw.syslog.Crit(msg)
	default:
		err = errors.New("invalid syslog level")
	}

	return err
}

func (sw *SysLogHandler) Close() {
	sw.syslog.Close()
}

var facilities = map[string]syslog.Priority{
	"user":   syslog.LOG_USER,
	"daemon": syslog.LOG_DAEMON,
	"local0": syslog.LOG_LOCAL0,
	"local1": syslog.LOG_LOCAL1,
	"local2": syslog.LOG_LOCAL2,
	"local3": syslog.LOG_LOCAL3,
	"local4": syslog.LOG_LOCAL4,
	"local5": syslog.LOG_LOCAL5,
	"local6": syslog.LOG_LOCAL6,
	"local7": syslog.LOG_LOCAL7,
}

func parseFacility(facility string) (syslog.Priority, error) {
	prio, ok := facilities[facility]
	if !ok {
		return syslog.LOG_LOCAL0, errors.New("invalid syslog facility")
	}

	return prio, nil
}
`
	LogSyslogWindows = `
//+build windows

package log

import (
	"github.com/inconshreveable/log15"
	"gopkg.in/ini.v1"
)

type SysLogHandler struct {
}

func NewSyslog(sec *ini.Section, format log15.Format) *SysLogHandler {
	return &SysLogHandler{}
}

func (sw *SysLogHandler) Log(r *log15.Record) error {
	return nil
}

func (sw *SysLogHandler) Close() {
}

`
)
