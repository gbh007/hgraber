package hgraber

import (
	"app/internal/domain"
	"context"
	"fmt"
	"io"
)

func (uc *UseCase) LoadPageWithUpdate(ctx context.Context, page domain.Page) error {
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
	data, err := uc.loader.LoadImage(ctx, URL)
	if err != nil {
		return fmt.Errorf("download page image: %w", err)
	}

	// создаем файл и загружаем туда изображение
	f, err := uc.files.CreatePageFile(ctx, id, page, ext)
	if err != nil {
		return fmt.Errorf("download page image: %w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		uc.logger.IfErr(ctx, f.Close())

		return fmt.Errorf("download page image: %w", err)
	}

	return f.Close()
}

func (uc *UseCase) PageWithBody(ctx context.Context, bookID int, pageNumber int) (*domain.Page, io.ReadCloser, error) {
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
