package service

import (
	"context"

	"github.com/Sna-ken/videoweb/biz/model/social"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/Sna-ken/videoweb/pkg/e"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type SocialService struct {
	ctx context.Context
}

func NewSocialService(ctx context.Context) *SocialService {
	return &SocialService{ctx: ctx}
}

func (s *SocialService) FollowActionService(req *social.FollowActionReq, userID string) error {
	if userID == "" {
		return e.New(consts.StatusUnauthorized, "User ID not found", nil)
	}

	var toUserID string
	var toUserName string
	var err error
	var _object model.SocialObject
	if req.ToUserid != "" && req.ToUsername != "" {
		_toUserID, err := dao.FindUserIDByName(s.ctx, req.ToUsername)
		if err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err)
		}
		if _toUserID != req.ToUserid {
			return e.New(consts.StatusBadRequest, "User ID and username are inconsistent", nil)
		}
		toUserID = req.ToUserid
		toUserName = req.ToUsername
	} else if req.ToUserid != "" {
		toUserID = req.ToUserid
		toUserName, err = dao.FindUsernameByID(s.ctx, toUserID)
		if err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err)
		}
	} else if req.ToUsername != "" {
		toUserName = req.ToUsername
		toUserID, err = dao.FindUserIDByName(s.ctx, req.ToUsername)
		if err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err)
		}
	}

	if err != nil || toUserID == "" {
		return e.New(consts.StatusNotFound, "User not found", nil)
	}

	if userID == toUserID {
		return e.New(consts.StatusBadRequest, "Cannot follow yourself", nil)
	}

	avatarURL, err := dao.FindAvatarByUserID(s.ctx, toUserID)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err)
	}

	if err = dao.FindSocialObject(s.ctx, &_object, userID, toUserID); err != nil {
		if req.ActionType == "1" {
			_object = model.SocialObject{
				ID:         uuid.NewString(),
				UserID:     userID,
				ObjectID:   toUserID,
				Username:   toUserName,
				Avatar_url: avatarURL,
			}
			if err = dao.CreateSocialObject(s.ctx, &_object); err != nil {
				return e.New(consts.StatusInternalServerError, "Create social object failed", err)
			}
		}
		if req.ActionType == "0" {
			return e.New(consts.StatusBadRequest, "Operation repeated", nil)
		}
	} else {
		if req.ActionType == "1" {
			return e.New(consts.StatusBadRequest, "Operation repeated", nil)
		}
		if req.ActionType == "0" {
			if err = dao.RemoveSocialObject(s.ctx, userID, toUserID); err != nil {
				return e.New(consts.StatusInternalServerError, "Remove social object failed", err)
			}
		}
	}

	return nil
}

func (s *SocialService) FollowListService(req *social.FollowListReq, userID string) (error, *social.FollowListResp) {
	offset := int((req.PageNum - 1) * req.PageSize)

	var _objectList []model.SocialObject

	if err := dao.FindFollowByUserID(s.ctx, userID, offset, int(req.PageSize), &_objectList); err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err), nil
	}

	items := make([]*social.Item, 0, len(_objectList))

	for _, v := range _objectList {
		item := &social.Item{
			UserID:    v.ObjectID,
			Username:  v.Username,
			AvatarURL: v.Avatar_url,
		}
		items = append(items, item)
	}

	return nil, &social.FollowListResp{
		Base: &social.Base{Code: consts.StatusOK, Msg: "follow list fetched successfully"},
		Data: &social.Data{Item: items, Total: int32(len(items))}}
}

func (s *SocialService) FollowerListService(req *social.FollowerListReq, userID string) (error, *social.FollowerListResp) {
	offset := int((req.PageNum - 1) * req.PageSize)

	var _objectList []model.SocialObject

	if err := dao.FindFollowerByUserID(s.ctx, userID, offset, int(req.PageSize), &_objectList); err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err), nil
	}

	items := make([]*social.Item, 0, len(_objectList))

	for _, v := range _objectList {
		item := &social.Item{
			UserID:    v.UserID,
			Username:  v.Username,
			AvatarURL: v.Avatar_url,
		}
		items = append(items, item)
	}

	return nil, &social.FollowerListResp{
		Base: &social.Base{Code: consts.StatusOK, Msg: "follower list fetched successfully"},
		Data: &social.Data{Item: items, Total: int32(len(items))}}
}

func (s *SocialService) FriendListService(req *social.FriendListReq, userID string) (error, *social.FriendListResp) {
	offset := int((req.PageNum - 1) * req.PageSize)

	var _objectList []model.SocialObject
	if err := dao.FindFriendByUserID(s.ctx, userID, offset, int(req.PageSize), &_objectList); err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err), nil
	}
	items := make([]*social.Item, 0, len(_objectList))

	for _, v := range _objectList {
		item := &social.Item{
			UserID:    v.ObjectID,
			Username:  v.Username,
			AvatarURL: v.Avatar_url,
		}
		items = append(items, item)
	}

	return nil, &social.FriendListResp{
		Base: &social.Base{Code: consts.StatusOK, Msg: "friend list fetched successfully"},
		Data: &social.Data{Item: items, Total: int32(len(items))}}
}
