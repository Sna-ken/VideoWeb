package service

import (
	"context"
	"time"

	"github.com/Sna-ken/videoweb/biz/model/user"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/pkg/utils"
	"github.com/google/uuid"
)

const defaultAvatar = "http://127.0.0.1:8888/static/default-avatar.png"

type UserService struct {
	ctx context.Context
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{ctx: ctx}
}

func (s *UserService) RegisterService(req *user.RegisterReq) error {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashedPassword

	_user := dao.User{
		ID:         uuid.New().String(),
		Username:   req.Username,
		Password:   req.Password,
		Avatar_url: defaultAvatar,
		Create_at:  time.Now(),
		Update_at:  time.Now(),
	}
	return dao.CreateUser(s.ctx, &_user)
}
