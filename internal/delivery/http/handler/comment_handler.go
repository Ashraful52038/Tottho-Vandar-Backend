package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/email"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase"
	notificationUsecase "github.com/Ashraful52038/tottho-vandar-backend/internal/usecase/notification"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	commentUsecase usecase.CommentUsecase
	likeUsecase    usecase.LikeUsecase
	postUsecase    usecase.PostUsecase
	userUsecase    usecase.UserUsecase
	notifUsecase   *notificationUsecase.NotificationUsecase
	emailSender    *email.MailpitSender
}

func NewCommentHandler(
	commentUsecase usecase.CommentUsecase,
	likeUsecase usecase.LikeUsecase,
	postUsecase usecase.PostUsecase,
	userUsecase usecase.UserUsecase,
	notifUsecase *notificationUsecase.NotificationUsecase,
	emailSender *email.MailpitSender) *CommentHandler {
	return &CommentHandler{
		commentUsecase: commentUsecase,
		likeUsecase:    likeUsecase,
		postUsecase:    postUsecase,
		userUsecase:    userUsecase,
		notifUsecase:   notifUsecase,
		emailSender:    emailSender,
	}
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

	// ✅ DEBUG LOG 1
	log.Printf("🔍 [COMMENT DEBUG] Comment created: userID=%d, postID=%d, notifUsecase=%v",
		userID, postID, h.notifUsecase != nil)

	if comment != nil && h.notifUsecase != nil {
		// Get post to find author
		post, err := h.postUsecase.GetByID(c.Request().Context(), uint(postID))

		// ✅ DEBUG LOG 2
		log.Printf("🔍 [COMMENT DEBUG] Post found: err=%v, post=%v, post.AuthorID=%d, userID=%d",
			err, post != nil, func() uint {
				if post != nil {
					return post.AuthorID
				}
				return 0
			}(), userID)

		if err == nil && post != nil && post.AuthorID != userID {
			// Get commenter info
			currentUser, err := h.userUsecase.GetByID(c.Request().Context(), userID)
			userName := "Someone"
			if err == nil && currentUser != nil {
				userName = currentUser.Name
			}

			// ✅ DEBUG LOG 3
			log.Printf("🔔 [COMMENT TRIGGER] Sending notification: author=%d, commenter=%s, post=%d",
				post.AuthorID, userName, postID)

			// Send notification to post author
			h.notifUsecase.NotifyNewComment(post.AuthorID, uint(postID), userName, req.Content)

			// ✅ Email notification
			if h.emailSender != nil {
				author, err := h.userUsecase.GetByID(c.Request().Context(), post.AuthorID)
				if err == nil && author != nil && author.Email != "" {
					postURL := fmt.Sprintf("http://localhost:3000/posts/%d", post.ID)

					// Get comment excerpt
					commentExcerpt := req.Content
					if len(commentExcerpt) > 100 {
						commentExcerpt = commentExcerpt[:100] + "..."
					}

					emailBody := fmt.Sprintf(email.CommentEmailTemplate,
						userName,        // %s - commenter name
						post.Title,      // %s - post title
						commentExcerpt,  // %s - comment excerpt
						post.Comments+1, // %d - total comments
						post.Likes,      // %d - total likes
						postURL,         // %s - post URL
						userName,        // %s - commenter name for footer
					)

					go func() {
						subject := fmt.Sprintf("💬 %s commented on your post \"%s\"", userName, post.Title)
						if err := h.emailSender.SendEmail(author.Email, subject, emailBody); err != nil {
							log.Printf("❌ [EMAIL FAILED] Comment email: %v", err)
						} else {
							log.Printf("✅ [EMAIL SUCCESS] Comment email sent to %s", author.Email)
						}
					}()
				}
			}
		} else {
			// ✅ DEBUG LOG 4 - কেন notification send হচ্ছে না
			log.Printf("⚠️ [COMMENT SKIP] Reason: err=%v, postNil=%v, sameUser=%v",
				err, post == nil, post != nil && post.AuthorID == userID)
		}
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
