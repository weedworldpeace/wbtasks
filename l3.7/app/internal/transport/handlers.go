package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

var secretKey = []byte("your-secret-key")

type handlers struct {
	ctx     context.Context
	service ServiceInterface
}

func (h *handlers) jwtChecker(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	co, err := c.Cookie("jwt")
	if err != nil {
		lg.Error().Err(err).Send()
		c.Status(401)
		c.Abort()
	} else {
		tok, err := jwt.Parse(co, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, models.ErrInvalidSigningMethod
			}
			return secretKey, nil
		})
		if err != nil || !tok.Valid {
			lg.Error().Err(err).Send()
			c.Status(401)
			c.Abort()
		} else {
			if claims, ok := tok.Claims.(jwt.MapClaims); ok {
				c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "user_id", claims["user_id"]))
			} else {
				lg.Error().Err(models.ErrUnexpected).Send()
				c.Status(401)
				c.Abort()
			}
		}
	}
}

func (h *handlers) middleware(c *ginext.Context) {
	lg := logger.LoggerFromCtx(h.ctx)
	requestId := uuid.NewString()
	lgWithReqId := lg.LoggerWithRequestId(requestId)

	lgWithReqId.Lg.Info().Str("method", c.Request.Method).Str("url", c.Request.URL.String()).Msg("received request")

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), logger.LoggerKey, lgWithReqId))
}

func (h *handlers) createItem(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	res, err := h.service.CreateItem(c.Request.Context(), item)
	if err != nil {
		lg.Error().Err(err).Send()
		if errors.Is(err, models.ErrNoPermission) {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *handlers) listItems(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	limit := c.Query("limit")
	offset := c.Query("offset")

	res, err := h.service.ListItems(c.Request.Context(), limit, offset)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) getItem(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	id := c.Param("id")

	res, err := h.service.GetItem(c.Request.Context(), id)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) updateItem(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	id := c.Param("id")

	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	item.ID = id

	res, err := h.service.UpdateItem(c.Request.Context(), item)
	if err != nil {
		lg.Error().Err(err).Send()
		if errors.Is(err, models.ErrNoPermission) {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) deleteItem(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	id := c.Param("id")

	err := h.service.DeleteItem(c.Request.Context(), id)
	if err != nil {
		lg.Error().Err(err).Send()
		if errors.Is(err, models.ErrNoPermission) {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error: err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: err.Error(),
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handlers) listHistory(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	limit := c.Query("limit")
	offset := c.Query("offset")

	res, err := h.service.ListHistory(c.Request.Context(), limit, offset)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handlers) getToken(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg

	role := c.Query("role")

	tok, err := h.service.GetToken(c.Request.Context(), role)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.SetCookie("jwt", tok, 600, "/", "", false, true)

	c.Status(204)
}
