package logger

import (
	"app/pkg/ctxtool"
	"context"
	"fmt"
	"os"
	"time"
)

type CtxLogger struct {
	innerCtx context.Context
}

func (l *CtxLogger) Error(err error) {
	if err == nil {
		return
	}

	l.print(l.innerCtx, logLevelError, err.Error())
}

func (l *CtxLogger) Info(s string) {
	l.print(l.innerCtx, logLevelInfo, s)
}

func (l *CtxLogger) print(ctx context.Context, level string, args ...interface{}) {
	fmt.Fprintf(
		os.Stderr,
		"[%s] [%s] %s - %s",
		level,
		ctxtool.GetRequestID(ctx),
		time.Now().Format(timeFormat),
		fmt.Sprintln(args...),
	)
}
