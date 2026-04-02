package domain

import (
	"time"
)

type Like struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"userId" gorm:"not null;index"`
	PostID    *uint      `json:"postId" gorm:"index"`
	CommentID *uint      `json:"commentId,omitempty" gorm:"index;uniqueIndex:idx_user_comment"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-" gorm:"index"`

	// Associations
	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Post    Post     `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Comment *Comment `json:"comment,omitempty" gorm:"foreignKey:CommentID"`
}

type LikedPostResponse struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        User      `json:"author"`
	FeaturedImage string    `json:"featuredImage"`
	CoverImage    string    `json:"coverImage"`
	CreatedAt     time.Time `json:"createdAt"`
	Likes         int       `json:"likes"`
	Comments      int       `json:"comments"`
}

type LikeResponse struct {
	ID        uint              `json:"id"`
	Post      LikedPostResponse `json:"post"`
	CreatedAt time.Time         `json:"createdAt"`
}
