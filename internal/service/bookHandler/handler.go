package bookHandler

import (
	"app/internal/domain"
	"app/internal/service/bookHandler/internal/parser"
	"context"
	"errors"
	"strings"
)

// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
func (s *Service) FirstHandle(ctx context.Context, u string) error {
	s.logger.Info(ctx, "начата обработка", u)
	defer s.logger.Info(ctx, "завершена обработка", u)

	u = strings.TrimSpace(u)

	_, err := parser.Parse(ctx, u)
	if err != nil {
		return err
	}

	_, err = s.storage.NewBook(ctx, "", u, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FirstHandleMultiple(ctx context.Context, data []string) domain.FirstHandleMultipleResult {
	res := domain.FirstHandleMultipleResult{}

	for _, link := range data {
		res.TotalCount++

		err := s.FirstHandle(ctx, link)

		switch {
		case errors.Is(err, domain.BookAlreadyExistsError):
			res.DuplicateCount++

		case errors.Is(err, parser.ErrInvalidLink):
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++

			s.logger.Warning(ctx, "не поддерживаемая ссылка", link)

		case err != nil:
			res.NotHandled = append(res.NotHandled, link)
			res.ErrorCount++

			s.logger.Error(ctx, err)
		default:
			res.LoadedCount++
		}
	}

	return res
}
