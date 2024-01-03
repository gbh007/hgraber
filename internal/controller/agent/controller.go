package agent

import (
	"app/internal/domain/agent"
	"context"
)

type logger interface {
	Debug(ctx context.Context, args ...any)
	Info(ctx context.Context, args ...any)
}

type useCases interface {
	Books(ctx context.Context) []agent.BookToHandle
	BookHandle(ctx context.Context, book agent.BookToHandle)

	Pages(ctx context.Context) []agent.PageToHandle
	PageHandle(ctx context.Context, page agent.PageToHandle)
}

type Controller struct {
	logger logger

	useCases useCases
}

func New(logger logger, useCases useCases) *Controller {
	return &Controller{
		logger:   logger,
		useCases: useCases,
	}
}
