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
	logCh   chan models.ToLog
	notifCh chan models.EventTask
}

func (h *handlers) middleware(c *ginext.Context) {
	lg := logger.LoggerFromCtx(h.ctx).Lg
	requestId := uuid.NewString()

	lg.Info().Str(string(logger.RequestIdKey), requestId).Str("method", c.Request.Method).Str("url", c.Request.URL.String()).Msg("received request")

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), logger.RequestIdKey, requestId))
}

func (h *handlers) createEvent(c *ginext.Context) {
	var event models.UserEvent

	em := c.Query("email")

	if err := c.ShouldBindJSON(&event); err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	res, err := h.service.CreateEvent(c.Request.Context(), event, em)
	if err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	h.notifCh <- models.EventTask{Event: *res, Email: em}

	c.JSON(http.StatusCreated, res)
}

func (h *handlers) updateEvent(c *ginext.Context) {
	var event models.UserEvent

	if err := c.ShouldBindJSON(&event); err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	res, err := h.service.UpdateEvent(c.Request.Context(), event)
	if err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) deleteEvent(c *ginext.Context) {
	var event models.UserEvent

	if err := c.ShouldBindJSON(&event); err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	err := h.service.DeleteEvent(c.Request.Context(), event)
	if err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handlers) eventsForDay(c *ginext.Context) {
	userId := c.Query("user_id")
	date := c.Query("date")

	events, err := h.service.ReadEvents(c, userId, date, "day")
	if err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *handlers) eventsForWeek(c *ginext.Context) {
	userId := c.Query("user_id")
	date := c.Query("date")

	events, err := h.service.ReadEvents(c, userId, date, "week")
	if err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *handlers) eventsForMonth(c *ginext.Context) {
	userId := c.Query("user_id")
	date := c.Query("date")

	events, err := h.service.ReadEvents(c, userId, date, "month")
	if err != nil {
		h.logCh <- models.ToLog{Level: logger.ErrLevelKey, Error: err, Ctx: c.Request.Context()}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, events)
}
