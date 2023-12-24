package agentapi

import (
	"app/internal/domain/agent"
	"context"
	"errors"
	"io"
)

var unimplementedError = errors.New("unimplemented")

type API struct {
	prefixes []string
}

func New(prefixes []string) *API {
	return &API{
		prefixes: prefixes,
	}
}

func (api *API) UnprocessedBooks(ctx context.Context, limit int) ([]agent.BookToHandle, error) {
	return nil, unimplementedError
}

func (api *API) UnprocessedPages(ctx context.Context, limit int) ([]agent.PageToHandle, error) {
	return nil, unimplementedError
}

func (api *API) UpdateBook(ctx context.Context, book agent.BookToUpdate) error {
	return unimplementedError
}

func (api *API) UploadPage(ctx context.Context, info agent.PageInfoToUpload, body io.Reader) error {
	return unimplementedError
}
