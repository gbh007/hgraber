package pageHandler

import (
	"app/internal/domain"
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
}

type Config struct {
	Storage   storage
	Files     files
	Requester requester
	Monitor   monitor
}

// Deprecated: устаревший подход
func Init(storage storage, files files, requester requester, monitor monitor) *Service {
	s := &Service{
		storage:   storage,
		files:     files,
		requester: requester,
	}

	s.worker = worker.New[qPage](
		queueSize,
		interval,
		s.handle,
		s.getter,
	)

	monitor.Register(s.Name(), s.worker)

	return s
}

func New(cfg Config) *Service {
	s := &Service{
		storage:   cfg.Storage,
		files:     cfg.Files,
		requester: cfg.Requester,
	}

	s.worker = worker.New[qPage](
		queueSize,
		interval,
		s.handle,
		s.getter,
	)

	cfg.Monitor.Register(s.Name(), s.worker)

	return s
}
