package domain

import (
	"time"
)

// FeedQueryParams - ফিডের ক্যোয়ারী প্যারামিটার
type FeedQueryParams struct {
	Page      int    `json:"page" query:"page"`
	Limit     int    `json:"limit" query:"limit"`
	SortBy    string `json:"sortBy" query:"sort_by"`       // created_at, likes_count, comments_count
	SortOrder string `json:"sortOrder" query:"sort_order"` // asc, desc
}

// FeedPost - ফিডে দেখানোর জন্য পোস্টের স্ট্রাকচার
type FeedPost struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      uint      `json:"authorId"`
	AuthorName    string    `json:"authorName"`
	AuthorEmail   string    `json:"authorEmail"`
	CreatedAt     time.Time `json:"createdAt"`
	LikesCount    int64     `json:"likesCount"`
	CommentsCount int64     `json:"commentsCount"`
	Tags          []string  `json:"tags"`
	IsLiked       bool      `json:"isLiked"`      // বর্তমান ইউজার লাইক করেছে কিনা
	IsBookmarked  bool      `json:"isBookmarked"` // বর্তমান ইউজার বুকমার্ক করেছে কিনা
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

// FollowTagRequest - ট্যাগ ফলো/আনফলো রিকোয়েস্ট
type FollowTagRequest struct {
	TagID uint `json:"tagId" validate:"required"`
}

// UserFollowedTag - ইউজারের ফলো করা ট্যাগের মডেল
type UserFollowedTag struct {
	UserID    uint      `json:"userId" gorm:"primaryKey"`
	TagID     uint      `json:"tagId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName - টেবিলের নাম নির্ধারণ
func (UserFollowedTag) TableName() string {
	return "user_followed_tags"
}
