package domain

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	PostID    uint      `json:"postId" gorm:"not nul;index"`
	UserID    uint      `json:"userId" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Post      Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
