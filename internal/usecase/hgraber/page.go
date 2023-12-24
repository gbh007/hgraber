package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"io"
)

func (uc *UseCase) LoadPageWithUpdate(ctx context.Context, page hgraber.Page) error {
	err := uc.downloadPageImage(ctx, page.BookID, page.PageNumber, page.URL, page.Ext)
	if err != nil {
		return fmt.Errorf("load page with update: %w", err)
	}

	err = uc.storage.UpdatePageSuccess(ctx, page.BookID, page.PageNumber, true)
	if err != nil {
		return fmt.Errorf("load page with update: %w", err)
	}

	return nil
}

func (uc *UseCase) downloadPageImage(ctx context.Context, id, page int, URL, ext string) error {
	// скачиваем изображение
	body, err := uc.loader.LoadImage(ctx, URL)
	if err != nil {
		return fmt.Errorf("download page image: %w", err)
	}

	defer uc.logger.IfErrFunc(ctx, body.Close)

	// создаем файл и загружаем туда изображение
	err = uc.files.CreatePageFile(ctx, id, page, ext, body)
	if err != nil {
		return fmt.Errorf("download page image: %w", err)
	}

	return nil
}

func (uc *UseCase) PageWithBody(ctx context.Context, bookID int, pageNumber int) (*hgraber.Page, io.ReadCloser, error) {
	info, err := uc.storage.GetPage(ctx, bookID, pageNumber)
	if err != nil {
		return nil, nil, fmt.Errorf("page with body: %w", err)
	}

	r, err := uc.files.OpenPageFile(ctx, bookID, pageNumber, info.Ext)
	if err != nil {
		return nil, nil, fmt.Errorf("page with body: %w", err)
	}

	return info, r, nil
}
