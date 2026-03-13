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
