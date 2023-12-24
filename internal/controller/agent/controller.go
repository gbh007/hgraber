package agent

import (
	"app/internal/domain/agent"
	"app/pkg/logger"
	"context"
)

type useCases interface {
	Books(ctx context.Context) []agent.BookToHandle
	BookHandle(ctx context.Context, book agent.BookToHandle)

	Pages(ctx context.Context) []agent.PageToHandle
	PageHandle(ctx context.Context, page agent.PageToHandle)
}

type Controller struct {
	logger *logger.Logger

	useCases useCases
}

func New(logger *logger.Logger, useCases useCases) *Controller {
	return &Controller{
		logger:   logger,
		useCases: useCases,
	}
}
