package hgraberworker

import (
	"app/internal/domain/hgraber"
	"context"
	"log/slog"
	"sync"
)

type useCases interface {
	GetUnsuccessPages(ctx context.Context) []hgraber.Page
	LoadPageWithUpdate(ctx context.Context, page hgraber.Page) error

	ParseWithUpdate(ctx context.Context, book hgraber.Book)
	GetUnloadedBooks(ctx context.Context) []hgraber.Book

	ExportBook(ctx context.Context, id int) error
	ExportList(ctx context.Context) []int
}

type Controller struct {
	workers map[string]hgraber.WorkerStat
	mutex   *sync.RWMutex

	hasAgent bool

	useCases useCases
	logger   *slog.Logger
}

func New(useCases useCases, logger *slog.Logger, hasAgent bool) *Controller {
	return &Controller{
		useCases: useCases,
		logger:   logger,

		hasAgent: hasAgent,

		workers: make(map[string]hgraber.WorkerStat),
		mutex:   new(sync.RWMutex),
	}
}
