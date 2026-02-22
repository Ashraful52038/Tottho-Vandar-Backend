package domain

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	PostID    uint      `json:"postId" gorm:"not null"`
	AuthorID  uint      `json:"authorId" gorm:"not null"`
	Author    User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
