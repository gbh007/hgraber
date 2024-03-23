package hasher

import (
	"app/internal/domain/hgraber"
	"context"
	"io"
)

type storage interface {
	GetUnHashedPages(ctx context.Context) []hgraber.Page
	UpdatePageHash(ctx context.Context, id int, page int, hash string, size int64) error
}

type files interface {
	OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error)
}

type UseCase struct {
	storage storage
	files   files
}

func New(storage storage, files files) *UseCase {
	return &UseCase{
		storage: storage,
		files:   files,
	}
}
