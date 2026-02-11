package dao

import (
	"context"

	"github.com/Sna-ken/videoweb/config"
)

func CreateUser(ctx context.Context, user *User) error {
	return config.MYSQLDB.WithContext(ctx).Create(user).Error
}
