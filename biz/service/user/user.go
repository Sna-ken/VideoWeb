package service

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	"github.com/Sna-ken/videoweb/biz/model/user"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/Sna-ken/videoweb/pkg/e"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/Sna-ken/videoweb/pkg/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const defaultAvatar = "http://127.0.0.1:8888/static/avatar/default/default-avatar.png"

type UserService struct {
	ctx context.Context
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{ctx: ctx}
}

func (s *UserService) RegisterService(req *user.RegisterReq) error {
	var tempUser model.User
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Hash password failed", err)
	}
	req.Password = hashedPassword

	if err := dao.FindUserByName(s.ctx, &tempUser, req.Username); err == nil {
		return e.New(consts.StatusConflict, "User has existed", err)
	} //先通过校检再初始化，FindUser和CreateUser不要重复使用同样的变量

	_user := model.User{
		ID:         uuid.New().String(),
		Username:   req.Username,
		Password:   req.Password,
		Avatar_url: defaultAvatar,
		Create_at:  time.Now(),
		Update_at:  time.Now(),
	}

	if err := dao.CreateUser(s.ctx, &_user); err != nil {
		return e.New(consts.StatusInternalServerError, "Create user failed", err)
	}

	return nil
}

func (s *UserService) LoginService(req *user.LoginReq) (error, *user.LoginResp) {
	var _user model.User

	if err := dao.FindUserByName(s.ctx, &_user, req.Username); err == gorm.ErrRecordNotFound {
		return e.New(consts.StatusNotFound, "User not found", err), nil
	}

	if !utils.CheckPasswordHash(req.Password, _user.Password) {
		return e.New(consts.StatusUnauthorized, "Wrong password", nil), nil
	}

	if req.MfaEnabled {
		if !_user.MFAEnabled {
			return e.New(consts.StatusUnauthorized, "MFA required", nil), nil
		}

		if !utils.ValidateMFA(req.Code, _user.MFASecret) {
			return e.New(consts.StatusUnauthorized, "Invalid MFA code", nil), nil
		}
	}

	accesstoken, err := jwt.GenerateAccessToken(_user.ID)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Generate access token failed", err), nil
	}

	refreshtoken, err := jwt.GenerateRefreshToken(_user.ID)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Generate refresh token failed", err), nil
	}

	if err := dao.SetRefreshToken(s.ctx, _user.ID, refreshtoken); err != nil {
		return e.New(consts.StatusInternalServerError, "Set refresh token failed", err), nil
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

func (s *UserService) UserInfoService(req *user.UserInfoReq, userID string) (error, *user.UserInfoResp) {
	if userID == "" {
		return e.New(consts.StatusNotFound, "User ID not found", nil), nil
	}

	var _user model.User
	if err := dao.FindUserByID(s.ctx, &_user, userID); err != nil {
		if err.Error() == "record not found" {
			return e.New(consts.StatusNotFound, "User not found", nil), nil
		}
		return e.New(consts.StatusInternalServerError, "Database error", err), nil
	}

	deletedAtStr := ""
	if _user.Delete_at != nil {
		deletedAtStr = _user.Delete_at.Format(time.DateTime)
	}

	return nil, &user.UserInfoResp{
		Base: &user.Base{
			Code: consts.StatusOK,
			Msg:  "user info fetched successfully",
		},
		Data: &user.Data{
			UserID:    _user.ID,
			Username:  _user.Username,
			AvatarURL: _user.Avatar_url,
			CreatedAt: _user.Create_at.Format(time.DateTime),
			UpdatedAt: _user.Update_at.Format(time.DateTime),
			DeletedAt: deletedAtStr,
		},
	}
}

func (s *UserService) UploadAvatarService(req *user.UploadAvatarReq, userID string, file *multipart.FileHeader) error {
	fileContent, err := file.Open()
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Open file failed", err)
	}

	defer fileContent.Close()
	if userID == "" {
		return e.New(consts.StatusNotFound, "User ID not found", nil)
	}

	avatarBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Read file failed", err)
	}

	if len(avatarBytes) == 0 {
		return e.New(consts.StatusBadRequest, "File is empty", nil)
	}
	var _user model.User
	if err := dao.FindUserByID(s.ctx, &_user, userID); err != nil {
		if err.Error() == "record not found" {
			return e.New(consts.StatusNotFound, "User not found", nil)
		}
		return e.New(consts.StatusInternalServerError, "Database error", err)
	}

	avatarURL, err := utils.StoreAvatar(avatarBytes, userID)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Save file failed", err)
	}

	if err := dao.UpdateUserAvatar(s.ctx, userID, avatarURL); err != nil {
		return e.New(consts.StatusInternalServerError, "Update user avatar failed", err)
	}

	return nil
}

func (s *UserService) GetMFAqrService(req *user.GetMFAqrReq, userID string) (error, *user.GetMFAqrResp) {
	if userID == "" {
		return e.New(consts.StatusNotFound, "User ID not found", nil), nil
	}
	username, err := dao.FindUsernameByID(s.ctx, userID)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err), nil
	}

	QRCode, err := utils.GenerateMFA(userID, username)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Generate MFA QR code failed", err), nil
	}

	return nil, &user.GetMFAqrResp{
		Base: &user.Base{
			Code: consts.StatusOK, Msg: "Get MFAqr successfully"},
		Qrcode: QRCode,
	}
}

func (s *UserService) BindMFAService(req *user.BindMFAReq, userID string) error {
	if userID == "" {
		return e.New(consts.StatusNotFound, "User ID not found", nil)
	}
	secret, err := dao.GetMFATemp(s.ctx, userID)
	if err != nil {
		return e.New(consts.StatusNotFound, "MFA temporary secret not found", nil)
	}

	if !utils.ValidateMFA(req.Code, secret) {
		return e.New(consts.StatusUnauthorized, "Invalid MFA code", nil)
	}

	if err := dao.StoreMFASecret(s.ctx, userID, secret); err != nil {
		return e.New(consts.StatusInternalServerError, "Store MFA secret failed", err)
	}

	if err := dao.DeleteMFATemp(s.ctx, userID); err != nil {
		return e.New(consts.StatusInternalServerError, "Delete MFA temporary secret failed", err)
	}

	if err := dao.EnableMFA(s.ctx, userID); err != nil {
		return e.New(consts.StatusInternalServerError, "Enable MFA failed", err)
	}

	return nil
}
