package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
	"io"
	"log/slog"
)

type storage interface {
	PagesCount(ctx context.Context) int
	BooksCount(ctx context.Context) int
	UnloadedPagesCount(ctx context.Context) int
	UnloadedBooksCount(ctx context.Context) int
	PagesSize(ctx context.Context) int64

	NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error)
	GetBook(ctx context.Context, id int) (hgraber.Book, error)

	GetBookIDByURL(ctx context.Context, url string) (int, error)

	GetUnloadedBooks(ctx context.Context) []hgraber.Book

	UpdateBookPages(ctx context.Context, id int, pages []hgraber.Page) error
	UpdateBookName(ctx context.Context, id int, name string) error
	UpdateAttributes(ctx context.Context, id int, attr hgraber.Attribute, data []string) error

	GetUnsuccessPages(ctx context.Context) []hgraber.Page
	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error

	GetPage(ctx context.Context, id int, page int) (*hgraber.Page, error)
	GetBooks(ctx context.Context, filter hgraber.BookFilter) []hgraber.Book
	UpdatePageRate(ctx context.Context, id int, page int, rating int) error
	UpdateBookRate(ctx context.Context, id int, rating int) error
}

type tempStorage interface {
	AddExport(ctx context.Context, bookID int)
	ExportList(ctx context.Context) []int
}

type files interface {
	CreatePageFile(ctx context.Context, id, page int, ext string, body io.Reader) error
	OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error)
	CreateExportFile(ctx context.Context, name string, body io.Reader) error
}

type loader interface {
	Collisions(ctx context.Context, u string) ([]string, error)
	Load(ctx context.Context, u string) (hgraber.BookParser, error)
	LoadImage(ctx context.Context, u string) (io.ReadCloser, error)
}

type UseCase struct {
	logger *slog.Logger

	storage storage
	files   files
	loader  loader

	hasAgent bool

	tempStorage tempStorage
}

func New(storage storage, logger *slog.Logger, loader loader, files files, tempStorage tempStorage, hasAgent bool) *UseCase {
	return &UseCase{
		storage:     storage,
		logger:      logger,
		loader:      loader,
		files:       files,
		tempStorage: tempStorage,
		hasAgent:    hasAgent,
	}
}
