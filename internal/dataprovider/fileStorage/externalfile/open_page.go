package externalfile

import (
	"app/internal/dataprovider/fileStorage/externalfile/dto"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (s *Storage) OpenPageFile(ctx context.Context, id int, page int, ext string) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url(dto.ExternalFileEndpointPage), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", storageName, err)
	}

	request.Header.Set(dto.ExternalFileToken, s.token)
	request.Header.Set(dto.ExternalFileBookID, strconv.Itoa(id))
	request.Header.Set(dto.ExternalFilePageNumber, strconv.Itoa(page))
	request.Header.Set(dto.ExternalFilePageExtension, ext)

	response, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", storageName, err)
	}

	if response.StatusCode == http.StatusOK {
		// TODO: возможно стоит использовать буффер, чтобы не произошла утечка по памяти
		return response.Body, nil
	}

	defer s.logger.IfErrFunc(ctx, response.Body.Close)

	switch response.StatusCode {
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%s: %w", storageName, dto.ExternalFileUnauthorizedError)
	case http.StatusForbidden:
		return nil, fmt.Errorf("%s: %w", storageName, dto.ExternalFileForbiddenError)
	case http.StatusNotFound:
		return nil, fmt.Errorf("%s: %w", storageName, dto.ExternalFileNotFoundError)
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
