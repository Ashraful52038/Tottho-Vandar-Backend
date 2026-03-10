package domain

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Title         string         `json:"title" gorm:"not null"`
	Content       string         `json:"content" gorm:"type:text;not null"`
	AuthorID      uint           `json:"authorId" gorm:"not null"`
	Author        User           `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	FeaturedImage string         `json:"featuredImage" gorm:"size:500"`
	CoverImage    string         `json:"coverImage" gorm:"size:500"`
	Tags          []Tag          `json:"tags" gorm:"many2many:post_tags;constraint:OnDelete:CASCADE;"`
	Likes         int            `json:"likes" gorm:"default:0"`
	Comments      int            `json:"comments" gorm:"default:0"`
	Published     bool           `json:"published" gorm:"default:false"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TagID type for type safety
type TagID uint

// CreatePostRequest
type CreatePostRequest struct {
	Title         string   `json:"title" validate:"required,min=3,max=200"`
	Content       string   `json:"content" validate:"required,min=10"`
	TagIDs        []uint   `json:"tagIds,omitempty"` // Explicit []uint instead of interface{}
	FeaturedImage string   `json:"featuredImage"`
	CoverImage    string   `json:"coverImage"`
	TagNames      []string `json:"tagNames" binding:"required,min=1"`
	Published     bool     `json:"published"`
}

// UpdatePostRequest
type UpdatePostRequest struct {
	Title         *string  `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Content       *string  `json:"content,omitempty" validate:"omitempty,min=10"`
	TagIDs        []uint   `json:"tagIds,omitempty"`        // Explicit []uint instead of interface{}
	FeaturedImage *string  `json:"featuredImage,omitempty"` // ✅ যোগ করুন
	CoverImage    *string  `json:"coverImage,omitempty"`
	TagNames      []string `json:"tagNames,omitempty"`
	Published     *bool    `json:"published,omitempty"`
}

// ToResponse - converts Post to response format
func (p *Post) ToResponse() *PostResponse {
	// Prepare author response if author exists
	var authorResponse *UserResponse
	if p.Author.ID != 0 {
		authorResponse = p.Author.ToResponse()
	}

	return &PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  p.AuthorID,
		Author:    authorResponse,
		Tags:      p.Tags,
		Likes:     p.Likes,
		Comments:  p.Comments,
		Published: p.Published,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		IsLiked:   false,
	}
}

// PostResponse - response struct for posts
type PostResponse struct {
	ID            uint          `json:"id"`
	Title         string        `json:"title"`
	Content       string        `json:"content"`
	AuthorID      uint          `json:"authorId"`
	Author        *UserResponse `json:"author,omitempty"`
	FeaturedImage string        `json:"featuredImage"` // ✅ যোগ করুন
	CoverImage    string        `json:"coverImage"`
	Tags          []Tag         `json:"tags"`
	Likes         int           `json:"likes"`
	Comments      int           `json:"comments"`
	Published     bool          `json:"published"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	IsLiked       bool          `json:"isLiked"`
}

// PostListResponse - paginated post list response
type PostListResponse struct {
	Posts      []PostResponse `json:"posts"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"totalPages"`
}
