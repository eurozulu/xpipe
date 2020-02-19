package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var DefaultLog *Logger

func init() {
	DefaultLog = &Logger{
		Level: ERROR,
	}
}

const ERROR = LogLevel(0)
const WARN = LogLevel(1)
const INFO = LogLevel(5)
const DEBUG = LogLevel(10)

type Logger struct {
	Level LogLevel
	Out   io.Writer

}

type LogLevel int

func (log Logger) Blog(l LogLevel, format string, a ...interface{}) {
	if l > log.Level {
		return
	}

	formt := format
	out := log.Out
	// If console output, leave message, otherwise prepend header of level and time
	if out != nil {
		formt = strings.Join([]string{
			fmt.Sprintf("%v:%s:", time.Now().String(), LogLevelString(l)),
			format}, "")
	} else {
		out = os.Stdout
	}

	if !strings.HasSuffix(formt, "\n") {
		formt = strings.Join([]string{formt, "\n"}, "")
	}
	if len(a) == 0 {
		fmt.Fprint(out)
	} else {
		fmt.Fprintf(out, formt, a...)
	}
}

func Debug(format string, a ...interface{}) {
	DefaultLog.Blog(DEBUG, format, a...)
}

func Info(format string, a ...interface{}) {
	DefaultLog.Blog(INFO, format, a...)
}

func Warn(format string, a ...interface{}) {
	DefaultLog.Blog(WARN, format, a...)
}

func Error(format string, a ...interface{}) {
	DefaultLog.Blog(ERROR, format, a...)
}

// ParseLogLevel attempts to parse the given string into a know levlevel.
// accepts : error, warn, info and debug.
// is string unknown, returns 0 (ERROR)
func ParseLogLevel(s string) LogLevel {
	switch s {
	case "error", "err":
		return ERROR
	case "warn", "warning":
		return WARN
	case "info", "information":
		return INFO
	case "debug":
		return DEBUG
	default:
		return 0
	}
}

func LogLevelString(ll LogLevel) string {
	switch ll {
	case ERROR:
		return "error"
	case WARN:
		return "warn"
	case INFO:
		return "info"
	case DEBUG:
		return "debug"
	default:
		return ""
	}
}