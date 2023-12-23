package hgraber

import (
	"app/internal/domain"
	"context"
)

func (uc *UseCase) GetUnsuccessPages(ctx context.Context) []domain.Page {
	return uc.storage.GetUnsuccessPages(ctx)
}

func (uc *UseCase) GetUnloadedBooks(ctx context.Context) []domain.Book {
	return uc.storage.GetUnloadedBooks(ctx)
}

func (uc *UseCase) GetBook(ctx context.Context, id int) (domain.Book, error) {
	return uc.storage.GetBook(ctx, id)
}

func (uc *UseCase) GetPage(ctx context.Context, id int, page int) (*domain.Page, error) {
	return uc.storage.GetPage(ctx, id, page)
}

func (uc *UseCase) GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book {
	return uc.storage.GetBooks(ctx, filter)
}

func (uc *UseCase) UpdatePageRate(ctx context.Context, id int, page int, rate int) error {
	return uc.storage.UpdatePageRate(ctx, id, page, rate)
}

func (uc *UseCase) UpdateBookRate(ctx context.Context, id int, rate int) error {
	return uc.storage.UpdateBookRate(ctx, id, rate)
}
