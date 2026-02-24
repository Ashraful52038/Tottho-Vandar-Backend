package handler

import (
	"net/http"
	"strconv"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	tagUsecase usecase.TagUsecase
}

func NewTagHandler(tagUsecase usecase.TagUsecase) *TagHandler {
	return &TagHandler{
		tagUsecase: tagUsecase,
	}
}

// CreateTag handler - নতুন ট্যাগ তৈরি (Admin only)
func (h *TagHandler) CreateTag(c echo.Context) error {
	// Optional: Check if user is admin
	userRole, ok := c.Get("userRole").(string)
	if !ok || userRole != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error": "admin access required",
		})
	}

	var req domain.CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	tag, err := h.tagUsecase.Create(c.Request().Context(), req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, tag)
}

// GetTagByID handler - নির্দিষ্ট ট্যাগ দেখুন
func (h *TagHandler) GetTagByID(c echo.Context) error {
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid tag id",
		})
	}

	tag, err := h.tagUsecase.GetByID(c.Request().Context(), uint(tagID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "tag not found",
		})
	}

	return c.JSON(http.StatusOK, tag)
}

// GetTagBySlug handler - স্লাগ দিয়ে ট্যাগ খুঁজুন
func (h *TagHandler) GetTagBySlug(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "slug is required",
		})
	}

	tag, err := h.tagUsecase.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "tag not found",
		})
	}

	return c.JSON(http.StatusOK, tag)
}

// GetAllTags handler - সব ট্যাগের তালিকা (পেজিনেটেড)
func (h *TagHandler) GetAllTags(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 20
	}

	tags, total, err := h.tagUsecase.GetAll(c.Request().Context(), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tags":  tags,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *TagHandler) GetPopularTags(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	tags, err := h.tagUsecase.GetPopular(c.Request().Context(), limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) GetTagsByPostID(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid post id",
		})
	}

	tags, err := h.tagUsecase.GetByPostID(c.Request().Context(), uint(postID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) UpdateTag(c echo.Context) error {
	// Optional: Check if user is admin
	userRole, ok := c.Get("userRole").(string)
	if !ok || userRole != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error": "admin access required",
		})
	}

	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid tag id",
		})
	}

	var req domain.UpdateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	tag, err := h.tagUsecase.Update(c.Request().Context(), uint(tagID), req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) DeleteTag(c echo.Context) error {
	// Optional: Check if user is admin
	userRole, ok := c.Get("userRole").(string)
	if !ok || userRole != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error": "admin access required",
		})
	}

	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid tag id",
		})
	}

	err = h.tagUsecase.Delete(c.Request().Context(), uint(tagID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Tag deleted successfully",
	})
}
