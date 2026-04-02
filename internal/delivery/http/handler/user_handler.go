package handler

import (
	"net/http"
	"strconv"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
	likeUsecase usecase.LikeUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase, likeUsecase usecase.LikeUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		likeUsecase: likeUsecase,
	}
}

// GetCurrentUser handler
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	user, err := h.userUsecase.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser handler
func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	var req domain.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	user, err := h.userUsecase.Update(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// Logout handler
func (h *UserHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Logged out successfully",
	})
}

// GetUserByID handler
func (h *UserHandler) GetUserByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	user, err := h.userUsecase.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "user not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserProfile(c echo.Context) error {
	userId := c.Param("userId")

	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	user, err := h.userUsecase.GetByID(c.Request().Context(), uint(userID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "user not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) validateAndGetUser(c echo.Context) (uint, error) {
	userId := c.Param("userId")
	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	_, err = h.userUsecase.GetByID(c.Request().Context(), uint(userID))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return uint(userID), nil
}

func (h *UserHandler) GetUserPosts(c echo.Context) error {
	_, err := h.validateAndGetUser(c)
	if err != nil {
		return err
	}

	userId := c.Param("userId")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	posts, total, err := h.userUsecase.GetUserPosts(c.Request().Context(), userId, page, limit)
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

func (h *UserHandler) GetUserComments(c echo.Context) error {
	_, err := h.validateAndGetUser(c)
	if err != nil {
		return err
	}

	userId := c.Param("userId")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	comments, total, err := h.userUsecase.GetUserComments(c.Request().Context(), userId, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"comments": comments,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (h *UserHandler) GetUserLikes(c echo.Context) error {
	userId := c.Param("userId")
	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	likes, total, err := h.likeUsecase.GetUserLikes(c.Request().Context(), uint(userID), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	responses := make([]domain.LikeResponse, len(likes))
	for i, like := range likes {
		responses[i] = domain.LikeResponse{
			ID:        like.ID,
			CreatedAt: like.CreatedAt,
			Post: domain.LikedPostResponse{
				ID:            like.Post.ID,
				Title:         like.Post.Title,
				Content:       like.Post.Content,
				Author:        like.Post.Author,
				FeaturedImage: like.Post.FeaturedImage,
				CoverImage:    like.Post.CoverImage,
				CreatedAt:     like.Post.CreatedAt,
				Likes:         like.Post.Likes,
				Comments:      like.Post.Comments,
			},
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"likes": responses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *UserHandler) GetUserFollowers(c echo.Context) error {

	_, err := h.validateAndGetUser(c)
	if err != nil {
		return err
	}

	userId := c.Param("userId")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	currentUserID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	followers, total, err := h.userUsecase.GetFollowers(c.Request().Context(), userId, currentUserID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"followers": followers,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

func (h *UserHandler) GetUserFollowing(c echo.Context) error {
	_, err := h.validateAndGetUser(c)
	if err != nil {
		return err
	}

	userId := c.Param("userId")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	// current user ID
	currentUserID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	following, total, err := h.userUsecase.GetFollowing(c.Request().Context(), userId, currentUserID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"following": following,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

func (h *UserHandler) GetFollowStatus(c echo.Context) error {
	userId := c.Param("userId")

	//Take current Id from JWT
	currentUserID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	targetUserID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	isFollowing, err := h.userUsecase.CheckFollowStatus(c.Request().Context(), currentUserID, uint(targetUserID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"isFollowing": isFollowing,
	})
}

func (h *UserHandler) FollowUser(c echo.Context) error {
	userId := c.Param("userId")

	currentUserID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	targetUserID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	err = h.userUsecase.FollowUser(c.Request().Context(), currentUserID, uint(targetUserID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "followed successfully",
		"following": true,
	})
}

func (h *UserHandler) UnfollowUser(c echo.Context) error {
	userId := c.Param("userId")

	currentUserID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user not authenticated",
		})
	}

	targetUserID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	err = h.userUsecase.UnfollowUser(c.Request().Context(), currentUserID, uint(targetUserID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "unfollowed successfully",
		"following": false,
	})
}
