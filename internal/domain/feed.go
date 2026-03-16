package domain

import (
	"time"
)

type FeedQueryParams struct {
	Page      int    `json:"page" query:"page"`
	Limit     int    `json:"limit" query:"limit"`
	SortBy    string `json:"sortBy" query:"sort_by"`       // created_at, likes_count, comments_count
	SortOrder string `json:"sortOrder" query:"sort_order"` // asc, desc
}

type FeedPost struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      uint      `json:"authorId"`
	AuthorName    string    `json:"authorName"`
	FeaturedImage string    `json:"featuredImage"`
	CoverImage    string    `json:"coverImage"`
	AuthorEmail   string    `json:"authorEmail"`
	CreatedAt     time.Time `json:"createdAt"`
	LikesCount    int64     `json:"likesCount"`
	CommentsCount int64     `json:"commentsCount"`
	Tags          []string  `json:"tags"`
	IsLiked       bool      `json:"isLiked"`
	IsBookmarked  bool      `json:"isBookmarked"`
	Published     bool      `json:"published"`
}

// FeedResponse - ফিড রেসপন্স
type FeedResponse struct {
	Posts      []FeedPost `json:"posts"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"totalPages"`
}

type FollowTagRequest struct {
	TagID uint `json:"tagId" validate:"required"`
}

type UserFollowedTag struct {
	UserID    uint      `json:"userId" gorm:"primaryKey"`
	TagID     uint      `json:"tagId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserFollowedTag) TableName() string {
	return "user_followed_tags"
}
