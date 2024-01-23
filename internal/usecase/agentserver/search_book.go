package agentserver

import (
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"fmt"
	"strings"
)

func (uc *UseCase) SearchBook(ctx context.Context, u string) (int, bool, error) {
	u = strings.TrimSpace(u)

	existsID, err := uc.storage.GetBookIDByURL(ctx, u)

	if errors.Is(err, hgraber.BookNotFoundError) {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("search url: %w", err)
	}

	return existsID, true, nil
}
