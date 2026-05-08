package dao

import (
	"context"
	"fmt"
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

func StoreMFASecret(ctx context.Context, userID string, secret string) error {
	return config.MYSQLDB.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("mfa_secret", secret).Error
}

func StoreMFATemp(ctx context.Context, userID string, secret string) error {
	return config.REDISDB.Set(ctx, fmt.Sprintf("mfa:%s", userID), secret, time.Minute*5).Err()
}

func GetMFATemp(ctx context.Context, userID string) (string, error) {
	secret, err := config.REDISDB.Get(ctx, fmt.Sprintf("mfa:%s", userID)).Result()
	return secret, err
}

func DeleteMFATemp(ctx context.Context, userID string) error {
	return config.REDISDB.Del(ctx, fmt.Sprintf("mfa:%s", userID)).Err()
}

func EnableMFA(ctx context.Context, userID string) error {
	return config.MYSQLDB.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("mfa_enabled", true).Error
}
