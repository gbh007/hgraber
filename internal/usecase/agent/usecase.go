package agent

import (
	"app/internal/domain/agent"
	"app/internal/domain/hgraber"
	"context"
	"io"
	"log/slog"
)

const (
	booksLimit = 30
	pageLimit  = 50
)

type agentAPI interface {
	UnprocessedBooks(ctx context.Context, limit int) ([]agent.BookToHandle, error)
	UnprocessedPages(ctx context.Context, limit int) ([]agent.PageToHandle, error)
	UpdateBook(ctx context.Context, book agent.BookToUpdate) error
	UploadPage(ctx context.Context, info agent.PageInfoToUpload, body io.Reader) error
}

type loader interface {
	Parse(ctx context.Context, URL string) (hgraber.Parser, error)
	Load(ctx context.Context, URL string) (hgraber.Parser, error)
	LoadImage(ctx context.Context, URL string) (io.ReadCloser, error)
}

type UseCase struct {
	logger *slog.Logger

	agentAPI agentAPI
	loader   loader
}

func New(logger *slog.Logger, agentAPI agentAPI, loader loader) *UseCase {
	return &UseCase{
		logger: logger,

		agentAPI: agentAPI,
		loader:   loader,
	}
}
