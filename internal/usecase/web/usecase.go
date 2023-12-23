package web

import (
	"app/pkg/logger"
)

type UseCase struct {
	logger *logger.Logger

	debug bool
}

func New(logger *logger.Logger, debug bool) *UseCase {
	return &UseCase{
		logger: logger,
		debug:  debug,
	}
}
