package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/zlog"
)

type stringContext string

const LoggerKey stringContext = "logger"
const requestIdKey = "requestId"

type Logger struct {
	Lg *zlog.Zerolog
}

func (l *Logger) LoggerWithRequestId(requestId string) *Logger {
	newLg := l.Lg.With().Str(requestIdKey, requestId).Logger()
	return &Logger{&newLg}
}

func New() *Logger {
	zlog.InitConsole()
	return &Logger{&zlog.Logger}
}

func (l *Logger) LoggerLevel(level int) *Logger {
	newLg := l.Lg.Level(zerolog.Level(level))
	return &Logger{&newLg}
}

func LoggerFromCtx(ctx context.Context) *Logger {
	if raw := ctx.Value(LoggerKey); raw != nil {
		return raw.(*Logger)
	}
	panic("failed to extract logger")
}
