package agentserver

import (
	"app/internal/domain/agent"
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"log/slog"
	"slices"
)

func (uc *UseCase) UpdateBook(ctx context.Context, book agent.BookToUpdate) error {
	err := uc.storage.UpdateBookName(ctx, book.ID, book.Name)
	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	for _, attr := range book.Attributes {
		err = uc.storage.UpdateAttributes(
			ctx, book.ID,
			hgraber.Attribute(attr.Code), // FIXME: по хорошему надо их матчить более явно
			attr.Values,
		)
		if err != nil {
			return fmt.Errorf("update book: %w", err)
		}
	}

	if len(book.Pages) > 0 {
		// Для обновления страниц требуется гарантия порядка в некоторых хранилищах данных
		slices.SortFunc(book.Pages, func(a, b agent.PageToUpdate) int {
			return a.PageNumber - b.PageNumber
		})

		pagesDB := make([]hgraber.Page, len(book.Pages))

		for i, page := range book.Pages {
			pagesDB[i] = hgraber.Page{
				BookID:     book.ID,
				PageNumber: page.PageNumber,
				URL:        page.URL,
				Ext:        page.Ext,
			}
		}

		err = uc.storage.UpdateBookPages(ctx, book.ID, pagesDB)
		if err != nil {
			return fmt.Errorf("update book: %w", err)
		}
	}

	uc.logger.InfoContext(ctx, "Обновлена книга", slog.Int("book_id", book.ID))

	uc.tempStorage.UnLockBookHandle(ctx, book.ID)

	return nil
}
