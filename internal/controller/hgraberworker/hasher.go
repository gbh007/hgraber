package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain/hgraber"
	"context"
	"log/slog"
	"time"
)

type hashWorkerUnitUseCases interface {
	UnHashedPages(ctx context.Context) []hgraber.Page
	HandlePage(ctx context.Context, page hgraber.Page) error
}

type HashWorkerUnit struct {
	*worker.Worker[hgraber.Page]

	useCases hashWorkerUnitUseCases

	interval      time.Duration
	queueSize     int
	handlersCount int

	logger *slog.Logger
}

func NewHashWorkerUnit(useCases hashWorkerUnitUseCases, logger *slog.Logger) *HashWorkerUnit {
	w := &HashWorkerUnit{
		useCases:      useCases,
		interval:      time.Second * 15,
		queueSize:     10000,
		handlersCount: 10,
		logger:        logger,
	}

	w.Worker = worker.New[hgraber.Page](
		w.queueSize,
		w.interval,
		w.logger,
		func(ctx context.Context, page hgraber.Page) {
			err := w.useCases.HandlePage(ctx, page)
			if err != nil {
				w.logger.ErrorContext(ctx, err.Error())
			}
		},
		w.useCases.UnHashedPages,
	)

	return w
}

func (w *HashWorkerUnit) Serve(ctx context.Context) {
	w.Worker.Serve(ctx, w.handlersCount)
}

func (w *HashWorkerUnit) Name() string {
	return "hasher"
}
