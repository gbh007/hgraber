package agentserver

import (
	"app/internal/domain/agent"
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (uc *UseCase) firstHandle(ctx context.Context, u string) (int, error) {
	uc.logger.Info(ctx, "начата обработка", u)
	defer uc.logger.Info(ctx, "завершена обработка", u)

	u = strings.TrimSpace(u)

	if u == "" {
		return 0, hgraber.ErrInvalidLink
	}

	_, err := url.Parse(u)
	if err != nil {
		return 0, fmt.Errorf("parse url: %w", err)
	}

	existsID, err := uc.storage.GetBookIDByURL(ctx, u)
	if err == nil { // Книга уже существует, выходим с ИД и ошибкой для дальнейшей обработки
		return existsID, hgraber.BookAlreadyExistsError
	}

	if err != nil && !errors.Is(err, hgraber.BookNotFoundError) {
		return 0, fmt.Errorf("search url: %w", err)
	}

	id, err := uc.storage.NewBook(ctx, "", u, false)
	if err != nil {
		return 0, fmt.Errorf("create book: %w", err)
	}

	return id, nil
}

func (uc *UseCase) CreateMultipleBook(ctx context.Context, data []string) (*agent.CreateBooksResult, error) {
	res := &agent.CreateBooksResult{
		NotHandled: make([]string, 0, len(data)),
		Details:    make([]agent.CreateBookResult, len(data)),
	}

	for i, link := range data {
		res.TotalCount++

		bookID, err := uc.firstHandle(ctx, link)
		res.Details[i] = agent.CreateBookResult{
			URL: link,
			ID:  bookID,
		}

		switch {
		case errors.Is(err, hgraber.BookAlreadyExistsError):
			res.DuplicateCount++
			res.Details[i].IsDuplicate = true
			res.Details[i].IsHandled = true

		case errors.Is(err, hgraber.ErrInvalidLink):
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++
			res.Details[i].ErrorReason = err.Error()

			uc.logger.Warning(ctx, "не поддерживаемая ссылка", link)

		case err != nil:
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++
			res.Details[i].ErrorReason = err.Error()

			uc.logger.Error(ctx, err)
		default:
			res.LoadedCount++
			res.Details[i].IsHandled = true
		}
	}

	return res, nil
}
