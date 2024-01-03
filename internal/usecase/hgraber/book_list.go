package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
)

func (uc *UseCase) GetBooks(ctx context.Context, filter hgraber.BookFilterOuter) hgraber.FilteredBooks {
	limit, offset := pageToLimit(filter.Page, filter.Count)

	books := uc.storage.GetBooks(ctx, hgraber.BookFilter{
		Limit:    limit,
		Offset:   offset,
		NewFirst: filter.NewFirst,
	})

	totalPages := totalToPages(uc.storage.BooksCount(ctx), filter.Count)

	return hgraber.FilteredBooks{
		Books:       books,
		Pages:       generatePagination(filter.Page, totalPages),
		CurrentPage: filter.Page,
	}
}
