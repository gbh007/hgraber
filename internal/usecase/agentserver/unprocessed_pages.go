package agentserver

import (
	"app/internal/domain/agent"
	hgDomain "app/internal/domain/hgraber"
	"context"
	"fmt"
	"strings"
)

func (uc *UseCase) UnprocessedPages(ctx context.Context, prefixes []string, limit int) ([]agent.PageToHandle, error) {
	if limit == 0 { // Нет смысла обрабатывать пустой запрос
		return nil, nil
	}

	// TODO: неоптимальное решение, нужна оптимизация
	pages := uc.storage.GetUnsuccessPages(ctx)
	if len(pages) == 0 { // Нет данных, нечего обрабатывать
		return nil, nil
	}

	books := make(map[int]hgDomain.Book)
	result := make([]agent.PageToHandle, 0, limit)

	for _, page := range pages {
		if uc.tempStorage.HasLockPageHandle(ctx, page.BookID, page.PageNumber) {
			continue
		}
		var (
			ok   bool
			book hgDomain.Book
			err  error
		)

		book, ok = books[page.BookID]
		if !ok {
			// TODO: неоптимальное решение, нужна оптимизация
			book, err = uc.storage.GetBook(ctx, page.BookID)
			if err != nil {
				return nil, fmt.Errorf("unprocessed pages: %w", err)
			}

			books[page.BookID] = book
		}

		ok = false // Сбрасываем для другого поиска

		for _, prefix := range prefixes {
			if strings.HasPrefix(book.URL, prefix) {
				ok = true

				break
			}
		}

		if !ok && len(prefixes) > 0 {
			continue
		}

		if !uc.tempStorage.TryLockPageHandle(ctx, page.BookID, page.PageNumber) {
			continue
		}

		result = append(result, agent.PageToHandle{
			BookID:     page.BookID,
			PageNumber: page.PageNumber,
			CreateAt:   book.Created, // FIXME: использовать данные страницы
			BookURL:    book.URL,
			PageURL:    page.URL,
			Ext:        page.Ext,
		})

		limit--

		if limit == 0 { // Получили все данные
			break
		}
	}

	return result, nil
}
