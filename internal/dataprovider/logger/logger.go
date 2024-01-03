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
	trace bool
}

func New(debug bool, trace bool) *Logger {
	return &Logger{debug: debug, trace: trace}
}

func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}

func (l *Logger) SetTrace(trace bool) {
	l.trace = trace
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

func (l *Logger) Info(ctx context.Context, args ...any) {
	l.print(ctx, logLevelInfo, args...)
}

func (l *Logger) Warning(ctx context.Context, args ...any) {
	l.print(ctx, logLevelWarning, args...)
}

func (l *Logger) Debug(ctx context.Context, args ...any) {
	if !l.debug {
		return
	}

	l.print(ctx, logLevelDebug, args...)
}

func (l *Logger) print(ctx context.Context, level string, args ...any) {
	if !l.trace {
		fmt.Fprintf(
			os.Stderr,
			"[%s] [%s] %s - %s",
			level,
			ctxtool.GetRequestID(ctx),
			time.Now().Format(timeFormat),
			fmt.Sprintln(args...),
		)

		return
	}

	fmt.Fprintf(
		os.Stderr,
		"[%s] [%s] %s - %s%s\n",
		level,
		ctxtool.GetRequestID(ctx),
		time.Now().Format(timeFormat),
		fmt.Sprintln(args...),
		simpleTrace(4, 5),
	)
}
