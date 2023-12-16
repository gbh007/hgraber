package pageHandler

import (
	"app/internal/domain"
	"app/pkg/worker"
	"context"
	"io"
)

type storage interface {
	GetBook(ctx context.Context, id int) (domain.Book, error)
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
