package system

import (
	"context"
	"io"
	"log"
)

type lWriter struct {
	ctx   context.Context
	level string
}

func (lw *lWriter) Write(data []byte) (int, error) {
	print(lw.ctx, lw.level, 2, string(data))

	return len(data), nil
}

func withLWriter(ctx context.Context, level string) io.Writer {
	return &lWriter{
		ctx:   ctx,
		level: level,
	}
}

func StdErrorLogger(ctx context.Context) *log.Logger {
	return log.New(withLWriter(ctx, logLevelError), "", 0)
}
