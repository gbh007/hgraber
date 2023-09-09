package pageHandler

import (
	"app/internal/domain"
	"app/pkg/worker"
	"context"
	"io"
)

type storage interface {
	GetBook(ctx context.Context, id int) (domain.Book, error)
	GetUnsuccessPages(ctx context.Context) []domain.PageFullInfo
	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error
}

type files interface {
	CreatePageFile(ctx context.Context, id, page int, ext string) (io.WriteCloser, error)
	OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error)
	CreateExportFile(ctx context.Context, name string) (io.WriteCloser, error)
}

type Service struct {
	storage storage
	files   files

	worker *worker.Worker[qPage]
}

func Init(storage storage, files files) *Service {
	s := &Service{
		storage: storage,
		files:   files,
	}

	s.worker = worker.New[qPage](
		queueSize,
		interval,
		s.handle,
		s.getter,
	)

	return s
}
