package stopwatch

import (
	"app/internal/domain"
	"app/system"
	"context"
)

const depth = 0

type storage interface {
	GetPage(ctx context.Context, id int, page int) (*domain.PageFullInfo, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book
	GetUnloadedBooks(ctx context.Context) []domain.Book
	GetUnsuccessPages(ctx context.Context) []domain.PageFullInfo
	NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error)
	PagesCount(ctx context.Context) int
	BooksCount(ctx context.Context) int
	UnloadedPagesCount(ctx context.Context) int
	UnloadedBooksCount(ctx context.Context) int

	UpdatePageRate(ctx context.Context, id int, page int, rate int) error
	UpdateBookRate(ctx context.Context, id int, rate int) error

	UpdateBookName(ctx context.Context, id int, name string) error

	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error
	UpdateBookPages(ctx context.Context, id int, pages []domain.Page) error

	UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error
}

type Stopwatch struct {
	storage storage
}

func WithStopwatch(storage storage) *Stopwatch {
	return &Stopwatch{storage: storage}
}

func (s *Stopwatch) UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateAttributes ("+string(attr)+")", depth)()

	return s.storage.UpdateAttributes(ctx, id, attr, data)
}

func (s *Stopwatch) GetPage(ctx context.Context, id int, page int) (*domain.PageFullInfo, error) {
	defer system.StopwatchWithDepth(ctx, "DB - GetPage", depth)()

	return s.storage.GetPage(ctx, id, page)
}

func (s *Stopwatch) GetBook(ctx context.Context, id int) (domain.Book, error) {
	defer system.StopwatchWithDepth(ctx, "DB - GetBook", depth)()

	return s.storage.GetBook(ctx, id)
}

func (s *Stopwatch) GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book {
	defer system.StopwatchWithDepth(ctx, "DB - GetBooks", depth)()

	return s.storage.GetBooks(ctx, filter)
}

func (s *Stopwatch) GetUnloadedBooks(ctx context.Context) []domain.Book {
	defer system.StopwatchWithDepth(ctx, "DB - GetUnloadedBooks", depth)()

	return s.storage.GetUnloadedBooks(ctx)
}

func (s *Stopwatch) GetUnsuccessPages(ctx context.Context) []domain.PageFullInfo {
	defer system.StopwatchWithDepth(ctx, "DB - GetUnsuccessPages", depth)()

	return s.storage.GetUnsuccessPages(ctx)
}

func (s *Stopwatch) NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error) {
	defer system.StopwatchWithDepth(ctx, "DB - NewBook", depth)()

	return s.storage.NewBook(ctx, name, URL, loaded)
}

func (s *Stopwatch) PagesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - PagesCount", depth)()

	return s.storage.PagesCount(ctx)
}

func (s *Stopwatch) BooksCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - BooksCount", depth)()

	return s.storage.BooksCount(ctx)
}

func (s *Stopwatch) UnloadedPagesCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - UnloadedPagesCount", depth)()

	return s.storage.UnloadedPagesCount(ctx)
}

func (s *Stopwatch) UnloadedBooksCount(ctx context.Context) int {
	defer system.StopwatchWithDepth(ctx, "DB - UnloadedBooksCount", depth)()

	return s.storage.UnloadedBooksCount(ctx)
}

func (s *Stopwatch) UpdatePageRate(ctx context.Context, id int, page int, rate int) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdatePageRate", depth)()

	return s.storage.UpdatePageRate(ctx, id, page, rate)
}

func (s *Stopwatch) UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdatePageSuccess", depth)()

	return s.storage.UpdatePageSuccess(ctx, id, page, success)
}
func (s *Stopwatch) UpdateBookName(ctx context.Context, id int, name string) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateBookName", depth)()

	return s.storage.UpdateBookName(ctx, id, name)
}

func (s *Stopwatch) UpdateBookPages(ctx context.Context, id int, pages []domain.Page) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateBookPages", depth)()

	return s.storage.UpdateBookPages(ctx, id, pages)
}

func (s *Stopwatch) UpdateBookRate(ctx context.Context, id int, rate int) error {
	defer system.StopwatchWithDepth(ctx, "DB - UpdateBookRate", depth)()

	return s.storage.UpdateBookRate(ctx, id, rate)
}
