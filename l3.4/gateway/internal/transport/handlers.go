package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

var validTypes = map[string]bool{
	"image/jpeg": true,
	// "image/png":  true,
	// "image/gif":  true,
	// "image/webp": true,
}

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

func (h *handlers) UploadImage(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	file, _, err := c.Request.FormFile("image")
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file from form"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}

	ext := http.DetectContentType(data)

	if !isValidImageType(ext) {
		lg.Error().Err(models.ErrInvalidImageType).Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image type"})
		return
	}

	res, err := h.service.UploadImage(&models.Image{Raw: data, Ext: ext, Size: len(data)})
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, res)
	lg.Debug().Str("task_id", res.ID).Msg("image uploaded successfully")
}

func (h *handlers) GetImage(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	id := c.Param("id")

	data, err := h.service.GetImage(id)
	if err != nil {
		lg.Error().Err(err).Send()
		if errors.Is(err, models.ErrImagePending) {
			c.JSON(http.StatusNotFound, ginext.H{"error": "pending"})
			return
		} else {
			c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, data)
	lg.Debug().Str("id", id).Msg("image served")
}

func (h *handlers) DeleteImage(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	id := c.Param("id")

	err := h.service.DeleteImage(id)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
	lg.Debug().Str("id", id).Msg("image deleted successfully")
}

func isValidImageType(contentType string) bool {
	return validTypes[contentType]
}
