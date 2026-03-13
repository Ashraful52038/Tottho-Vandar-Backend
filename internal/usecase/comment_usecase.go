package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/email"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type CommentUsecase interface {
	Create(ctx context.Context, userID uint, postID uint, content string) (*domain.Comment, error)
	GetByPostID(ctx context.Context, postID uint) ([]domain.Comment, error)
	GetByID(ctx context.Context, id uint) (*domain.Comment, error)
	Update(ctx context.Context, commentID uint, userID uint, content string) (*domain.Comment, error)
	Delete(ctx context.Context, commentID uint, userID uint) error
}

type commentUsecase struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	emailQueue  *email.EmailQueueService
}

func NewCommentUsecase(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	emailQueue *email.EmailQueueService,
) CommentUsecase {
	return &commentUsecase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		emailQueue:  emailQueue,
	}
}

func (u *commentUsecase) Create(ctx context.Context, userID uint, postID uint, content string) (*domain.Comment, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	if !user.Verified {
		return nil, errors.New("email not verified")
	}

	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil || post == nil {
		return nil, errors.New("post not found")
	}

	comment := &domain.Comment{
		Content:  content,
		PostID:   postID,
		AuthorID: userID,
	}

	err = u.commentRepo.Create(ctx, comment)

	if err == nil && u.emailQueue != nil {
		go u.sendReplyNotification(ctx, post, user, content)
	}

	return comment, err
}

func (u *commentUsecase) sendReplyNotification(ctx context.Context, post *domain.Post, commenter *domain.User, content string) {
	if post.AuthorID == commenter.ID {
		log.Printf("User commented on own post, skipping notification")
		return
	}

	postAuthor, err := u.userRepo.FindByID(ctx, post.AuthorID)
	if err != nil || postAuthor == nil {
		log.Printf("Failed to find post author: %v", err)
		return
	}

	// notification preference চেক করুন (future use)
	// if !postAuthor.NotificationsEnabled {
	//     return
	// }

	log.Printf("Sending reply notification to %s about comment from %s", postAuthor.Email, commenter.Name)

	// EmailQueueService এর মাধ্যমে notification পাঠান
	if u.emailQueue != nil {
		msg := email.EmailMessage{
			Type:      "reply",
			To:        postAuthor.Email,
			Username:  postAuthor.Name,
			Commenter: commenter.Name,
			Content:   content,
			PostID:    post.ID,
		}
		if err := u.emailQueue.PublishEmail(msg); err != nil {
			log.Printf("Failed to publish email notification: %v", err)
		} else {
			log.Printf("Reply notification queued successfully")
		}
	} else {
		log.Printf("Email queue not initialized, skipping notification")
	}
}

func (u *commentUsecase) GetByPostID(ctx context.Context, postID uint) ([]domain.Comment, error) {
	_, err := u.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return u.commentRepo.FindByPostID(ctx, postID)
}

func (u *commentUsecase) Update(ctx context.Context, commentID uint, userID uint, content string) (*domain.Comment, error) {
	comment, err := u.commentRepo.FindByID(ctx, commentID)
	if err != nil || comment == nil {
		return nil, errors.New("comment not found")
	}

	if comment.AuthorID != userID {
		return nil, errors.New("unauthorized")
	}

	comment.Content = content
	err = u.commentRepo.Update(ctx, comment)
	return comment, err
}

func (u *commentUsecase) Delete(ctx context.Context, commentID uint, userID uint) error {
	comment, err := u.commentRepo.FindByID(ctx, commentID)
	if err != nil || comment == nil {
		return errors.New("comment not found")
	}

	if comment.AuthorID != userID {
		return errors.New("unauthorized")
	}

	return u.commentRepo.Delete(ctx, commentID)
}

func (u *commentUsecase) GetByID(ctx context.Context, id uint) (*domain.Comment, error) {
	comment, err := u.commentRepo.FindByID(ctx, id)
	if err != nil || comment == nil {
		return nil, errors.New("comment not found")
	}
	return comment, nil
}
