package loader

import (
	"app/internal/dataprovider/loader/internal/parser/doujins_com"
	"app/internal/dataprovider/loader/internal/parser/mock"
	"app/internal/dataprovider/loader/internal/parser/multi_manga_com"
	"app/internal/dataprovider/loader/internal/request"
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type Loader struct {
	logger    *slog.Logger
	requester *request.Requester

	parsers []hgraber.Parser
}

func New(logger *slog.Logger) *Loader {
	requester := request.New(logger)

	return &Loader{
		logger:    logger,
		requester: requester,
		parsers: []hgraber.Parser{
			multi_manga_com.New(requester),
			doujins_com.New(requester),
			mock.New(requester),
		},
	}
}

func (l *Loader) Prefixes() []string {
	prefixes := make([]string, 0, len(l.parsers))

	for _, p := range l.parsers {
		prefixes = append(prefixes, p.Prefixes()...)
	}

	return prefixes
}

func (l *Loader) getParser(u string) (hgraber.Parser, error) {
	for _, p := range l.parsers {
		for _, prefix := range p.Prefixes() {
			if strings.HasPrefix(u, prefix) {
				return p, nil
			}
		}
	}

	return nil, hgraber.InvalidLinkError
}

func (l *Loader) Load(ctx context.Context, u string) (hgraber.BookParser, error) {
	p, err := l.getParser(u)
	if err != nil {
		return nil, fmt.Errorf("get parser: %w", err)
	}

	bookParser, err := p.Load(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	return bookParser, nil
}

func (l *Loader) LoadImage(ctx context.Context, u string) (io.ReadCloser, error) {
	data, err := l.requester.Request(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("load image: %w", err)
	}

	return data, nil
}

func (l *Loader) Collisions(ctx context.Context, u string) ([]string, error) {
	p, err := l.getParser(u)
	if err != nil {
		return nil, fmt.Errorf("get parser: %w", err)
	}

	for prefix, replacements := range p.Collisions() {
		if strings.HasPrefix(u, prefix) {
			res := make([]string, 0, len(replacements))

			for _, v := range replacements {
				res = append(res, strings.Replace(u, prefix, v, 1))
			}

			return res, nil
		}
	}

	return []string{}, nil // Коллизий нет
}
