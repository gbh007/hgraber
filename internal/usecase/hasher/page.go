package hasher

import (
	"app/internal/domain/hgraber"
	"context"
	"crypto/md5"
	"fmt"
	"io"
)

func (uc *UseCase) UnHashedPages(ctx context.Context) []hgraber.Page {
	return uc.storage.GetUnHashedPages(ctx)
}

func (uc *UseCase) HandlePage(ctx context.Context, page hgraber.Page) error {
	body, err := uc.files.OpenPageFile(ctx, page.BookID, page.PageNumber, page.Ext)
	if err != nil {
		return fmt.Errorf("get page body for hashing: %w", err)
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read page body for hashing: %w", err)
	}

	size := len(data)
	hash := fmt.Sprintf("%x", md5.Sum(data))

	err = uc.storage.UpdatePageHash(ctx, page.BookID, page.PageNumber, hash, int64(size))
	if err != nil {
		return fmt.Errorf("update page hash: %w", err)
	}

	return nil
}
