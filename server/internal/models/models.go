package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"` // 不返回给前端
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Role      string         `json:"role" gorm:"size:20;default:'user'"` // admin, user
	Bio       string         `json:"bio" gorm:"size:500"`
	Comments  []Comment      `json:"-" gorm:"foreignKey:UserID"`
	Likes     []Like         `json:"-" gorm:"foreignKey:UserID"`
}

// Comment 评论模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	PostID    string         `json:"post_id" gorm:"index;size:100;not null"` // 文章 ID（slug）
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	ParentID  *uint          `json:"parent_id" gorm:"index"` // 父评论 ID（用于嵌套回复）
	Content   string         `json:"content" gorm:"type:text;not null"`
	Status    string         `json:"status" gorm:"size:20;default:'approved'"` // pending, approved, spam
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Replies   []Comment      `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// Like 点赞模型
type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	PostID    string    `json:"post_id" gorm:"uniqueIndex:idx_post_user;size:100;not null"` // 文章 ID
	UserID    *uint     `json:"user_id" gorm:"uniqueIndex:idx_post_user;index"`             // 登录用户 ID
	IPHash    string    `json:"-" gorm:"uniqueIndex:idx_post_ip;size:64"`                   // 访客 IP Hash
}

// LikeCount 文章点赞计数（缓存表）
type LikeCount struct {
	PostID string `json:"post_id" gorm:"primaryKey;size:100"`
	Count  int    `json:"count" gorm:"default:0"`
}

// PostMeta 文章元数据（可选，用于存储文章额外信息）
type PostMeta struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PostID     string    `json:"post_id" gorm:"uniqueIndex;size:100;not null"`
	ViewCount  int       `json:"view_count" gorm:"default:0"`
	LikeCount  int       `json:"like_count" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}