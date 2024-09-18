package sloger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	var sb strings.Builder
	r.Attrs(func(a slog.Attr) bool {
		switch a.Value.Kind() {
		case slog.KindAny:
			switch v := a.Value.Any().(type) {
			case error:
				sb.WriteString(fmt.Sprintf("%s=%+v ", a.Key, v))
			default:
				raw, err := json.Marshal(v)
				if err != nil {
					sb.WriteString(fmt.Sprintf("%s=%v ", a.Key, v))
				} else {
					sb.WriteString(fmt.Sprintf("%s=%s ", a.Key, string(raw)))
				}
			}
		default:
			sb.WriteString(fmt.Sprintf("%s=%v ", a.Key, a.Value.Any()))
		}
		return true
	})

	timeStr := r.Time.Format("[2006/01/02 15:04:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, sourceFormat(source(r)), msg, color.WhiteString(sb.String()))

	return nil
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(w, opts),
		l:       log.New(w, "", 0),
	}
	return h
}

func source(r slog.Record) *slog.Source {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	return &slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}

func sourceFormat(s *slog.Source) string {
	return fmt.Sprintf("%s:%d %s", s.File, s.Line, s.Function)
}
