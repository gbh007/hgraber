package system

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

const (
	logLevelDebug   = "DBG"
	logLevelInfo    = "INF"
	logLevelWarning = "WRN"
	logLevelError   = "ERR"
)

var logWriter io.Writer = os.Stderr

type LogConfig struct {
	EnableFile   bool
	AppendMode   bool
	EnableStdErr bool
}

func Init(cnf LogConfig) {
	writers := []io.Writer{}

	if cnf.EnableFile && cnf.AppendMode {
		file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Fprint(logWriter, err)
			os.Exit(1)
		}
		writers = append(writers, file)
	} else if cnf.EnableFile {
		file, err := os.Create("log.txt")
		if err != nil {
			fmt.Fprint(logWriter, err)
			os.Exit(1)
		}
		writers = append(writers, file)
	}

	if cnf.EnableStdErr {
		writers = append(writers, os.Stderr)
	}

	logWriter = io.MultiWriter(writers...)
}

func IfErr(ctx context.Context, err error) {
	if err == nil {
		return
	}

	print(ctx, logLevelError, 0, err.Error())
}

func IfErrFunc(ctx context.Context, f func() error) {
	err := f()
	if err == nil {
		return
	}

	print(ctx, logLevelError, 0, err.Error())
}

func Error(ctx context.Context, err error) {
	print(ctx, logLevelError, 0, err.Error())
}

func ErrorText(ctx context.Context, args ...interface{}) {
	print(ctx, logLevelError, 0, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	print(ctx, logLevelInfo, 0, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	print(ctx, logLevelWarning, 0, args...)
}

func print(ctx context.Context, level string, depth int, args ...interface{}) {
	fmt.Fprintf(
		logWriter,
		"[%s] [%s] %s [%s] - %s",
		level,
		GetRequestID(ctx),
		time.Now().Format(timeFormat),
		from(3+depth),
		fmt.Sprintln(args...),
	)
}

func from(depth int) string {
	from := "???"
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		if !debugMode {
			_, file = path.Split(file)
		}
		from = fmt.Sprintf("%s:%d", file, line)
	}
	return from
}
