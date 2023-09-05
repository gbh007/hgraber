package system

import "context"

type Logger struct {
	innerCtx context.Context
}

func NewLogger(ctx context.Context) *Logger {
	return &Logger{
		innerCtx: ctx,
	}
}

func (l *Logger) Error(err error) {
	if err == nil {
		return
	}

	print(l.innerCtx, logLevelError, 0, err.Error())
}

func (l *Logger) Info(s string) {
	print(l.innerCtx, logLevelInfo, 0, s)
}
