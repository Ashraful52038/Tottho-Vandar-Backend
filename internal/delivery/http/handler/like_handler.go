package handler

import (
	"net/http"
	"strconv"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type LikeHandler struct {
	likeUsecase usecase.LikeUsecase
}

func NewLikeHandler(likeUsecase usecase.LikeUsecase) *LikeHandler {
	return &LikeHandler{likeUsecase: likeUsecase}
}

// TogglePost
func (h *LikeHandler) TogglePost(c echo.Context) error {
	userID := c.Get("userID").(uint)

	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid post id"})
	}

	like, err := h.likeUsecase.TogglePostLike(c.Request().Context(), userID, uint(postID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if like == nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "unliked"})
	}
	return c.JSON(http.StatusOK, like)
}

// GetPostLikes
func (h *LikeHandler) GetPostLikes(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid post id"})
	}

	likes, err := h.likeUsecase.GetPostLikes(c.Request().Context(), uint(postID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	count, _ := h.likeUsecase.GetPostLikesCount(c.Request().Context(), uint(postID))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"likes": likes,
		"count": count,
	})
}

// ToggleComment
func (h *LikeHandler) ToggleComment(c echo.Context) error {
	userID := c.Get("userID").(uint)

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid comment id"})
	}

	like, err := h.likeUsecase.ToggleCommentLike(c.Request().Context(), userID, uint(commentID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if like == nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "unliked"})
	}
	return c.JSON(http.StatusOK, like)
}

// GetCommentLikes
func (h *LikeHandler) GetCommentLikes(c echo.Context) error {
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid comment id"})
	}

	likes, err := h.likeUsecase.GetCommentLikes(c.Request().Context(), uint(commentID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	count, _ := h.likeUsecase.GetCommentLikesCount(c.Request().Context(), uint(commentID))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"likes": likes,
		"count": count,
	})
}
