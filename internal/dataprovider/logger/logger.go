package logger

import (
	"app/internal/dataprovider/slogHandler"
	"log/slog"
)

func New(debug bool, trace bool) *slog.Logger {
	opts := []slogHandler.HandlerOption{}

	if debug {
		opts = append(opts, slogHandler.WithDebug())
	}

	if trace {
		opts = append(opts, slogHandler.WithPrinter(printWithTrace))
	} else {
		opts = append(opts, slogHandler.WithPrinter(print))
	}

	return slog.New(slogHandler.New(opts...))
}
