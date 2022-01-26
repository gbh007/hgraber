package clog

import (
	"app/system/coreContext"
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

func Error(ctx coreContext.CoreContext, err error) {
	fmt.Fprintf(logWriter, "[%s]{%s} - %s\n", ctx.GetRequestID(), time.Now().Format(timeFormat), err.Error())
}
func Info(ctx coreContext.CoreContext, args ...interface{}) {
	fmt.Fprintf(logWriter, "[%s]{%s} - %s", ctx.GetRequestID(), time.Now().Format(timeFormat), fmt.Sprintln(args...))
}
