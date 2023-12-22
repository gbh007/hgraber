package pageHandler

import (
	"app/internal/domain"
	"app/pkg/logger"
	"app/pkg/worker"
	"context"
	"io"
)

type storage interface {
	GetBook(ctx context.Context, id int) (domain.Book, error) // FIXME: не ответственность сервиса страницы, перенести в сервис книг
	GetUnsuccessPages(ctx context.Context) []domain.Page
	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error
}

type files interface {
	CreatePageFile(ctx context.Context, id, page int, ext string) (io.WriteCloser, error)
	OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error)
	CreateExportFile(ctx context.Context, name string) (io.WriteCloser, error)
}

type monitor interface {
	Register(name string, worker domain.WorkerStat)
}

type requester interface {
	RequestBytes(ctx context.Context, URL string) ([]byte, error)
}

type Service struct {
	storage   storage
	files     files
	requester requester

	worker *worker.Worker[qPage]

	logger *logger.Logger
}

type Config struct {
	Storage   storage
	Files     files
	Requester requester
	Monitor   monitor
	Logger    *logger.Logger
}

func New(cfg Config) *Service {
	s := &Service{
		storage:   cfg.Storage,
		files:     cfg.Files,
		requester: cfg.Requester,
		logger:    cfg.Logger,
	}

	s.worker = worker.New[qPage](
		queueSize,
		interval,
		cfg.Logger,
		s.handle,
		s.getter,
	)

	cfg.Monitor.Register(s.Name(), s.worker)

	return s
}
