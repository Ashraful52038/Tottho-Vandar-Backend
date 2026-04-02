package handler

import (
	"net/http"
	"strconv"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type FeedHandler struct {
	postUsecase usecase.PostUsecase
}

func NewFeedHandler(postUsecase usecase.PostUsecase) *FeedHandler {
	return &FeedHandler{
		postUsecase: postUsecase,
	}
}

func (h *FeedHandler) GetPersonalizedFeed(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	// Parse query parameters
	params := &domain.FeedQueryParams{
		Page:      1,
		Limit:     20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}

	if limit := c.QueryParam("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			params.Limit = l
		}
	}

	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}

	if sortOrder := c.QueryParam("sort_order"); sortOrder != "" {
		params.SortOrder = sortOrder
	}

	feed, err := h.postUsecase.GetPersonalizedFeed(c.Request().Context(), userID, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, feed)
}

func (h *FeedHandler) GetPublicFeed(c echo.Context) error {
	// Parse query parameters
	params := &domain.FeedQueryParams{
		Page:      1,
		Limit:     20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}

	if limit := c.QueryParam("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			params.Limit = l
		}
	}

	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}

	if sortOrder := c.QueryParam("sort_order"); sortOrder != "" {
		params.SortOrder = sortOrder
	}

	feed, err := h.postUsecase.GetPublicFeed(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, feed)
}

func (h *FeedHandler) FollowTag(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid tag id",
		})
	}

	if err := h.postUsecase.FollowTag(c.Request().Context(), userID, uint(tagID)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "tag followed successfully",
	})
}

func (h *FeedHandler) UnfollowTag(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid tag id",
		})
	}

	if err := h.postUsecase.UnfollowTag(c.Request().Context(), userID, uint(tagID)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "tag unfollowed successfully",
	})
}

func (h *FeedHandler) GetFollowedTags(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	tags, err := h.postUsecase.GetFollowedTags(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tags)
}
