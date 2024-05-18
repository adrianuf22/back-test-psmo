package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal"
)

func main() {
	ctx := context.Background()
	initLogger()

	slog.Info("Initializing...")
	app := internal.NewServer(ctx)
	app.Run(ctx)
}

func initLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, s slog.Attr) slog.Attr {
			if s.Key == slog.TimeKey {
				s.Value = slog.StringValue(s.Value.Time().Format(time.RFC3339))
			}

			if s.Key == slog.MessageKey {
				s.Key = "message"
			}

			return s
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
