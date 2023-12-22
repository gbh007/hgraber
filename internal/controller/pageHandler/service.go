package pageHandler

import (
	"app/internal/domain"
	"app/pkg/logger"
	"app/pkg/worker"
	"context"
)

type useCases interface {
	GetUnsuccessPages(ctx context.Context) []domain.Page
	LoadPageWithUpdate(ctx context.Context, page domain.Page) error
}

type monitor interface {
	Register(name string, worker domain.WorkerStat)
}

type Service struct {
	useCases useCases

	worker *worker.Worker[domain.Page]

	logger *logger.Logger
}

type Config struct {
	UseCases useCases
	Monitor  monitor
	Logger   *logger.Logger
}

func New(cfg Config) *Service {
	s := &Service{
		useCases: cfg.UseCases,
		logger:   cfg.Logger,
	}

	s.worker = worker.New[domain.Page](
		queueSize,
		interval,
		cfg.Logger,
		s.handle,
		s.useCases.GetUnsuccessPages,
	)

	cfg.Monitor.Register(s.Name(), s.worker)

	return s
}
