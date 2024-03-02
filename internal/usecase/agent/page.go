package agent

import (
	"app/internal/domain/agent"
	"context"
	"fmt"
)

func (uc *UseCase) Pages(ctx context.Context) []agent.PageToHandle {
	pages, err := uc.agentAPI.UnprocessedPages(ctx, pageLimit)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
	}

	return pages
}

func (uc *UseCase) PageHandle(ctx context.Context, page agent.PageToHandle) {
	err := uc.pageHandle(ctx, page)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
	}
}

func (uc *UseCase) pageHandle(ctx context.Context, page agent.PageToHandle) error {
	// скачиваем изображение
	body, err := uc.loader.LoadImage(ctx, page.PageURL)
	if err != nil {
		return fmt.Errorf("page handle: download: %w", err)
	}

	defer func() {
		bodyCloseErr := body.Close()
		if bodyCloseErr != nil {
			uc.logger.ErrorContext(ctx, bodyCloseErr.Error())
		}
	}()

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
