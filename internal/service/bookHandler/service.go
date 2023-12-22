package bookHandler

import (
	"app/internal/domain"
	"app/pkg/logger"
	"app/pkg/worker"
	"context"
)

type storage interface {
	GetUnloadedBooks(ctx context.Context) []domain.Book
	NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error)
	UpdateBookPages(ctx context.Context, id int, pages []domain.Page) error

	UpdateBookName(ctx context.Context, id int, name string) error
	UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error
}

type monitor interface {
	Register(name string, worker domain.WorkerStat)
}

type requester interface {
	RequestString(ctx context.Context, URL string) (string, error)
}

type Service struct {
	storage   storage
	requester requester

	worker *worker.Worker[domain.Book]

	logger *logger.Logger
}

type Config struct {
	Storage   storage
	Requester requester
	Monitor   monitor

	Logger *logger.Logger
}

func New(cfg Config) *Service {
	s := &Service{
		storage:   cfg.Storage,
		requester: cfg.Requester,
		logger:    cfg.Logger,
	}

	s.worker = worker.New[domain.Book](
		titleQueueSize,
		titleInterval,
		cfg.Logger,
		s.updateForWorker,
		s.storage.GetUnloadedBooks,
	)

	cfg.Monitor.Register(s.Name(), s.worker)

	return s
}
