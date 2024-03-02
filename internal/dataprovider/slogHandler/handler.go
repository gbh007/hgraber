package slogHandler

import (
	"context"
	"log/slog"
)

var _ slog.Handler = (*Handler)(nil)

type handlerCtxHook func(context.Context) []slog.Attr

type Handler struct {
	isDebug bool
	hooks   []handlerCtxHook
	layers  []handlerLayer
	printer Printer
}

func New(options ...HandlerOption) *Handler {
	cfg := &handlerConfig{
		printer: stdoutPrinter,
	}

	for _, opt := range options {
		opt.apply(cfg)
	}

	return &Handler{
		isDebug: cfg.isDebug,
		hooks:   cfg.hooks,
		printer: cfg.printer,
	}
}

func (h *Handler) Enabled(_ context.Context, lv slog.Level) bool {
	if h.isDebug {
		return true
	}

	if lv == slog.LevelDebug {
		return false
	}

	return true
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
	recordAttrs := make([]slog.Attr, 0, rec.NumAttrs())
	rec.Attrs(func(a slog.Attr) bool {
		recordAttrs = append(recordAttrs, a)

		return true
	})

	attrs := layersToAttrs(layersWithAttrs(h.layers, recordAttrs))

	for _, hook := range h.hooks {
		attrs = append(attrs, hook(ctx)...)
	}

	return h.printer(ctx, rec.Time, rec.Message, rec.Level, attrs)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.copyWithLayers(layersWithAttrs(h.layers, attrs))
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h.copyWithLayers(layersWithGroup(h.layers, name))
}

func (h *Handler) copyWithLayers(layers []handlerLayer) *Handler {
	return &Handler{
		isDebug: h.isDebug,
		layers:  layers,
		hooks:   h.hooks,
		printer: h.printer,
	}
}
