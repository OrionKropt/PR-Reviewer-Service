package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
)

type LogHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type LogHandler struct {
	slog.Handler
	logger *log.Logger
}

func (h *LogHandler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String() + ":"

	fields := make(map[string]interface{}, record.NumAttrs())
	record.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	bytes, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := record.Time.Format("[00:00:00.000]")

	h.logger.Println(timeStr, level, record.Message, string(bytes))
	return nil
}

func NewLogHandler(out io.Writer, opts LogHandlerOptions) *LogHandler {
	handler := &LogHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		logger:  log.New(out, "", 0),
	}

	return handler
}

func ParseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func NewTest() *slog.Logger {
	h := slog.NewTextHandler(io.Discard, nil)
	return slog.New(h)
}
