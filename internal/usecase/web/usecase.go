package web

import "context"

type logger interface {
	Debug(ctx context.Context, args ...any)
	Error(ctx context.Context, err error)
	IfErr(ctx context.Context, err error)
	Warning(ctx context.Context, args ...any)
}

type UseCase struct {
	logger logger

	debug bool
}

func New(logger logger, debug bool) *UseCase {
	return &UseCase{
		logger: logger,
		debug:  debug,
	}
}
