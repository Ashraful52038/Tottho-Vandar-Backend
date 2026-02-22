package domain

import (
	"time"
)

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	AuthorID  uint      `json:"authorId" gorm:"not null"`
	Author    User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Tags      []Tag     `json:"tags" gorm:"many2many:post_tags;"`
	Likes     int       `json:"likes" gorm:"default:0"`
	Comments  int       `json:"comments" gorm:"default:0"`
	Published bool      `json:"published" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreatePostRequest - পোস্ট তৈরি করার জন্য রিকোয়েস্ট
type CreatePostRequest struct {
	Title     string   `json:"title" validate:"required"`
	Content   string   `json:"content" validate:"required"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}

// UpdatePostRequest - পোস্ট আপডেট করার জন্য রিকোয়েস্ট
type UpdatePostRequest struct {
	Title     *string  `json:"title,omitempty"`
	Content   *string  `json:"content,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Published *bool    `json:"published,omitempty"`
}
