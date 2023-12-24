package agent

import (
	"app/internal/domain/agent"
	"context"
	"fmt"
)

func (uc *UseCase) Pages(ctx context.Context) []agent.PageToHandle {
	pages, err := uc.agentAPI.UnprocessedPages(ctx, pageLimit)
	if err != nil {
		uc.logger.Error(ctx, err)
	}

	return pages
}

func (uc *UseCase) PageHandle(ctx context.Context, page agent.PageToHandle) {
	err := uc.pageHandle(ctx, page)
	if err != nil {
		uc.logger.Error(ctx, err)
	}
}

func (uc *UseCase) pageHandle(ctx context.Context, page agent.PageToHandle) error {
	// скачиваем изображение
	body, err := uc.loader.LoadImage(ctx, page.PageURL)
	if err != nil {
		return fmt.Errorf("page handle: download: %w", err)
	}

	defer uc.logger.IfErrFunc(ctx, body.Close)

	err = uc.agentAPI.UploadPage(
		ctx,
		agent.PageInfoToUpload{
			BookID:     page.BookID,
			PageNumber: page.PageNumber,
			Ext:        page.Ext,
		},
		body,
	)
	if err != nil {
		return fmt.Errorf("page handle: upload: %w", err)
	}

	return nil
}
