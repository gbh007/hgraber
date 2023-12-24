package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
)

func (uc *UseCase) Info(ctx context.Context) (*hgraber.MainInfo, error) {
	info := &hgraber.MainInfo{
		BookCount:        uc.storage.BooksCount(ctx),
		NotLoadBookCount: uc.storage.UnloadedBooksCount(ctx),
		PageCount:        uc.storage.PagesCount(ctx),
		NotLoadPageCount: uc.storage.UnloadedPagesCount(ctx),
	}

	return info, nil
}
