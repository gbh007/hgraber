package externalfile

import (
	"app/internal/dto"
	"app/system"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type fileWriter struct {
	token string

	bookID     int
	pageNumber int
	pageExt    string

	filename string

	buff *bytes.Buffer

	client *http.Client
	url    string
	// FIXME: это очеееееееееееень плохо, но с текущей сигнатурой пока никак по другому.
	requestContext context.Context
}

func (s *Storage) newPageFileWriter(ctx context.Context, bookID int, pageNumber int, pageExt string) *fileWriter {
	return &fileWriter{
		requestContext: ctx,

		token: s.token,

		bookID:     bookID,
		pageNumber: pageNumber,
		pageExt:    pageExt,

		buff:   &bytes.Buffer{},
		client: s.client,
		url: (&url.URL{
			Scheme: s.scheme,
			Host:   s.hostWithPort,
			Path:   dto.ExternalFileEndpointPage,
		}).String(),
	}
}

func (s *Storage) newExportFileWriter(ctx context.Context, filename string) *fileWriter {
	return &fileWriter{
		requestContext: ctx,

		token: s.token,

		filename: filename,

		buff:   &bytes.Buffer{},
		client: s.client,
		url: (&url.URL{
			Scheme: s.scheme,
			Host:   s.hostWithPort,
			Path:   dto.ExternalFileEndpointExport,
		}).String(),
	}
}

func (fw *fileWriter) Write(p []byte) (int, error) {
	if fw.buff == nil {
		return 0, fmt.Errorf("%s: write to nil buffer", storageName)
	}

	n, err := fw.buff.Write(p)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", storageName, err)
	}

	return n, nil
}

func (fw *fileWriter) Close() error {
	if fw.buff == nil {
		return nil
	}

	request, err := http.NewRequestWithContext(fw.requestContext, http.MethodPost, fw.url, fw.buff)
	if err != nil {
		return fmt.Errorf("%s: %w", storageName, err)
	}

	fw.setHeaders(request)

	response, err := fw.client.Do(request)
	if err != nil {
		return fmt.Errorf("%s: %w", storageName, err)
	}

	defer system.IfErrFunc(fw.requestContext, response.Body.Close)

	// Принудительно удаляем буффер на случай некорректных вызовов.
	fw.buff = nil

	switch response.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil

	case http.StatusUnauthorized:
		return fmt.Errorf("%s: %w", storageName, dto.ExternalFileUnauthorizedError)
	case http.StatusForbidden:
		return fmt.Errorf("%s: %w", storageName, dto.ExternalFileForbiddenError)
	}

	partOfBodyData := make([]byte, 100)
	n, err := response.Body.Read(partOfBodyData)
	if err != nil {
		return fmt.Errorf("%s: %w", storageName, err)
	}

	if n < len(partOfBodyData) {
		partOfBodyData = partOfBodyData[:n]
	}

	return fmt.Errorf("%s: unknown error: %s", storageName, string(partOfBodyData))
}

func (fw *fileWriter) setHeaders(request *http.Request) {
	request.Header.Set(dto.ExternalFileToken, fw.token)

	if fw.bookID != 0 {
		request.Header.Set(dto.ExternalFileBookID, strconv.Itoa(fw.bookID))
	}

	if fw.pageNumber != 0 {
		request.Header.Set(dto.ExternalFilePageNumber, strconv.Itoa(fw.pageNumber))
	}

	if fw.pageExt != "" {
		request.Header.Set(dto.ExternalFilePageExtension, fw.pageExt)
	}

	if fw.filename != "" {
		request.Header.Set(dto.ExternalFileFilename, fw.filename)
	}
}
