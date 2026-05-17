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

type LikeHandler struct {
	likeUsecase    usecase.LikeUsecase
	postUsecase    usecase.PostUsecase
	commentUsecase usecase.CommentUsecase
	userUsecase    usecase.UserUsecase
	notifUsecase   *notificationUsecase.NotificationUsecase
	emailSender    *email.MailpitSender
}

func NewLikeHandler(
	likeUsecase usecase.LikeUsecase,
	postUsecase usecase.PostUsecase,
	commentUsecase usecase.CommentUsecase,
	userUsecase usecase.UserUsecase,
	notifUsecase *notificationUsecase.NotificationUsecase,
	emailSender *email.MailpitSender,
) *LikeHandler {
	return &LikeHandler{
		likeUsecase:    likeUsecase,
		postUsecase:    postUsecase,
		commentUsecase: commentUsecase,
		userUsecase:    userUsecase,
		notifUsecase:   notifUsecase,
		emailSender:    emailSender,
	}
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

	post, err := h.postUsecase.GetByID(c.Request().Context(), uint(postID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch updated post"})
	}

	// Get liker info once
	liker, _ := h.userUsecase.GetByID(c.Request().Context(), userID)
	likerName := "Someone"
	if liker != nil {
		likerName = liker.Name
	}

	if like != nil && h.notifUsecase != nil && post.AuthorID != userID {
		log.Printf("🔔 [DEBUG] Sending like notification: author=%d, liker=%s, post=%d",
			post.AuthorID, likerName, postID)

		// WebSocket notification
		h.notifUsecase.NotifyNewLike(post.AuthorID, uint(postID), likerName)

		// Like handler এর TogglePost ফাংশনে, notification এর পর এই কোড যোগ করো:

		// ✅ ইমেইল সেন্ড করার আগে লগ
		log.Printf("📧 [EMAIL CHECK] emailSender=%v, post.AuthorID=%d", h.emailSender != nil, post.AuthorID)

		if h.emailSender != nil {
			author, err := h.userUsecase.GetByID(c.Request().Context(), post.AuthorID)
			log.Printf("📧 [EMAIL CHECK] author found: err=%v, author=%v, email=%s", err, author != nil, func() string {
				if author != nil {
					return author.Email
				}
				return ""
			}())

			if author != nil && author.Email != "" {
				liker, _ := h.userUsecase.GetByID(c.Request().Context(), userID)
				likerName := "Someone"
				if liker != nil {
					likerName = liker.Name
				}

				log.Printf("📧 [EMAIL SENDING] to=%s, subject=❤️ %s liked your post", author.Email, likerName)

				go func() {
					subject := fmt.Sprintf("❤️ %s liked your post", likerName)
					body := fmt.Sprintf(`
                <h2>New Like on Your Post</h2>
                <p><strong>%s</strong> liked your post "<strong>%s</strong>"</p>
                <a href="http://localhost:3000/posts/%d">View Post</a>
            `, likerName, post.Title, post.ID)

					err := h.emailSender.SendEmail(author.Email, subject, body)
					if err != nil {
						log.Printf("❌ [EMAIL FAILED] %v", err)
					} else {
						log.Printf("✅ [EMAIL SUCCESS] Sent to %s", author.Email)
					}
				}()
			} else {
				log.Printf("⚠️ [EMAIL SKIP] No email for author %d", post.AuthorID)
			}
		} else {
			log.Printf("⚠️ [EMAIL SKIP] emailSender is nil")
		}
	} else {
		log.Printf("⚠️ [DEBUG] Notification NOT sent: like=%v, notifUsecase=%v, sameUser=%v",
			like != nil, h.notifUsecase != nil, post.AuthorID == userID)
	}

	response := post.ToResponse()
	response.IsLiked = (like != nil)
	return c.JSON(http.StatusOK, response)
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

	comment, err := h.commentUsecase.GetByID(c.Request().Context(), uint(commentID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch comment"})
	}

	if like == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "unliked", "comment": comment})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "liked", "like": like, "comment": comment})
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
