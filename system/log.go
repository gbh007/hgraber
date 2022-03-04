package system

import (
	"fmt"
	"io"
	"os"
	"time"
)

var logWriter io.Writer = os.Stderr

func init() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprint(logWriter, err)
		os.Exit(1)
	}
	logWriter = io.MultiWriter(os.Stderr, file)
}

const timeFormat = "2006-01-02 15:04:05"

func IfErr(ctx Context, err error) {
	if err == nil {
		return
	}
	Error(ctx, err)
}

func IfErrFunc(ctx Context, f func() error) { IfErr(ctx, f()) }

func Error(ctx Context, err error) {
	fmt.Fprintf(logWriter, "[ERR][%s][%s] - %s\n", ctx.GetRequestID(), time.Now().Format(timeFormat), err.Error())
}

func Info(ctx Context, args ...interface{}) {
	fmt.Fprintf(logWriter, "[INF][%s][%s] - %s", ctx.GetRequestID(), time.Now().Format(timeFormat), fmt.Sprintln(args...))
}
