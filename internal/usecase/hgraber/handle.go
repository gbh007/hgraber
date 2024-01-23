package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"net/url"
	"strings"
)

func (uc *UseCase) FirstHandle(ctx context.Context, u string) error {
	uc.logger.Info(ctx, "начата обработка", u)
	defer uc.logger.Info(ctx, "завершена обработка", u)

	u = strings.TrimSpace(u)

	if u == "" {
		return hgraber.ErrInvalidLink
	}

	var err error

	if uc.hasAgent { // Для обработки агентом может быть любой валидный адрес
		_, err = url.Parse(u)
	} else {
		_, err = uc.loader.Parse(ctx, u)
	}

	if err != nil {
		return err
	}

	_, err = uc.storage.NewBook(ctx, "", u, false)
	if err != nil {
		return err
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

		case errors.Is(err, hgraber.ErrInvalidLink):
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++
			res.Details[i] = hgraber.BookHandleResult{
				URL:         link,
				ErrorReason: err.Error(),
			}

			uc.logger.Warning(ctx, "не поддерживаемая ссылка", link)

		case err != nil:
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++
			res.Details[i] = hgraber.BookHandleResult{
				URL:         link,
				ErrorReason: err.Error(),
			}

			uc.logger.Error(ctx, err)
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
