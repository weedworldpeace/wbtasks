package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"net/http"

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

func (h *handlers) CreateEvent(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	var ev models.Event

	if err := c.ShouldBindJSON(&ev); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	res, err := h.service.CreateEvent(c.Request.Context(), ev)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *handlers) BookEvent(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	eventID := c.Param("id")

	book := models.BookRequest{
		EventID: eventID,
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	res, err := h.service.BookEvent(c.Request.Context(), book)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *handlers) ConfirmEvent(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	eventID := c.Param("id")

	conf := models.ConfirmRequest{}
	conf.EventID = eventID

	if err := c.ShouldBindJSON(&conf); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	if err := h.service.ConfirmEvent(c.Request.Context(), conf); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handlers) GetEvent(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	eventId := c.Param("id")

	res, err := h.service.GetEvent(c.Request.Context(), eventId)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) ListEvents(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	res, err := h.service.ListEvents(c.Request.Context())
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
