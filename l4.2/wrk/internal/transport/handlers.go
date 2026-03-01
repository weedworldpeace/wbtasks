package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

type handlers struct {
	ctx     context.Context
	service ServiceInterface
}

func (h *handlers) middleware(c *ginext.Context) {
	lg := logger.LoggerFromCtx(h.ctx)
	requestId := uuid.NewString()
	lgWithReqId := lg.LoggerWithRequestId(requestId)

	lgWithReqId.Lg.Info().Str("method", c.Request.Method).Str("url", c.Request.URL.String()).Msg("received request")

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), logger.LoggerKey, lgWithReqId))
}

func (h *handlers) Cut(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	var ent models.Entity
	err := c.ShouldBindJSON(&ent)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(400, models.ErrorResponse{Error: err.Error()})
		return
	}

	res, err := h.service.Cut(ent)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(200, models.Response{Data: res})
	lg.Info().Msg("request processed")

}
