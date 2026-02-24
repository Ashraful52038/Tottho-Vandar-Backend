package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
}

func NewPostHandler(postUsecase usecase.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

// CreatePost handler
func (h *PostHandler) CreatePost(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	var req domain.CreatePostRequest
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

	post, err := h.postUsecase.Create(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, post)
}

// UpdatePost handler
func (h *PostHandler) UpdatePost(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid post id",
		})
	}

	var req domain.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	post, err := h.postUsecase.Update(c.Request().Context(), uint(postID), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, post)
}

// DeletePost handler
func (h *PostHandler) DeletePost(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid post id",
		})
	}

	err = h.postUsecase.Delete(c.Request().Context(), uint(postID), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Post deleted successfully",
	})
}

// GetPostByID handler
func (h *PostHandler) GetPostByID(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid post id",
		})
	}

	post, err := h.postUsecase.GetByID(c.Request().Context(), uint(postID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "post not found",
		})
	}

	return c.JSON(http.StatusOK, post)
}

// GetAllPosts handler
func (h *PostHandler) GetAllPosts(c echo.Context) error {
	// Parse tag IDs for filtering
	var tagIDs []uint
	tagParam := c.QueryParam("tagIds")
	if tagParam != "" {
		for _, idStr := range strings.Split(tagParam, ",") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
			if err == nil {
				tagIDs = append(tagIDs, uint(id))
			}
		}
	}

	// Parse author ID
	var authorID *uint
	if authorParam := c.QueryParam("authorId"); authorParam != "" {
		if id, err := strconv.ParseUint(authorParam, 10, 32); err == nil {
			aid := uint(id)
			authorID = &aid
		}
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	// যদি কোনো ফিল্টার থাকে, তাহলে search ব্যবহার করুন
	if len(tagIDs) > 0 || authorID != nil {
		params := &repository.SearchParams{
			TagIDs:   tagIDs,
			AuthorID: authorID,
			Page:     page,
			Limit:    limit,
		}
		posts, total, err := h.postUsecase.SearchPosts(c.Request().Context(), params)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		postResponses := make([]*domain.PostResponse, len(posts))
		for i, post := range posts {
			postResponses[i] = post.ToResponse()
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"posts": postResponses,
			"total": total,
			"page":  page,
			"limit": limit,
		})
	}

	posts, total, err := h.postUsecase.GetAll(c.Request().Context(), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts": posts,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetMyPosts handler
func (h *PostHandler) GetMyPosts(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	posts, err := h.postUsecase.GetByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) SearchPosts(c echo.Context) error {
	// Parse query parameters
	query := c.QueryParam("q")

	// Parse tag IDs (comma-separated or multiple)
	var tagIDs []uint
	tagParam := c.QueryParam("tagIds")
	if tagParam != "" {
		for _, idStr := range strings.Split(tagParam, ",") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
			if err == nil {
				tagIDs = append(tagIDs, uint(id))
			}
		}
	}

	// Parse author ID
	var authorID *uint
	if authorParam := c.QueryParam("authorId"); authorParam != "" {
		if id, err := strconv.ParseUint(authorParam, 10, 32); err == nil {
			aid := uint(id)
			authorID = &aid
		}
	}

	// Pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 20
	}

	// Create search params
	params := &repository.SearchParams{
		Query:    query,
		TagIDs:   tagIDs,
		AuthorID: authorID,
		Page:     page,
		Limit:    limit,
	}

	// Execute search
	posts, total, err := h.postUsecase.SearchPosts(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	// Convert to response format
	postResponses := make([]*domain.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = post.ToResponse()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts": postResponses,
		"total": total,
		"page":  page,
		"limit": limit,
		"query": query,
	})
}
