package lib

import (
	"log/slog"
	"os"
)

type CustomHandler struct {
	slog.Handler
}

func NewCustomLogger() *slog.Logger {
	svcName := os.Getenv("SERVICE_NAME")
	if svcName == "" {
		if os.Getenv("LOCAL_ONLY") != "true" {
			slog.Error(
				"SERVICE_NAME is not set",
				slog.String("error", "SERVICE_NAME must be set"),
			)
			panic("SERVICE_NAME must be set")
		} else {
			svcName = "fetch-chat-function"
		}
	}

	handler := CustomHandler{slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case slog.MessageKey:
					a = slog.Attr{
						Key:   "message",
						Value: a.Value,
					}
				case slog.LevelKey:
					a = slog.Attr{
						Key:   "severity",
						Value: a.Value,
					}
				case slog.SourceKey:
					a = slog.Attr{
						Key:   "logging.googleapis.com/sourceLocation",
						Value: a.Value,
					}
				}
				return a
			},
		}),
	}

	logger := slog.New(&handler).With(
		slog.Group("logging.googleapis.com/labels",
			slog.String("service", svcName),
		))

	return logger
}
