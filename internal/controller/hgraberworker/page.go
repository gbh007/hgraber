package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain/hgraber"
	"context"
	"log/slog"
	"time"
)

type pageWorkerUnitUseCases interface {
	GetUnsuccessPages(ctx context.Context) []hgraber.Page
	LoadPageWithUpdate(ctx context.Context, page hgraber.Page) error
}

type PageWorkerUnit struct {
	*worker.Worker[hgraber.Page]

	useCases pageWorkerUnitUseCases

	interval      time.Duration
	queueSize     int
	handlersCount int

	logger *slog.Logger
}

func NewPageWorkerUnit(useCases pageWorkerUnitUseCases, logger *slog.Logger) *PageWorkerUnit {
	w := &PageWorkerUnit{
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
			err := w.useCases.LoadPageWithUpdate(ctx, page)
			if err != nil {
				w.logger.ErrorContext(ctx, err.Error())
			}
		},
		w.useCases.GetUnsuccessPages,
	)

	return w
}

func (w *PageWorkerUnit) Serve(ctx context.Context) {
	w.Worker.Serve(ctx, w.handlersCount)
}

func (w *PageWorkerUnit) Name() string {
	return "page"
}
