package log

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	stderr = os.Stderr
	stdout = os.Stdout
)

const (
	Debug = iota
	Info
	Warn
	Error
	Fatal
)

// https://en.wikipedia.org/wiki/ANSI_escape_code
const (
	escape = "\x1b"
	reset  = escape + "[0m"
	// Control Sequence Introducer
	csi = "[1;%dm[%-5s] " + reset
	// ANSI_escape_code
	ansi = escape + csi
)

// foreground color
const (
	b_red = iota + 31
	b_green
	b_yellow
	b_blue
	b_magenta
	b_cyan
)

const (
	ldebug level = "Debug"
	linfo  level = "Info"
	lwarn  level = "Warn"
	lerror level = "Error"
	lfatal level = "Fatal"
)

const (
	detailFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
)

var StderrLogger *log.Logger = log.New(stderr, "", detailFlag)

var Log = New()

func New() *Logger {
	return &Logger{
		debug: log.New(stderr, ldebug.colorize(b_green), detailFlag),
		info:  log.New(stderr, linfo.colorize(b_blue), detailFlag),
		warn:  log.New(stderr, lwarn.colorize(b_yellow), detailFlag),
		er:    log.New(stderr, lerror.colorize(b_red), detailFlag),
		fatal: log.New(stderr, lfatal.colorize(b_red), detailFlag),
	}
}

type level string

func (l level) String() string {
	return string(l)
}

func (l level) colorize(color int) string {
	return fmt.Sprintf(ansi, color, l.String())
}

type LevelLogger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

type Logger struct {
	mu    sync.Mutex
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	er    *log.Logger
	fatal *log.Logger
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debug.Printf(format+"\n", args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.info.Printf(format+"\n", args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.warn.Printf(format+"\n", args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.er.Printf(format+"\n", args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fatal.Printf(format+"\n", args...)
}
