package dao

import (
	"context"
	"time"

	"github.com/Sna-ken/videoweb/config"
)

func CreateUser(ctx context.Context, user *User) error {
	return config.MYSQLDB.WithContext(ctx).Create(user).Error
}

func FindUserByName(ctx context.Context, user *User, username string) error {
	return config.MYSQLDB.Where("username = ?", username).First(user).Error
}

func SetRefreshToken(ctx context.Context, userID string, refreshtoken string, duration time.Duration) error {
	return config.REDISDB.Set(ctx, "user_rftoken:"+userID, refreshtoken, duration).Err()
}

func FindUserByID(ctx context.Context, user *User, userID string) error {
	return config.MYSQLDB.Where("id = ?", userID).First(user).Error
}
