package externalfile

import (
	"app/internal/domain/externalfile"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (s *Storage) OpenPageFile(ctx context.Context, id int, page int, ext string) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url(externalfile.EndpointPage), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", storageName, err)
	}

	request.Header.Set(externalfile.HeaderToken, s.token)
	request.Header.Set(externalfile.HeaderBookID, strconv.Itoa(id))
	request.Header.Set(externalfile.HeaderPageNumber, strconv.Itoa(page))
	request.Header.Set(externalfile.HeaderPageExtension, ext)

	response, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", storageName, err)
	}

	if response.StatusCode == http.StatusOK {
		// TODO: возможно стоит использовать буффер, чтобы не произошла утечка по памяти
		return response.Body, nil
	}

	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			s.logger.ErrorContext(ctx, closeErr.Error())
		}
	}()

	switch response.StatusCode {
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%s: %w", storageName, externalfile.UnauthorizedError)
	case http.StatusForbidden:
		return nil, fmt.Errorf("%s: %w", storageName, externalfile.ForbiddenError)
	case http.StatusNotFound:
		return nil, fmt.Errorf("%s: %w", storageName, externalfile.NotFoundError)
	}

	partOfBodyData := make([]byte, 100)
	n, err := response.Body.Read(partOfBodyData)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", storageName, err)
	}

	if n < len(partOfBodyData) {
		partOfBodyData = partOfBodyData[:n]
	}

	return nil, fmt.Errorf("%s: unknown error: %s", storageName, string(partOfBodyData))
}
