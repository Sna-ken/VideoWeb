package service

import (
	"context"
	"time"

	"github.com/Sna-ken/videoweb/biz/model/interact"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/Sna-ken/videoweb/pkg/e"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type InteractService struct {
	ctx context.Context
}

func NewInteractService(ctx context.Context) *InteractService {
	return &InteractService{ctx: ctx}
}

func (s *InteractService) LikeActionService(req *interact.LikeActionReq, userID string) error {
	hasVideoID := req.VideoID != ""
	hasCommentID := req.CommentID != ""

	if hasVideoID == hasCommentID {
		if hasVideoID {
			return e.New(consts.StatusBadRequest, "only choose one between video_id and comment_id", nil)
		} else {
			return e.New(consts.StatusBadRequest, "video_id and comment_id can't be both empty", nil)
		}
	}

	if userID == "" {
		return e.New(consts.StatusUnauthorized, "User ID not found", nil)
	}

	var _like model.Like
	var err error

	if hasVideoID {
		err = dao.FindLikeByUserIDAndVideoID(s.ctx, &_like, userID, req.VideoID)
	} else {
		err = dao.FindLikeByUserIDAndCommentID(s.ctx, &_like, userID, req.CommentID)
	}

	if err != nil {
		if req.ActionType == "1" {
			_like = model.Like{
				ID:        uuid.NewString(),
				UserID:    userID,
				VideoID:   req.VideoID,
				CommentID: req.CommentID,
				Create_at: time.Now(),
			}
			if err = dao.CreateLike(s.ctx, &_like); err != nil {
				return e.New(consts.StatusInternalServerError, "Create like failed", err)
			}
		}
		if req.ActionType == "0" {
			return e.New(consts.StatusBadRequest, "Like not existed", nil)
		}
	} else {
		if req.ActionType == "1" {
			return e.New(consts.StatusBadRequest, "Operation repeated", nil)
		}
		if req.ActionType == "0" {
			if req.VideoID != "" {
				if err := dao.CancelVideoLike(s.ctx, userID, req.VideoID); err != nil {
					return e.New(consts.StatusInternalServerError, "Cancel video like failed", err)
				}
			}
			if req.CommentID != "" {
				if err := dao.CancelCommentLike(s.ctx, userID, req.CommentID); err != nil {
					return e.New(consts.StatusInternalServerError, "Cancel comment like failed", err)
				}
			}
		}
	}

	if req.VideoID != "" {
		if err := dao.UpdateVideoLike(s.ctx, req.VideoID, req.ActionType); err != nil {
			return e.New(consts.StatusInternalServerError, "Update video like count failed", err)
		}
	}
	if req.CommentID != "" {
		if err := dao.UpdateCommentLike(s.ctx, req.CommentID, req.ActionType); err != nil {
			return e.New(consts.StatusInternalServerError, "Update comment like count failed", err)
		}
	}
	return nil
}

func (s *InteractService) LikeListService(req *interact.LikeListReq) (error, *interact.LikeListResp) {
	offset := int((req.PageNum - 1) * req.PageSize)

	var _videoList []model.Video
	var userID string
	var err error

	if req.UserID != "" && req.Username != "" {
		_userID, err := dao.FindUserIDByName(s.ctx, req.Username)
		if err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err), nil
		}
		if _userID != req.UserID {
			return e.New(consts.StatusBadRequest, "User ID and username are inconsistent", nil), nil
		}
		userID = req.UserID
	} else if req.UserID != "" {
		userID = req.UserID
	} else if req.Username != "" {
		userID, err = dao.FindUserIDByName(s.ctx, req.Username)
		if err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err), nil
		}
	}

	if err = dao.FindLikedVideoByUserID(s.ctx, userID, offset, int(req.PageSize), &_videoList); err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err), nil
	}
	items := make([]*interact.Item, 0)

	for _, v := range _videoList {
		item := &interact.Item{
			Username:     v.UserName,
			VideoID:      v.ID,
			UserID:       v.UserID,
			VideoURL:     v.Video_url,
			CoverURL:     v.Cover_url,
			Title:        v.Title,
			Description:  v.Description,
			VisitCount:   v.Visit_count,
			LikeCount:    v.Like_count,
			CommentCount: v.Comment_count,
			CreatedAt:    v.Create_at.Format(time.DateTime),
		}
		items = append(items, item)
	}

	return nil, &interact.LikeListResp{
		Base: &interact.Base{Code: consts.StatusOK, Msg: "Liked videos fetched successfully"},
		Data: &interact.Data{Item: items, Total: int32(len(items))}}
}

