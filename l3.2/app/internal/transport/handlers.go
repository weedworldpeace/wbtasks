package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"errors"
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

func (h *handlers) createLink(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	var sh models.ShortenRequest

	err := c.ShouldBindJSON(&sh)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "bad request body"})
		return
	}

	res, err := h.service.CreateLink(&sh)
	if err != nil {
		lg.Error().Err(err).Send()
		if errors.Is(err, models.ErrAlreadyExistURL) {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "url already exists"})
		} else if errors.Is(err, models.ErrBadURL) {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "bad url value"})
		} else {
			c.JSON(http.StatusInternalServerError, ginext.H{"error": "server error"})
		}
	} else {
		c.JSON(http.StatusOK, ginext.H{"original_url": sh.URL, "short_code": res})
	}
}

func (h *handlers) redirect(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	res, err := h.service.Redirect(c.Request)
	if err != nil {
		lg.Error().Err(err).Send()
		if err == models.ErrNonExistURL {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "non exist url"})
		} else {
			c.JSON(http.StatusInternalServerError, ginext.H{"error": "server error"})
		}
	} else {
		c.Redirect(http.StatusPermanentRedirect, res)
	}
}

func (h *handlers) getAnalytics(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	data, err := h.service.GetAnalytics(c.Request.URL.String()[11:])
	if err != nil {
		lg.Error().Err(err).Send()
		if err == models.ErrNonExistURL {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "non exist url"})
		} else {
			c.JSON(http.StatusInternalServerError, ginext.H{"error": "server error"})
		}
	} else {
		c.JSON(http.StatusOK, data)
	}
}
