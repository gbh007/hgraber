package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

func (uc *UseCase) FirstHandle(ctx context.Context, u string) error {
	uc.logger.InfoContext(ctx, "начата обработка", slog.String("url", u))
	defer uc.logger.InfoContext(ctx, "завершена обработка", slog.String("url", u))

	u = strings.TrimSpace(u)

	if u == "" {
		return hgraber.InvalidLinkError
	}

	if uc.hasAgent { // Для обработки агентом может быть любой валидный адрес
		_, err := url.Parse(u)
		if err != nil {
			return fmt.Errorf("handle: agent: %w", err)
		}
	} else {
		collisions, err := uc.loader.Collisions(ctx, u)
		if err != nil {
			return fmt.Errorf("handle: loader: collisions: %w", err)
		}

		for _, u := range collisions {
			u = strings.TrimSpace(u)
			_, err := uc.storage.GetBookIDByURL(ctx, u)

			if errors.Is(err, hgraber.BookNotFoundError) {
				continue
			}

			if err != nil {
				return fmt.Errorf("handle: collisions: storage: %w", err)
			}

			// Найдена коллизия
			return fmt.Errorf("%w: found collision", hgraber.BookAlreadyExistsError)
		}
	}

	_, err := uc.storage.NewBook(ctx, "", u, false)
	if err != nil {
		return fmt.Errorf("handle: storage: %w", err)
	}

	return nil
}

func (uc *UseCase) FirstHandleMultiple(ctx context.Context, data []string) (*hgraber.FirstHandleMultipleResult, error) {
	res := &hgraber.FirstHandleMultipleResult{
		NotHandled: make([]string, 0, len(data)),
		Details:    make([]hgraber.BookHandleResult, len(data)),
	}

	for i, link := range data {
		res.TotalCount++

		err := uc.FirstHandle(ctx, link)

		switch {
		case errors.Is(err, hgraber.BookAlreadyExistsError):
			res.DuplicateCount++
			res.Details[i] = hgraber.BookHandleResult{
				URL:         link,
				IsDuplicate: true,
				IsHandled:   true,
			}

		case errors.Is(err, hgraber.InvalidLinkError):
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++
			res.Details[i] = hgraber.BookHandleResult{
				URL:         link,
				ErrorReason: err.Error(),
			}

			uc.logger.WarnContext(ctx, "не поддерживаемая ссылка", slog.String("link", link))

		case err != nil:
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++
			res.Details[i] = hgraber.BookHandleResult{
				URL:         link,
				ErrorReason: err.Error(),
			}

			uc.logger.ErrorContext(ctx, err.Error())
		default:
			res.LoadedCount++
			res.Details[i] = hgraber.BookHandleResult{
				URL:       link,
				IsHandled: true,
			}
		}
	}

	return res, nil
}
