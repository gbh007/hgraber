package converter

import (
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type storageFrom interface {
	GetBooks(ctx context.Context, filter hgraber.BookFilter) []hgraber.Book
	BooksCount(ctx context.Context) int
}

type storageTo interface {
	NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error)
	UpdateBookRate(ctx context.Context, id int, rating int) error
	UpdateBookPages(ctx context.Context, id int, pages []hgraber.Page) error
	UpdateAttributes(ctx context.Context, id int, attr hgraber.Attribute, data []string) error
}

type Builder struct {
	src storageFrom
	dst storageTo

	logger *slog.Logger
}

func New(logger *slog.Logger) *Builder {
	return &Builder{
		logger: logger,
	}
}

func (b *Builder) WithFrom(src storageFrom) *Builder {
	b.src = src

	return b
}

func (b *Builder) WithTo(dst storageTo) *Builder {
	b.dst = dst

	return b
}

func (b *Builder) Convert(ctx context.Context, offset int, notUniqWorkaround bool) {
	books := b.src.GetBooks(ctx, hgraber.BookFilter{
		Limit:  b.src.BooksCount(ctx),
		Offset: offset,
	})

	for _, book := range books {
		b.logger.InfoContext(ctx, "Начат", slog.Int("book_id", book.ID))

		id, err := b.dst.NewBook(ctx, book.Data.Name, book.URL, book.Data.Parsed.Name)
		if err != nil {
			b.logger.ErrorContext(ctx, err.Error())

			b.logger.DebugContext(ctx, "book data", slog.Any("book", book))

			if !notUniqWorkaround || !errors.Is(err, hgraber.BookAlreadyExistsError) {
				return
			}

			id, err = b.dst.NewBook(ctx, book.Data.Name, fmt.Sprintf("err (%d): %s", book.ID, book.URL), book.Data.Parsed.Name)
			if err != nil {
				b.logger.ErrorContext(ctx, err.Error())

				return
			}
		}

		if id != book.ID {
			b.logger.WarnContext(ctx, fmt.Sprintf("ID %d изменился на %d", book.ID, id))
		}

		err = b.dst.UpdateBookPages(ctx, id, book.Pages)
		if err != nil {
			b.logger.ErrorContext(ctx, err.Error())

			return
		}

		err = b.dst.UpdateBookRate(ctx, id, book.Data.Rating)
		if err != nil {
			b.logger.ErrorContext(ctx, err.Error())

			return
		}

		for attr := range book.Data.Parsed.Attributes {
			err = b.dst.UpdateAttributes(ctx, id, attr, book.Data.Attributes[attr])
			if err != nil {
				b.logger.ErrorContext(ctx, err.Error())

				return
			}
		}
	}
}
