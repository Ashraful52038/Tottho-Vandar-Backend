package domain

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	PostID    uint      `json:"postId" gorm:"not null;column:post_id;index"`
	AuthorID  uint      `json:"authorId" gorm:"not null;column:author_id;index"`
	User      User      `json:"author" gorm:"foreignKey:AuthorID"`
	ParentID  *uint     `json:"parentId,omitempty" gorm:"index"`
	Likes     int       `json:"likes" gorm:"default:0"`
	Post      Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Comment) TableName() string {
	return "comments"
}
