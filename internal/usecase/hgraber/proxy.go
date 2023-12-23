package hgraber

import (
	"app/internal/domain"
	"context"
)

func (uc *UseCases) GetUnsuccessPages(ctx context.Context) []domain.Page {
	return uc.storage.GetUnsuccessPages(ctx)
}

func (uc *UseCases) GetUnloadedBooks(ctx context.Context) []domain.Book {
	return uc.storage.GetUnloadedBooks(ctx)
}

func (uc *UseCases) GetBook(ctx context.Context, id int) (domain.Book, error) {
	return uc.storage.GetBook(ctx, id)
}

func (uc *UseCases) GetPage(ctx context.Context, id int, page int) (*domain.Page, error) {
	return uc.storage.GetPage(ctx, id, page)
}

func (uc *UseCases) GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book {
	return uc.storage.GetBooks(ctx, filter)
}

func (uc *UseCases) UpdatePageRate(ctx context.Context, id int, page int, rate int) error {
	return uc.storage.UpdatePageRate(ctx, id, page, rate)
}

func (uc *UseCases) UpdateBookRate(ctx context.Context, id int, rate int) error {
	return uc.storage.UpdateBookRate(ctx, id, rate)
}
