package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

var DefaultLog Logger

const ErrorText = "ERROR:"

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

	out := log.Out
	if out == nil {
		out = os.Stdout
	}
	if !strings.HasSuffix(format, "\n") {
		format = strings.Join([]string{format, "\n"}, "")
	}
	fmt.Fprintf(out, format, a)
}

func (log Logger) Error(format string, err ...error) {
	buf := bytes.NewBuffer(nil)
	for _, e := range err {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(e.Error())
	}

	out := log.Out
	if out != nil {
		// For file output, preceed with error text.
		fmt.Fprint(out, ErrorText)
	} else {
		out = os.Stderr
	}

	fmt.Fprintf(out, format, buf.String())
}

func Log(l LogLevel, format string,  m ...interface{}) {
	DefaultLog.Blog(l, format, m...)
}

func Error(format string, err ...error) {
	DefaultLog.Error(format, err...)
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