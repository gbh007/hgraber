package externalfile

import (
	"app/internal/domain/externalfile"
	"context"
	"fmt"
	"net/http"
)

func (s *Storage) post(ctx context.Context, request *http.Request) error {
	request.Header.Set(externalfile.HeaderToken, s.token)

	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("%s: do: %w", storageName, err)
	}

	defer s.logger.IfErrFunc(ctx, response.Body.Close)

	switch response.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil

	case http.StatusUnauthorized:
		return fmt.Errorf("%s: %w", storageName, externalfile.UnauthorizedError)
	case http.StatusForbidden:
		return fmt.Errorf("%s: %w", storageName, externalfile.ForbiddenError)
	}

	partOfBodyData := make([]byte, 100)
	n, err := response.Body.Read(partOfBodyData)
	if err != nil {
		return fmt.Errorf("%s: read code(%d) : %w", storageName, response.StatusCode, err)
	}

	if n < len(partOfBodyData) {
		partOfBodyData = partOfBodyData[:n]
	}

	return fmt.Errorf("%s: unknown error code(%d): %s", storageName, response.StatusCode, string(partOfBodyData))
}
