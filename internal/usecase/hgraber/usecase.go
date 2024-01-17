package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
	"io"
)

type logger interface {
	Error(ctx context.Context, err error)
	IfErrFunc(ctx context.Context, f func() error)
	Info(ctx context.Context, args ...any)
	Warning(ctx context.Context, args ...any)
}

type storage interface {
	PagesCount(ctx context.Context) int
	BooksCount(ctx context.Context) int
	UnloadedPagesCount(ctx context.Context) int
	UnloadedBooksCount(ctx context.Context) int

	NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error)
	GetBook(ctx context.Context, id int) (hgraber.Book, error)

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
	Parse(ctx context.Context, URL string) (hgraber.Parser, error)
	Load(ctx context.Context, URL string) (hgraber.Parser, error)
	LoadImage(ctx context.Context, URL string) (io.ReadCloser, error)
}

type UseCase struct {
	logger logger

	storage storage
	files   files
	loader  loader

	hasAgent bool

	tempStorage tempStorage
}

func New(storage storage, logger logger, loader loader, files files, tempStorage tempStorage, hasAgent bool) *UseCase {
	return &UseCase{
		storage:     storage,
		logger:      logger,
		loader:      loader,
		files:       files,
		tempStorage: tempStorage,
		hasAgent:    hasAgent,
	}
}
