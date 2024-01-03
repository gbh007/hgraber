package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
)

func (uc *UseCase) GetUnsuccessPages(ctx context.Context) []hgraber.Page {
	return uc.storage.GetUnsuccessPages(ctx)
}

func (uc *UseCase) GetUnloadedBooks(ctx context.Context) []hgraber.Book {
	return uc.storage.GetUnloadedBooks(ctx)
}

func (uc *UseCase) GetBook(ctx context.Context, id int) (hgraber.Book, error) {
	return uc.storage.GetBook(ctx, id)
}

func (uc *UseCase) GetPage(ctx context.Context, id int, page int) (*hgraber.Page, error) {
	return uc.storage.GetPage(ctx, id, page)
}

func (uc *UseCase) UpdatePageRate(ctx context.Context, id int, page int, rate int) error {
	return uc.storage.UpdatePageRate(ctx, id, page, rate)
}

func (uc *UseCase) UpdateBookRate(ctx context.Context, id int, rate int) error {
	return uc.storage.UpdateBookRate(ctx, id, rate)
}

func (uc *UseCase) ExportList(ctx context.Context) []int {
	return uc.tempStorage.ExportList(ctx)
}
