package hgraberworker

import (
	"app/internal/domain/hgraber"
	"context"
	"log/slog"
	"sync"
)

type hgraberUseCases interface {
	GetUnsuccessPages(ctx context.Context) []hgraber.Page
	LoadPageWithUpdate(ctx context.Context, page hgraber.Page) error

	ParseWithUpdate(ctx context.Context, book hgraber.Book)
	GetUnloadedBooks(ctx context.Context) []hgraber.Book

	ExportBook(ctx context.Context, id int) error
	ExportList(ctx context.Context) []int
}

type hasherUseCases interface {
	UnHashedPages(ctx context.Context) []hgraber.Page
	HandlePage(ctx context.Context, page hgraber.Page) error
}

type Controller struct {
	workers map[string]hgraber.WorkerStat
	mutex   *sync.RWMutex

	hasAgent bool

	hgraberUseCases hgraberUseCases
	hasherUseCases  hasherUseCases
	logger          *slog.Logger
}

func New(useCases hgraberUseCases, hasherUseCases hasherUseCases, logger *slog.Logger, hasAgent bool) *Controller {
	return &Controller{
		hgraberUseCases: useCases,
		hasherUseCases:  hasherUseCases,
		logger:          logger,

		hasAgent: hasAgent,

		workers: make(map[string]hgraber.WorkerStat),
		mutex:   new(sync.RWMutex),
	}
}
