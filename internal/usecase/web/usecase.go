package web

import (
	"log/slog"
)

type UseCase struct {
	logger *slog.Logger

	debug bool
}

func New(logger *slog.Logger, debug bool) *UseCase {
	return &UseCase{
		logger: logger,
		debug:  debug,
	}
}
