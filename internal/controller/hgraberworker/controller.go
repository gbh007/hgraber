package hgraberworker

import (
	"context"
	"log/slog"
)

type WorkerUnit interface {
	Name() string
	Serve(ctx context.Context)

	InQueueCount() int
	InWorkCount() int
	RunnersCount() int
}

type Controller struct {
	workerUnits []WorkerUnit
	logger      *slog.Logger
}

func New(logger *slog.Logger, workerUnits []WorkerUnit) *Controller {
	return &Controller{
		logger:      logger,
		workerUnits: workerUnits,
	}
}
