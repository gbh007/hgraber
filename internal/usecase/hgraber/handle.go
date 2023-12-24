package hgraber

import (
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"strings"
)

func (uc *UseCase) FirstHandle(ctx context.Context, u string) error {
	uc.logger.Info(ctx, "начата обработка", u)
	defer uc.logger.Info(ctx, "завершена обработка", u)

	u = strings.TrimSpace(u)

	_, err := uc.loader.Parse(ctx, u)
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
	res := &hgraber.FirstHandleMultipleResult{}

	for _, link := range data {
		res.TotalCount++

		err := uc.FirstHandle(ctx, link)

		switch {
		case errors.Is(err, hgraber.BookAlreadyExistsError):
			res.DuplicateCount++

		case errors.Is(err, hgraber.ErrInvalidLink):
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++

			uc.logger.Warning(ctx, "не поддерживаемая ссылка", link)

		case err != nil:
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++

			uc.logger.Error(ctx, err)
		default:
			res.LoadedCount++
		}
	}

	return res, nil
}
