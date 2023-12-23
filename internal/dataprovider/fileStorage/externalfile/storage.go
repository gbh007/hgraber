package externalfile

import (
	"app/internal/dataprovider/fileStorage/externalfile/dto"
	"app/pkg/logger"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const storageName = "external file"

type Storage struct {
	token string

	scheme       string
	hostWithPort string

	client *http.Client

	logger *logger.Logger
}

func New(token string, scheme string, hostWithPort string, logger *logger.Logger) *Storage {
	return &Storage{
		token:        token,
		scheme:       scheme,
		hostWithPort: hostWithPort,
		client: &http.Client{
			Timeout: time.Minute,
		},
		logger: logger,
	}
}

func (s *Storage) CreateExportFile(ctx context.Context, name string) (io.WriteCloser, error) {
	return s.newExportFileWriter(ctx, name), nil
}

func (s *Storage) CreatePageFile(ctx context.Context, id int, page int, ext string) (io.WriteCloser, error) {
	return s.newPageFileWriter(ctx, id, page, ext), nil
}

func (s *Storage) OpenPageFile(ctx context.Context, id int, page int, ext string) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, (&url.URL{
		Scheme: s.scheme,
		Host:   s.hostWithPort,
		Path:   dto.ExternalFileEndpointPage,
	}).String(), nil)
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
