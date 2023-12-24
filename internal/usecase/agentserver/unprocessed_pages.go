package agentserver

import (
	"app/internal/domain/agent"
	"context"
	"strings"
)

func (uc *UseCase) UnprocessedPages(ctx context.Context, prefixes []string, limit int) ([]agent.PageToHandle, error) {
	if limit == 0 { // Нет смысла обрабатывать пустой запрос
		return nil, nil
	}

	// TODO: неоптимальное решение, нужна оптимизация
	books := uc.storage.GetUnsuccessPages(ctx)
	if len(books) == 0 { // Нет данных, нечего обрабатывать
		return nil, nil
	}

	result := make([]agent.PageToHandle, 0, limit)

	for _, page := range books {
		if uc.tempStorage.HasLockPageHandle(ctx, page.BookID, page.PageNumber) {
			continue
		}

		var ok bool

		for _, prefix := range prefixes {
			if strings.HasPrefix(page.URL, prefix) {
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
			// FIXME: работать с данными
			// CreateAt: ,
			// BookURL: ,
			PageURL: page.URL,
		})

		limit--

		if limit == 0 { // Получили все данные
			break
		}
	}

	return result, nil
}
