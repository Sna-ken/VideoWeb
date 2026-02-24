package service

import (
	"context"
	"time"

	"github.com/Sna-ken/videoweb/biz/model/user"
	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/pkg/e"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/Sna-ken/videoweb/pkg/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const defaultAvatar = "http://127.0.0.1:8888/static/default-avatar.png"

type UserService struct {
	ctx context.Context
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{ctx: ctx}
}

func (s *UserService) RegisterService(req *user.RegisterReq) error {
	var tempUser dao.User
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return e.ErrHasedPassword
	}
	req.Password = hashedPassword

	if err := dao.FindUser(s.ctx, &tempUser, req.Username); err == nil {
		return e.ErrUserHasExisted
	} //先通过校检再初始化，FindUser和CreateUser不要重复使用同样的变量

	_user := dao.User{
		ID:         uuid.New().String(),
		Username:   req.Username,
		Password:   req.Password,
		Avatar_url: defaultAvatar,
		Create_at:  time.Now(),
		Update_at:  time.Now(),
	}

	if err := dao.CreateUser(s.ctx, &_user); err != nil {
		return e.ErrDB
	}

	return nil
}

func (s *UserService) LoginService(req *user.LoginReq) (error, *user.LoginResp) {
	var _user dao.User

	if err := dao.FindUser(s.ctx, &_user, req.Username); err == gorm.ErrRecordNotFound {
		return e.ErrUserNotFound, nil
	}

	if !utils.CheckPasswordHash(req.Password, _user.Password) {
		return e.ErrWrongPassword, nil
	}

	accesstoken, err := jwt.GenerateAccessToken(_user.ID)
	if err != nil {
		return e.ErrGenerateToken, nil
	}

	refreshtoken, err := jwt.GenerateRefreshToken(_user.ID)
	if err != nil {
		return e.ErrGenerateToken, nil
	}

	duration := time.Duration(config.JWTConfig.RefreshTokenExpiry) * time.Second
	if err := dao.SetRefreshToken(s.ctx, _user.ID, refreshtoken, duration); err != nil {
		return e.ErrDB, nil
	}

	deletedAtStr := ""
	if _user.Delete_at != nil {
		deletedAtStr = _user.Delete_at.Format(time.DateTime)
	}

	return nil, &user.LoginResp{
		Base: &user.Base{Code: consts.StatusOK, Msg: "User login successfully"},
		Data: &user.Data{
			UserID:    _user.ID,
			Username:  _user.Username,
			AvatarURL: _user.Avatar_url,
			CreatedAt: _user.Create_at.Format(time.DateTime),
			UpdatedAt: _user.Update_at.Format(time.DateTime),
			DeletedAt: deletedAtStr, //没删除用户指针为nil,直接调用Format会Panic
		},
		AccessToken:  accesstoken,
		RefreshToken: refreshtoken,
	}
}
