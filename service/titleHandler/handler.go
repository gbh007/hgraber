package titleHandler

import (
	"app/service/parser"
	"app/system"
	"context"
	"strings"
)

// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
func (s *Service) FirstHandle(ctx context.Context, u string) error {
	system.Info(ctx, "начата обработка", u)
	defer system.Info(ctx, "завершена обработка", u)

	u = strings.TrimSpace(u)

	_, err := parser.Parse(ctx, u)
	if err != nil {
		return err
	}

	_, err = s.Storage.NewTitle(ctx, "", u, false)
	if err != nil {
		return err
	}

	return nil
}
