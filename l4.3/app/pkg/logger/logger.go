package logger

import (
	"app/internal/models"
	"context"
	"sync"

	"github.com/wb-go/wbf/zlog"
)

type stringContext string

const LoggerKey stringContext = "logger"
const BufSize = 50
const RequestIdKey stringContext = "requestId"
const ErrLevelKey = "error"
const InfoLevelKey = "info"

type Logger struct {
	Lg *zlog.Zerolog
}

func (l *Logger) LoggerWithRequestId(requestId string) *Logger {
	newLg := l.Lg.With().Str(string(RequestIdKey), requestId).Logger()
	return &Logger{&newLg}
}

func (l *Logger) Start(ch chan models.ToLog, wg *sync.WaitGroup) {
	for ent := range ch {
		rawReqId := ent.Ctx.Value(RequestIdKey)
		reqId, ok := rawReqId.(string)
		if !ok {
			l.Lg.Debug().Msg(models.ErrUnexpected.Error())
			continue
		}
		switch ent.Level {
		case InfoLevelKey:
			l.Lg.Info().Str(string(RequestIdKey), reqId).Msg(ent.Message)
		case ErrLevelKey:
			l.Lg.Error().Str(string(RequestIdKey), reqId).Err(ent.Error).Msg(ent.Message)
		default:
			l.Lg.Debug().Str(string(RequestIdKey), reqId).Msg(models.ErrUnsupportedLevel.Error())
		}
	}
	wg.Done()
}

func New() *Logger {
	zlog.InitConsole()
	return &Logger{&zlog.Logger}
}

func LoggerFromCtx(ctx context.Context) *Logger {
	if raw := ctx.Value(LoggerKey); raw != nil {
		return raw.(*Logger)
	}
	panic("failed to extract logger")
}
