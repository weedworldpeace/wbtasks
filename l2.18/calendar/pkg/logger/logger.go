package logger

import (
	"context"
	"log/slog"
	"os"
)

const (
	ContextLogger = "logger"
)

type Logger struct {
	Lg *slog.Logger
}

func New() *Logger {
	return &Logger{slog.New(slog.NewJSONHandler(os.Stdout, nil))}
}

func LoggerFromCtx(ctx context.Context) *Logger {
	if raw := ctx.Value(ContextLogger); raw != nil {
		return raw.(*Logger)
	}
	panic("failed to extract logger")
}
