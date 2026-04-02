// internal/domain/follow.go
package domain

import "time"

type Follow struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FollowerID  uint      `json:"followerId" gorm:"not null;index:idx_follow_follower"`
	FollowingID uint      `json:"followingId" gorm:"not null;index:idx_follow_following"`
	CreatedAt   time.Time `json:"createdAt"`

	// Associations
	Follower  *User `json:"follower,omitempty" gorm:"foreignKey:FollowerID"`
	Following *User `json:"following,omitempty" gorm:"foreignKey:FollowingID"`
}

type FollowUser struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar,omitempty"`
	Bio         string `json:"bio,omitempty"`
	IsFollowing bool   `json:"isFollowing"`
}

func (Follow) TableName() string {
	return "follows"
}
