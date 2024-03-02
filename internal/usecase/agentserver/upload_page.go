package agentserver

import (
	"app/internal/domain/agent"
	"context"
	"fmt"
	"io"
)

func (uc *UseCase) UploadPage(ctx context.Context, info agent.PageInfoToUpload, body io.Reader) error {
	err := uc.files.CreatePageFile(ctx, info.BookID, info.PageNumber, info.Ext, body)
	if err != nil {
		return fmt.Errorf("upload page: create file: %w", err)
	}

	if info.URL == "" {
		err = uc.storage.UpdatePageSuccess(ctx, info.BookID, info.PageNumber, true)
	} else {
		err = uc.storage.UpdatePage(ctx, info.BookID, info.PageNumber, true, info.URL)
	}

	if err != nil {
		return fmt.Errorf("upload page: update info: %w", err)
	}

	uc.tempStorage.UnLockPageHandle(ctx, info.BookID, info.PageNumber)

	return nil
}
