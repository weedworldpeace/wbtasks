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

func (h *handlers) CreateTransaction(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	var trans models.Transaction

	if err := c.ShouldBindJSON(&trans); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	res, err := h.service.CreateTransaction(c.Request.Context(), trans)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *handlers) ListTransactions(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	from := c.Query("from")
	to := c.Query("to")
	limit := c.Query("limit")
	offset := c.Query("offset")

	res, err := h.service.ListTransactions(c.Request.Context(), from, to, limit, offset)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) GetTransaction(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	id := c.Param("id")

	res, err := h.service.GetTransaction(c.Request.Context(), id)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) UpdateTransaction(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	id := c.Param("id")

	var trans models.Transaction

	if err := c.ShouldBindJSON(&trans); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	trans.ID = id

	res, err := h.service.UpdateTransaction(c.Request.Context(), trans)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) DeleteTransaction(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	id := c.Param("id")

	err := h.service.DeleteTransaction(c.Request.Context(), id)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handlers) GetAnalytics(c *ginext.Context) {
	lg := logger.LoggerFromCtx(h.ctx).Lg

	from := c.Query("from")
	to := c.Query("to")

	res, err := h.service.GetAnalytics(c.Request.Context(), from, to)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
