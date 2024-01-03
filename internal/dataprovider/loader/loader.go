package loader

import (
	"app/internal/dataprovider/loader/internal/parser"
	"app/internal/dataprovider/loader/internal/request"
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"io"
)

type logger interface {
	Error(ctx context.Context, err error)
	IfErrFunc(ctx context.Context, f func() error)
}

type Loader struct {
	logger    logger
	requester *request.Requester
}

func New(logger logger) *Loader {
	return &Loader{
		logger:    logger,
		requester: request.New(logger),
	}
}

func (l *Loader) Prefixes() []string {
	return parser.Prefixes
}

func (l *Loader) Parse(ctx context.Context, URL string) (hgraber.Parser, error) {
	p, err := parser.Parse(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	return p, nil
}

func (l *Loader) Load(ctx context.Context, URL string) (hgraber.Parser, error) {
	p, err := parser.Parse(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	err = p.Load(ctx, l.requester, URL)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	return p, nil
}

func (l *Loader) LoadImage(ctx context.Context, URL string) (io.ReadCloser, error) {
	data, err := l.requester.Request(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("load image: %w", err)
	}

	return data, nil
}
