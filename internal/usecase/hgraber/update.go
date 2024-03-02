package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func (uc *UseCase) ParseWithUpdate(ctx context.Context, book hgraber.Book) {
	err := uc.parseWithUpdate(ctx, book)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
	}
}

func (uc *UseCase) parseWithUpdate(ctx context.Context, book hgraber.Book) error {
	uc.logger.InfoContext(ctx, "начата обработка", slog.Int("book_id", book.ID), slog.String("url", book.URL))
	defer uc.logger.InfoContext(ctx, "завершена обработка", slog.Int("book_id", book.ID), slog.String("url", book.URL))

	p, err := uc.loader.Load(ctx, strings.TrimSpace(book.URL))
	if err != nil {
		return fmt.Errorf("parse with update: %w", err)
	}

	if !book.Data.Parsed.Name {
		err = uc.storage.UpdateBookName(ctx, book.ID, p.ParseName(ctx))
		if err != nil {
			return fmt.Errorf("parse with update: %w", err)
		}

		uc.logger.InfoContext(ctx, "обновлено название", slog.Int("book_id", book.ID), slog.String("url", book.URL))
	}

	for _, attr := range hgraber.AllAttributes {
		if book.Data.Parsed.Attributes[attr] {
			continue
		}

		err = uc.storage.UpdateAttributes(ctx, book.ID, attr, hgraber.ParseAttr(ctx, p, attr))
		if err != nil {
			return fmt.Errorf("parse with update: %w", err)
		}

		uc.logger.InfoContext(ctx, "обновлен аттрибут "+string(attr), slog.Int("book_id", book.ID), slog.String("url", book.URL))
	}

	if !book.Data.Parsed.Page {
		pages := p.ParsePages(ctx)
		if len(pages) > 0 {
			pagesDB := make([]hgraber.Page, len(pages))

			for i, page := range pages {
				pagesDB[i] = hgraber.Page{
					BookID:     book.ID,
					PageNumber: page.PageNumber,
					URL:        page.URL,
					Ext:        page.Ext,
				}
			}

			err = uc.storage.UpdateBookPages(ctx, book.ID, pagesDB)
			if err != nil {
				return fmt.Errorf("parse with update: %w", err)
			}

			uc.logger.InfoContext(ctx, "обновлены страницы", slog.Int("book_id", book.ID), slog.String("url", book.URL))
		}
	}

	return nil
}
