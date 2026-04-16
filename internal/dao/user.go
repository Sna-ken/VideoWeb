package dao

import (
	"context"
	"time"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/model"
)

func CreateUser(ctx context.Context, user *model.User) error {
	return config.MYSQLDB.WithContext(ctx).Create(user).Error
}

func FindUserByName(ctx context.Context, user *model.User, username string) error {
	return config.MYSQLDB.WithContext(ctx).Where("username = ?", username).First(user).Error
}

func SetRefreshToken(ctx context.Context, userID string, refreshtoken string) error {
	duration := time.Duration(config.JWT.RefreshTokenExpiry) * time.Second
	return config.REDISDB.Set(ctx, "user_rftoken:"+userID, refreshtoken, duration).Err()
}

func FindUserByID(ctx context.Context, user *model.User, userID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("id = ?", userID).First(user).Error
}

func UpdateUserAvatar(ctx context.Context, userID string, avatarURL string) error {
	return config.MYSQLDB.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"avatar_url": avatarURL,
		"update_at":  time.Now(),
	}).Error
}
