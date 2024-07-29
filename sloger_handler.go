package sloger

import (
	"context"
	"io"
	"log/slog"
)

type Handler struct {
	slog.Handler
	format  LogFormat
	ctxKeys []string
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var attrs []slog.Attr
	for _, key := range h.ctxKeys {
		if data := ctx.Value(key); data != nil {
			attrs = append(attrs, slog.Any(key, data))
		}
	}
	r.AddAttrs(attrs...)
	return h.Handler.Handle(ctx, r)
}

func NewHandler(opts *slog.HandlerOptions, format LogFormat, ctxKeys []string, mw ...io.Writer) *Handler {
	var handler slog.Handler
	switch format {
	case Json:
		handler = slog.NewJSONHandler(NewMultiWriter(mw...), opts)
	case Text:
		handler = slog.NewTextHandler(NewMultiWriter(mw...), opts)
	default:
		handler = slog.NewJSONHandler(NewMultiWriter(mw...), opts)
	}
	return &Handler{
		Handler: handler,
		format:  format,
		ctxKeys: ctxKeys,
	}
}
