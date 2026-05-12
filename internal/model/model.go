package model

import (
	"time"
)

type User struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	Username   string     `json:"username" gorm:"unique;not null"`
	Password   string     `json:"password" gorm:"not null"`
	Avatar_url string     `json:"avatar_url" gorm:"not null"`
	MFASecret  string     `json:"mfa_secret" gorm:"not null;default:''"`
	MFAEnabled bool       `json:"mfa_enabled" gorm:"not null;default:false"`
	Create_at  time.Time  `json:"create_at"`
	Update_at  time.Time  `json:"update_at"`
	Delete_at  *time.Time `json:"delete_at"`
}

type Video struct {
	UserName      string     `json:"username" gorm:"not null"`
	ID            string     `json:"id" gorm:"primaryKey;comment:video_id"`
	UserID        string     `json:"user_id" gorm:"not null;comment:user_id"`
	Video_url     string     `json:"video_url" gorm:"not null"`
	Cover_url     string     `json:"cover_url" gorm:"not null"`
	Title         string     `json:"title" gorm:"not null"`
	Description   string     `json:"description" gorm:"not null"`
	Visit_count   int32      `json:"visit_count" gorm:"not null;default:0"`
	Like_count    int32      `json:"like_count" gorm:"not null;default:0;check:like_count >= 0"`
	Comment_count int32      `json:"comment_count" gorm:"not null;default:0;check:like_count >= 0"`
	Create_at     time.Time  `json:"create_at"`
	Update_at     time.Time  `json:"update_at"`
	Delete_at     *time.Time `json:"delete_at"`
}

type Like struct {
	ID        string    `json:"id" gorm:"primaryKey;comment:like_id"`
	UserID    string    `json:"user_id" gorm:"not null;comment:user_id"`
	CommentID string    `json:"comment_id" gorm:"comment:comment_id"`
	VideoID   string    `json:"video_id" gorm:"comment:video_id"`
	Create_at time.Time `json:"create_at"`
}

type Comment struct {
	UserName    string     `json:"username" gorm:"not null"`
	ID          string     `json:"id" gorm:"primaryKey;comment:comment_id"`
	UserID      string     `json:"user_id" gorm:"not null;comment:user_id"`
	VideoID     string     `json:"video_id" gorm:"comment:video_id"`
	Root_ID     string     `json:"root_id" gorm:"comment:root_id"`
	Parent_ID   string     `json:"parent_id" gorm:"comment:parent_id"`
	Like_count  int32      `json:"like_count" gorm:"not null;default:0;check:like_count >= 0"`
	Child_count int32      `json:"child_count" gorm:"not null;default:0;check:child_count >= 0"`
	Content     string     `json:"content" gorm:"not null"`
	Create_at   time.Time  `json:"create_at"`
	Update_at   time.Time  `json:"update_at"`
	Delete_at   *time.Time `json:"delete_at"`
}

type SocialObject struct {
	ID         string `json:"id" gorm:"primaryKey"`
	UserID     string `json:"user_id" gorm:"not null;comment:user_id"`
	ObjectID   string `json:"object_id" gorm:"primaryKey;comment:object_id"`
	Username   string `json:"username" gorm:"not null"`
	Avatar_url string `json:"avatar_url" gorm:"not null"`
}

type Timestamp struct {
	Create_at time.Time  `json:"create_at"`
	Update_at time.Time  `json:"update_at"`
	Delete_at *time.Time `json:"delete_at"`
}

type Response struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type Message struct {
	ID        string    `json:"id" gorm:"primaryKey;comment:message_id"`
	UserID    string    `json:"user_id" gorm:"not null;comment:user_id"`
	Username  string    `json:"username" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	Create_at time.Time `json:"create_at"`
}
