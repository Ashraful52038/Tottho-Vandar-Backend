package domain

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null;size:50"`
	Posts     []Post         `json:"posts,omitempty" gorm:"many2many:post_tags;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateTagRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UpdateTagRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type TagResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	PostsCount int64     `json:"postsCount,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (t *Tag) ToResponse() *TagResponse {
	return &TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	t.Slug = generateSlug(t.Name)
	return nil
}

func (t *Tag) BeforeUpdate(tx *gorm.DB) error {
	t.Slug = generateSlug(t.Name)
	return nil
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	slug = strings.ReplaceAll(slug, ",", "")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, `"`, "")
	return slug
}

type PostTag struct {
	PostID    uint      `json:"postId" gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	TagID     uint      `json:"tagId" gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `json:"createdAt"`

	Post Post `json:"-" gorm:"foreignKey:PostID"`
	Tag  Tag  `json:"-" gorm:"foreignKey:TagID"`
}

func (PostTag) TableName() string {
	return "post_tags"
}
