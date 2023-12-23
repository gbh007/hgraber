package hgraberworker

import (
	"app/internal/domain"
	"app/pkg/logger"
	"context"
	"sync"
)

type useCases interface {
	GetUnsuccessPages(ctx context.Context) []domain.Page
	LoadPageWithUpdate(ctx context.Context, page domain.Page) error

	ParseWithUpdate(ctx context.Context, book domain.Book)
	GetUnloadedBooks(ctx context.Context) []domain.Book
}

type Controller struct {
	workers map[string]domain.WorkerStat
	mutex   *sync.RWMutex

	useCases useCases
	logger   *logger.Logger
}

func New(useCases useCases, logger *logger.Logger) *Controller {
	return &Controller{
		useCases: useCases,
		logger:   logger,

		workers: make(map[string]domain.WorkerStat),
		mutex:   new(sync.RWMutex),
	}
}
