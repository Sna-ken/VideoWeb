package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/Sna-ken/videoweb/biz/model/video"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/pkg/e"
	"github.com/Sna-ken/videoweb/pkg/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type VideoService struct {
	ctx context.Context
}

func NewVideoService(ctx context.Context) *VideoService {
	return &VideoService{ctx: ctx}
}

func (s *VideoService) PublishService(req *video.PublishReq, userID string, coverfile *multipart.FileHeader, videofile *multipart.FileHeader) error {
	if userID == "" {
		return e.ErrUserIDNotFound
	}

	username, err := dao.FindUsernameByID(s.ctx, userID)
	if err != nil {
		return e.ErrDB
	}

	coverContent, err := coverfile.Open()
	if err != nil {
		return e.ErrFileOpenFailed
	}
	defer coverContent.Close()
	videoContent, err := videofile.Open()
	if err != nil {
		return e.ErrVideoOpenFailed
	}
	defer videoContent.Close()

	coverBytes, err := io.ReadAll(coverContent)
	if err != nil {
		return e.ErrFileOpenFailed
	}
	videoBytes, err := io.ReadAll(videoContent)
	if err != nil {
		return e.ErrVideoOpenFailed
	}

	if len(coverBytes) == 0 {
		return e.ErrFileRequired
	}
	if len(videoBytes) == 0 {
		return e.ErrVideoRequired
	}

	videoURL, coverURL, err := utils.StoreVideo(videoBytes, coverBytes, userID)
	if err != nil {
		return e.ErrFileSaveFailed
	}

	_video := dao.Video{
		UserName:    username,
		ID:          uuid.New().String(),
		UserID:      userID,
		Video_url:   videoURL,
		Cover_url:   coverURL,
		Title:       req.Title,
		Description: req.Description,
		Create_at:   time.Now(),
		Update_at:   time.Now(),
	}

	if err := dao.CreateVideo(s.ctx, &_video); err != nil {
		return e.ErrDB
	}
	return nil
}

func (s *VideoService) ListService(req *video.ListReq, userID string) (error, *video.ListResp) {
	offset := int((req.PageNum - 1) * req.PageSize)

	var _videoList []dao.Video
	if userID == "" {
		return e.ErrUserIDNotFound, nil
	}

	if err := dao.FindVideoByUserID(s.ctx, &_videoList, userID, offset, int(req.PageSize)); err != nil {
		return e.ErrDB, nil
	}

	items := make([]*video.Item, 0)

	for _, v := range _videoList {
		item := &video.Item{
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
	return nil, &video.ListResp{
		Base: &video.Base{Code: consts.StatusOK, Msg: "video list fetched successfully"},
		Data: &video.Data{Item: items, Total: int32(len(items))}}
}

func (s *VideoService) PopularService(req *video.PopularReq) (error, *video.PopularResp) {
	offset := int((req.PageNum - 1) * req.PageSize)
	var _videoList []dao.Video

	if err := dao.OrderPopular(s.ctx, &_videoList, offset, int(req.PageSize)); err != nil {
		return e.ErrDB, nil
	}
	items := make([]*video.Item, 0)

	for _, v := range _videoList {
		item := &video.Item{
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
		}
		items = append(items, item)
	}
	return nil, &video.PopularResp{
		Base: &video.Base{Code: consts.StatusOK, Msg: "popular videos fetched successfully"},
		Data: &video.Data{Item: items, Total: int32(len(items))}}

}

func (s *VideoService) SearchService(req *video.SearchReq) (error, *video.SearchResp) {
	offset := int((req.PageNum - 1) * req.PageSize)
	var _videoList []dao.Video
	if req.Keyword != "" {
		err := dao.SearchByKeyword(s.ctx, &_videoList, offset, int(req.PageSize), req.Keyword)
		if err != nil {
			return e.ErrDB, nil
		}
	}

	if req.Username != "" {
		err := dao.SearchByName(s.ctx, &_videoList, offset, int(req.PageSize), req.Username)
		if err != nil {
			return e.ErrDB, nil
		}
	}
	items := make([]*video.Item, 0)

	for _, v := range _videoList {
		item := &video.Item{
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
		}
		items = append(items, item)
	}
	return nil, &video.SearchResp{
		Base: &video.Base{Code: consts.StatusOK, Msg: fmt.Sprintf("search results for keyword '%s' or username '%s' fetched successfully", req.Keyword, req.Username)},
		Data: &video.Data{Item: items, Total: int32(len(items))}}
}