func (s *InteractService) CommentPublishService(req *interact.CommentPublishReq, userID string) error {
	if userID == "" {
		return e.New(consts.StatusUnauthorized, "User ID not found", nil)
	}

	hasVideoID := req.VideoID != ""
	hasParentID := req.ParentID != ""

	if hasVideoID == hasParentID {
		if hasVideoID {
			return e.New(consts.StatusBadRequest, "only choose one between video_id and parent_id", nil)
		} else {
			return e.New(consts.StatusBadRequest, "video_id and parent_id can't be both empty", nil)
		}
	}

	username, err := dao.FindUsernameByID(s.ctx, userID)
	if err != nil {
		return e.New(consts.StatusNotFound, "User not found", err)
	}

	var rootID string
	var videoID string

	if hasVideoID {
		rootID = ""
		videoID = req.VideoID
	} else {
		rootID, err = dao.FindRootIDByParentID(s.ctx, req.ParentID)
		if err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err)
		}
		videoID, err = dao.FindVideoIDByCommentID(s.ctx, req.ParentID)
		if err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err)
		}
	}

	_comment := model.Comment{
		UserName:  username,
		UserID:    userID,
		ID:        uuid.NewString(),
		VideoID:   videoID,
		Root_ID:   rootID,
		Parent_ID: req.ParentID,
		Content:   req.Content,
		Create_at: time.Now(),
		Update_at: time.Now(),
	}

	if err := dao.CreateComment(s.ctx, &_comment); err != nil {
		return e.New(consts.StatusInternalServerError, "Create comment failed", err)
	}

	if err := dao.AddVideoCommentCount(s.ctx, videoID, 1); err != nil {
		return e.New(consts.StatusInternalServerError, "Update video comment count failed", err)
	}

	if hasParentID {
		if err := dao.AddChildCommentCount(s.ctx, req.ParentID, 1); err != nil {
			return e.New(consts.StatusInternalServerError, "Update child comment count failed", err)
		}
		if rootID != req.ParentID {
			if err := dao.AddChildCommentCount(s.ctx, rootID, 1); err != nil {
				return e.New(consts.StatusInternalServerError, "Update child comment count failed", err)
			}
		}
	}

	return nil
}

func (s *InteractService) CommentListService(req *interact.CommentListReq) (error, *interact.CommentListResp) {
	offset := int((req.PageNum - 1) * req.PageSize)
	var _commentList []model.Comment

	hasVideoID := req.VideoID != ""
	hasRootCommentID := req.RootCommentID != ""

	if hasVideoID == hasRootCommentID {
		if hasVideoID {
			return e.New(consts.StatusBadRequest, "only choose one between video_id and root_comment_id", nil), nil
		} else {
			return e.New(consts.StatusBadRequest, "video_id and root_comment_id can't be both empty", nil), nil
		}
	}

	if hasVideoID {
		if err := dao.FindRootCommentByVideoID(s.ctx, req.VideoID, offset, int(req.PageSize), &_commentList); err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err), nil
		}
	} else {
		if err := dao.FindCommentByRootCommentID(s.ctx, req.RootCommentID, offset, int(req.PageSize), &_commentList); err != nil {
			return e.New(consts.StatusInternalServerError, "Database error", err), nil
		}
	}

	items := make([]*interact.CommentItem, 0, len(_commentList))

	for _, v := range _commentList {
		item := &interact.CommentItem{
			Username:   v.UserName,
			CommentID:  v.ID,
			VideoID:    v.VideoID,
			UserID:     v.UserID,
			RootID:     v.Root_ID,
			ParentID:   v.Parent_ID,
			Content:    v.Content,
			LikeCount:  v.Like_count,
			ChildCount: v.Child_count,
			CreatedAt:  v.Create_at.Format(time.DateTime),
		}
		items = append(items, item)
	}

	return nil, &interact.CommentListResp{
		Base: &interact.Base{Code: consts.StatusOK, Msg: "comments fetched successfully"},
		Data: &interact.CommentData{Item: items, Total: int32(len(items))}}
}

func (s *InteractService) CommentDeleteService(req *interact.CommentDeleteReq, userID string) error {
	if userID == "" {
		return e.New(consts.StatusUnauthorized, "User ID not found", nil)
	}
	videoID, err := dao.FindVideoIDByCommentID(s.ctx, req.CommentID)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err)
	}

	rootID, err := dao.FindRootIDByCommentID(s.ctx, req.CommentID)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err)
	}

	parentID, err := dao.FindParentIDByCommentID(s.ctx, req.CommentID)
	if err != nil {
		return e.New(consts.StatusInternalServerError, "Database error", err)
	}

	// 删除
	var deletecnt, childcnt int
	if err := dao.DeleteComment(s.ctx, req.CommentID, userID); err != nil {
		return e.New(consts.StatusNotFound, "Comment not found or no permission", err)
	}
	if err := dao.DeleteCommentLike(s.ctx, req.CommentID); err != nil {
		return e.New(consts.StatusInternalServerError, "Delete comment likes failed", err)
	}

	deletecnt = 1
	// 只有删除根评论时才删除所有子评论，如果是删除子评论，不删除其下的子评论
	if rootID == "" {
		if childcnt, err = dao.DeleteChildCommentByRootID(s.ctx, req.CommentID); err != nil {
			return e.New(consts.StatusNotFound, "Child comment not found or no permission", err)
		}
		if err := dao.DeleteCommentLikeByRootID(s.ctx, req.CommentID); err != nil {
			return e.New(consts.StatusInternalServerError, "Delete child comment likes failed", err)
		}
		deletecnt += childcnt
	}

	// 更新评论数
	if err := dao.ReduceVideoCommentCount(s.ctx, videoID, deletecnt); err != nil {
		return e.New(consts.StatusInternalServerError, "Update video comment count failed", err)
	}
	if parentID != "" {
		if err := dao.ReduceChildCommentCount(s.ctx, parentID, 1); err != nil {
			return e.New(consts.StatusInternalServerError, "Update child comment count failed", err)
		}
	}
	if rootID != "" && rootID != parentID {
		if err := dao.ReduceChildCommentCount(s.ctx, rootID, 1); err != nil {
			return e.New(consts.StatusInternalServerError, "Update child comment count failed", err)
		}
	}

	return nil
}
