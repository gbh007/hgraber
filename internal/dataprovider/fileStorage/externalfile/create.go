package externalfile

import (
	"app/internal/dataprovider/fileStorage/externalfile/dto"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (s *Storage) CreateExportFile(ctx context.Context, name string, body io.Reader) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url(dto.ExternalFileEndpointExport), body)
	if err != nil {
		return fmt.Errorf("%s: request: %w", storageName, err)
	}

	request.Header.Set(dto.ExternalFileFilename, name)

	return s.post(ctx, request)
}

func (s *Storage) CreatePageFile(ctx context.Context, id int, page int, ext string, body io.Reader) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url(dto.ExternalFileEndpointPage), body)
	if err != nil {
		return fmt.Errorf("%s: request: %w", storageName, err)
	}

	request.Header.Set(dto.ExternalFileBookID, strconv.Itoa(id))
	request.Header.Set(dto.ExternalFilePageNumber, strconv.Itoa(page))
	request.Header.Set(dto.ExternalFilePageExtension, ext)

	return s.post(ctx, request)
}
