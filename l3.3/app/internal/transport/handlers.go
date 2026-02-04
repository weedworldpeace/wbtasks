package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"net/http"
	"strconv"

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

func (h *handlers) CreateComment(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	var req models.CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	comment, err := h.service.CreateComment(c.Request.Context(), req)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *handlers) GetComments(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	parentID := c.Query("parent")
	query := c.Query("query")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	req := models.GetCommentsRequest{
		Page:   page,
		Limit:  limit,
		Query:  query,
		SortBy: sortBy,
		Order:  order,
	}

	if parentID != "" {
		req.ParentID = &parentID
	}

	result, err := h.service.GetComments(c.Request.Context(), req)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get comments",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handlers) DeleteComment(c *ginext.Context) {
	lg := logger.LoggerFromCtx(c.Request.Context()).Lg
	id := c.Param("id")
	if id == "" {
		lg.Error().Err(models.ErrInvalidInput).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Comment ID is required",
		})
		return
	}

	err := h.service.DeleteComment(c.Request.Context(), id)
	if err != nil {
		lg.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
