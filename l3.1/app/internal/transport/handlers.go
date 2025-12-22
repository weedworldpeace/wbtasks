package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/retry"
)

type handlers struct {
	ctx               context.Context
	service           ServiceInterface
	emailSenderConfig EmailSenderConfig
}

func (h *handlers) middleware(c *ginext.Context) {
	lg := logger.LoggerFromCtx(h.ctx)
	requestId := uuid.NewString()
	lgWithReqId := lg.LoggerWithRequestId(requestId)

	lgWithReqId.Lg.Info().Str("method", c.Request.Method).Str("url", c.Request.URL.String()).Msg("received request")

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), logger.LoggerKey, lgWithReqId))
}

func (h *handlers) createNotification(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	notif := models.NewNotification()
	if err := c.ShouldBindJSON(notif); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "bad request"})
		return
	}

	if id, err := h.service.CreateNotification(notif); err != nil {
		lg.Error().Err(err).Send()

		c.JSON(http.StatusInternalServerError, ginext.H{"error": "server error"})
	} else {
		c.JSON(http.StatusOK, ginext.H{"id": id})
	}
}

func (h *handlers) readNotification(c *ginext.Context) {
	id := c.Param(idParam)

	if notif, err := h.service.ReadNotification(id); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, notif)
	}
}

func (h *handlers) deleteNotification(c *ginext.Context) {
	id := c.Param(idParam)

	if err := h.service.DeleteNotification(id); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, ginext.H{"message": "notification deleted"})
	}
}

func (h *handlers) consume(ch chan []byte) {
	lg := logger.LoggerFromCtx(h.ctx).Lg
	for d := range ch {
		var notif models.Notification
		if err := json.Unmarshal(d, &notif); err != nil {
			lg.Error().Err(err).Msg("while processing on rabbit consumer")
			continue
		}

		if _, err := h.service.ReadNotification(notif.Id); err != nil {
			continue
		}

		if err := retry.Do(func() error { return h.sendEmail(notif.Email, d) }, retry.Strategy{Attempts: 10, Delay: 3, Backoff: 3}); err != nil {
			lg.Error().Err(err).Msg("while sending an email")
		}

		err := h.service.DeleteNotification(notif.Id)
		if err != nil {
			lg.Error().Err(err).Msg("while deleting notification after sending email")
		}
	}
}

func (h *handlers) sendEmail(to string, bd []byte) error {
	subject := "Notification"
	body := string(bd)
	message := []byte(subject + "\n" + body)

	// auth := smtp.PlainAuth("", "", "", h.emailSenderConfig.Host)

	return smtp.SendMail(fmt.Sprintf("%s:%s", h.emailSenderConfig.Host, h.emailSenderConfig.Port), nil, h.emailSenderConfig.From, []string{to}, message)
}
