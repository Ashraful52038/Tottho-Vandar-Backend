package handler

import (
	"net/http"
	"strconv"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	commentUsecase usecase.CommentUsecase
	likeUsecase    usecase.LikeUsecase
}

func NewCommentHandler(commentUsecase usecase.CommentUsecase, likeUsecase usecase.LikeUsecase) *CommentHandler {
	return &CommentHandler{commentUsecase: commentUsecase, likeUsecase: likeUsecase}
}

func (h *CommentHandler) Create(c echo.Context) error {
	userID := c.Get("userID").(uint)
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid post id"})
	}

	var req struct {
		Content string `json:"content" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	comment, err := h.commentUsecase.Create(c.Request().Context(), userID, uint(postID), req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetByPost(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid post id"})
	}

	comments, err := h.commentUsecase.GetByPostID(c.Request().Context(), uint(postID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	userID, _ := c.Get("userID").(uint)

	var response []map[string]interface{}
	for _, comment := range comments {
		commentMap := map[string]interface{}{
			"id":       comment.ID,
			"content":  comment.Content,
			"postId":   comment.PostID,
			"parentId": comment.ParentID,
			"author": map[string]interface{}{
				"id":     comment.User.ID,
				"name":   comment.User.Name,
				"avatar": comment.User.Avatar,
			},
			"likes":     comment.Likes,
			"createdAt": comment.CreatedAt,
			"updatedAt": comment.UpdatedAt,
		}

		if userID > 0 {
			isLiked, _ := h.likeUsecase.CheckUserLikedComment(c.Request().Context(), userID, comment.ID)
			commentMap["isLiked"] = isLiked
		} else {
			commentMap["isLiked"] = false
		}

		response = append(response, commentMap)
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CommentHandler) Update(c echo.Context) error {
	userID := c.Get("userID").(uint)
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid comment id"})
	}

	var req struct {
		Content string `json:"content" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	comment, err := h.commentUsecase.Update(c.Request().Context(), uint(commentID), userID, req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) Delete(c echo.Context) error {
	userID := c.Get("userID").(uint)
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid comment id"})
	}

	err = h.commentUsecase.Delete(c.Request().Context(), uint(commentID), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "comment deleted"})
}
