package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain/hgraber"
	"context"
	"log/slog"
	"time"
)

type bookWorkerUnitUseCases interface {
	ParseWithUpdate(ctx context.Context, book hgraber.Book)
	GetUnloadedBooks(ctx context.Context) []hgraber.Book
}

type BookWorkerUnit struct {
	*worker.Worker[hgraber.Book]

	useCases bookWorkerUnitUseCases

	interval      time.Duration
	queueSize     int
	handlersCount int

	logger *slog.Logger
}

func NewBookWorkerUnit(useCases bookWorkerUnitUseCases, logger *slog.Logger) *BookWorkerUnit {
	w := &BookWorkerUnit{
		useCases:      useCases,
		interval:      time.Second * 15,
		queueSize:     10000,
		handlersCount: 10,
		logger:        logger,
	}

	w.Worker = worker.New[hgraber.Book](
		w.queueSize,
		w.interval,
		w.logger,
		w.useCases.ParseWithUpdate,
		w.useCases.GetUnloadedBooks,
	)

	return w
}

func (w *BookWorkerUnit) Serve(ctx context.Context) {
	w.Worker.Serve(ctx, w.handlersCount)
}

func (w *BookWorkerUnit) Name() string {
	return "book"
}
