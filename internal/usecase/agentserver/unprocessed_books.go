package agentserver

import (
	"app/internal/domain/agent"
	"context"
	"strings"
)

func (uc *UseCase) UnprocessedBooks(ctx context.Context, prefixes []string, limit int) ([]agent.BookToHandle, error) {
	if limit == 0 { // Нет смысла обрабатывать пустой запрос
		return nil, nil
	}

	// TODO: неоптимальное решение, нужна оптимизация
	books := uc.storage.GetUnloadedBooks(ctx)
	if len(books) == 0 { // Нет данных, нечего обрабатывать
		return nil, nil
	}

	result := make([]agent.BookToHandle, 0, limit)

	for _, book := range books {
		if uc.tempStorage.HasLockBookHandle(ctx, book.ID) {
			continue
		}

		var ok bool

		for _, prefix := range prefixes {
			if strings.HasPrefix(book.URL, prefix) {
				ok = true

				break
			}
		}

		if !ok && len(prefixes) > 0 {
			continue
		}

		if !uc.tempStorage.TryLockBookHandle(ctx, book.ID) {
			continue
		}

		result = append(result, agent.BookToHandle{
			ID:       book.ID,
			URL:      book.URL,
			CreateAt: book.Created,
		})

		limit--

		if limit == 0 { // Получили все данные
			break
		}
	}

	return result, nil
}
