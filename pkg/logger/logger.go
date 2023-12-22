package logger

import (
	"app/pkg/ctxtool"
	"context"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	debug bool
}

func New(
	debug bool,
) *Logger {
	return &Logger{debug: debug}
}

func (l *Logger) WithCtx(ctx context.Context) *CtxLogger {
	return &CtxLogger{
		innerCtx: ctx,
	}
}

func (l *Logger) IfErr(ctx context.Context, err error) {
	if err == nil {
		return
	}

	l.print(ctx, logLevelError, err.Error())
}

func (l *Logger) IfErrFunc(ctx context.Context, f func() error) {
	err := f()
	if err == nil {
		return
	}

	l.print(ctx, logLevelError, err.Error())
}

func (l *Logger) Error(ctx context.Context, err error) {
	l.print(ctx, logLevelError, err.Error())
}

func (l *Logger) ErrorText(ctx context.Context, args ...interface{}) {
	l.print(ctx, logLevelError, args...)
}

func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	l.print(ctx, logLevelInfo, args...)
}

func (l *Logger) Warning(ctx context.Context, args ...interface{}) {
	l.print(ctx, logLevelWarning, args...)
}

func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	if !ctxtool.IsDebug(ctx) {
		return
	}

	l.print(ctx, logLevelDebug, args...)
}

func (l *Logger) print(ctx context.Context, level string, args ...interface{}) {
	fmt.Fprintf(
		os.Stderr,
		"[%s] [%s] %s - %s",
		level,
		ctxtool.GetRequestID(ctx),
		time.Now().Format(timeFormat),
		fmt.Sprintln(args...),
	)
}
