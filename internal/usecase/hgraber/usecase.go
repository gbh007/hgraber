package hgraber

import (
	"app/internal/domain"
	"app/pkg/logger"
	"context"
	"io"
)

type storage interface {
	PagesCount(ctx context.Context) int
	BooksCount(ctx context.Context) int
	UnloadedPagesCount(ctx context.Context) int
	UnloadedBooksCount(ctx context.Context) int

	NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)

	GetUnloadedBooks(ctx context.Context) []domain.Book

	UpdateBookPages(ctx context.Context, id int, pages []domain.Page) error
	UpdateBookName(ctx context.Context, id int, name string) error
	UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error

	GetUnsuccessPages(ctx context.Context) []domain.Page
	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error

	GetPage(ctx context.Context, id int, page int) (*domain.Page, error)
	GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book
	UpdatePageRate(ctx context.Context, id int, page int, rate int) error
	UpdateBookRate(ctx context.Context, id int, rate int) error
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
	Parse(ctx context.Context, URL string) (domain.Parser, error)
	Load(ctx context.Context, URL string) (domain.Parser, error)
	LoadImage(ctx context.Context, URL string) (io.ReadCloser, error)
}

type UseCase struct {
	logger *logger.Logger

	storage storage
	files   files
	loader  loader

	tempStorage tempStorage
}

func New(storage storage, logger *logger.Logger, loader loader, files files, tempStorage tempStorage) *UseCase {
	return &UseCase{
		storage:     storage,
		logger:      logger,
		loader:      loader,
		files:       files,
		tempStorage: tempStorage,
	}
}
