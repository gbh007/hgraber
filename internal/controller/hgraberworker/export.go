package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"context"
	"log/slog"
	"time"
)

type exportWorkerUnitUseCases interface {
	ExportBook(ctx context.Context, id int) error
	ExportList(ctx context.Context) []int
}

type ExportWorkerUnit struct {
	*worker.Worker[int]

	useCases exportWorkerUnitUseCases

	interval      time.Duration
	queueSize     int
	handlersCount int

	logger *slog.Logger
}

func NewExportWorkerUnit(useCases exportWorkerUnitUseCases, logger *slog.Logger) *ExportWorkerUnit {
	w := &ExportWorkerUnit{
		useCases:      useCases,
		interval:      time.Second * 5,
		queueSize:     1000,
		handlersCount: 3,
		logger:        logger,
	}

	w.Worker = worker.New[int](
		w.queueSize,
		w.interval,
		w.logger,
		func(ctx context.Context, bookID int) {
			err := w.useCases.ExportBook(ctx, bookID)
			if err != nil {
				w.logger.ErrorContext(ctx, err.Error())
			}
		},
		w.useCases.ExportList,
	)

	return w
}

func (w *ExportWorkerUnit) Serve(ctx context.Context) {
	w.Worker.Serve(ctx, w.handlersCount)
}

func (w *ExportWorkerUnit) Name() string {
	return "export"
}
