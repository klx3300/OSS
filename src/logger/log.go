package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

var readableMu sync.Mutex

// Logger represents a logger
type Logger struct {
	LogLevel int
	Name     string
}

// Log is the default logger
var Log = Logger{
	LogLevel: 0,
	Name:     "",
}

// These are debug levels
const (
	DEBUG = iota
	INFO
	SUCC
	WARN
	FAIL
)

// Logln is the log liner
func (lg Logger) Logln(level int, levelname string, args ...interface{}) {
	if lg.LogLevel <= level {
		variParams := append([]interface{}{"[" + levelname + "]", "(" + time.Now().Format("2006/01/02 15:04:05") + ") " +
			lg.Name}, args...)
		readableMu.Lock()
		fmt.Fprintln(os.Stderr, variParams...)
		readableMu.Unlock()
	}
}

var debugColor = color.New(color.FgWhite).SprintFunc()
var succColor = color.New(color.FgGreen).SprintFunc()
var warnColor = color.New(color.FgYellow).SprintFunc()
var failColor = color.New(color.FgRed).SprintFunc()

// Debugln is the shortcut for debugging.
func (lg Logger) Debugln(args ...interface{}) {
	lg.Logln(DEBUG, debugColor("DEBG"), args...)
}

// Infoln is the shortcut
func (lg Logger) Infoln(args ...interface{}) {
	lg.Logln(INFO, "INFO", args...)
}

// Succln is the shortcut
func (lg Logger) Succln(args ...interface{}) {
	lg.Logln(SUCC, succColor("SUCC"), args...)
}

// Warnln is the shortcut
func (lg Logger) Warnln(args ...interface{}) {
	lg.Logln(WARN, warnColor("WARN"), args...)
}

// FailLn is the shortcut
func (lg Logger) FailLn(args ...interface{}) {
	lg.Logln(FAIL, failColor("FAIL"), args...)
}
