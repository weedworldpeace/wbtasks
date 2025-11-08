package logger

import (
	"context"
	"log/slog"
	"os"
)

type stringContext string

const LoggerKey stringContext = "logger"
const requestIdKey = "requestId"

type Logger struct {
	Lg *slog.Logger
}

func (l *Logger) LoggerWithRequestId(requestId string) *Logger {
	return &Logger{l.Lg.With(requestIdKey, requestId)}
}

func New() *Logger {
	return &Logger{slog.New(slog.NewJSONHandler(os.Stdout, nil))}
}

func LoggerFromCtx(ctx context.Context) *Logger {
	if raw := ctx.Value(LoggerKey); raw != nil {
		return raw.(*Logger)
	}
	panic("failed to extract logger")
}
