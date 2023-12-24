package externalfile

import (
	"app/internal/domain/externalfile"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (s *Storage) CreateExportFile(ctx context.Context, name string, body io.Reader) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url(externalfile.EndpointExport), body)
	if err != nil {
		return fmt.Errorf("%s: request: %w", storageName, err)
	}

	request.Header.Set(externalfile.HeaderFilename, name)

	return s.post(ctx, request)
}

func (s *Storage) CreatePageFile(ctx context.Context, id int, page int, ext string, body io.Reader) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url(externalfile.EndpointPage), body)
	if err != nil {
		return fmt.Errorf("%s: request: %w", storageName, err)
	}

	request.Header.Set(externalfile.HeaderBookID, strconv.Itoa(id))
	request.Header.Set(externalfile.HeaderPageNumber, strconv.Itoa(page))
	request.Header.Set(externalfile.HeaderPageExtension, ext)

	return s.post(ctx, request)
}
