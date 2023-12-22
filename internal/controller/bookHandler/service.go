package bookHandler

import (
	"app/internal/domain"
	"app/pkg/logger"
	"app/pkg/worker"
	"context"
)

type monitor interface {
	Register(name string, worker domain.WorkerStat)
}

type useCases interface {
	ParseWithUpdate(ctx context.Context, title domain.Book)
	GetUnloadedBooks(ctx context.Context) []domain.Book
}

type Service struct {
	useCases useCases

	worker *worker.Worker[domain.Book]

	logger *logger.Logger
}

type Config struct {
	UseCases useCases
	Monitor  monitor

	Logger *logger.Logger
}

func New(cfg Config) *Service {
	s := &Service{
		useCases: cfg.UseCases,
		logger:   cfg.Logger,
	}

	s.worker = worker.New[domain.Book](
		titleQueueSize,
		titleInterval,
		cfg.Logger,
		s.useCases.ParseWithUpdate,
		s.useCases.GetUnloadedBooks,
	)

	cfg.Monitor.Register(s.Name(), s.worker)

	return s
}
