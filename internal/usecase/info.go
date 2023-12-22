package usecase

import (
	"app/internal/domain"
	"context"
)

func (uc *UseCases) Info(ctx context.Context) (*domain.MainInfo, error) {
	info := &domain.MainInfo{
		BookCount:        uc.storage.BooksCount(ctx),
		NotLoadBookCount: uc.storage.UnloadedBooksCount(ctx),
		PageCount:        uc.storage.PagesCount(ctx),
		NotLoadPageCount: uc.storage.UnloadedPagesCount(ctx),
	}

	return info, nil
}
