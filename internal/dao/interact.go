package dao

import (
	"context"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/model"
	"gorm.io/gorm"
)

func FindLikeByUserIDAndVideoID(ctx context.Context, like *model.Like, userID string, videoID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).First(like).Error
}

func FindLikeByUserIDAndCommentID(ctx context.Context, like *model.Like, userID string, commentID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).First(like).Error
}

func CreateLike(ctx context.Context, like *model.Like) error {
	return config.MYSQLDB.WithContext(ctx).Create(like).Error
}

func CancelVideoLike(ctx context.Context, userID string, videoID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&model.Like{}).Error
}

func CancelCommentLike(ctx context.Context, userID string, commentID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&model.Like{}).Error
}

func UpdateVideoLike(ctx context.Context, videoID string, statu string) error {
	if statu == "1" {
		return config.MYSQLDB.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoID).Update("like_count", gorm.Expr("like_count + ?", 1)).Error
	}
	if statu == "0" {
		return config.MYSQLDB.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoID).Update("like_count", gorm.Expr("like_count - ?", 1)).Error
	}
	return nil
}

func UpdateCommentLike(ctx context.Context, commetID string, statu string) error {
	if statu == "1" {
		return config.MYSQLDB.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", commetID).Update("like_count", gorm.Expr("like_count + ?", 1)).Error
	}
	if statu == "0" {
		return config.MYSQLDB.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", commetID).Update("like_count", gorm.Expr("like_count - ?", 1)).Error
	}
	return nil
}

func FindLikedVideoByUserID(ctx context.Context, userID string, offset int, pagesize int, video *[]model.Video) error {
	rsl := config.MYSQLDB.WithContext(ctx).Model(&model.Video{}).Select("videos.*").Joins("INNER JOIN likes ON likes.video_id = videos.id").
		Where("likes.user_id = ? AND videos.delete_at IS NULL", userID)

	if err := rsl.Offset(offset).Limit(pagesize).Find(video).Error; err != nil {
		return err
	}

	return nil
}

func FindUserIDByName(ctx context.Context, username string) (userID string, err error) {
	err = config.MYSQLDB.WithContext(ctx).Table("users").
		Select("id").Where("username = ?", username).Scan(&userID).Error

	return userID, err
}

func CreateComment(ctx context.Context, comment *model.Comment) error {
	return config.MYSQLDB.WithContext(ctx).Create(comment).Error
}

func AddVideoCommentCount(ctx context.Context, videoID string, deletecnt int) error {
	return config.MYSQLDB.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count + ?", deletecnt)).Error
}
func AddChildCommentCount(ctx context.Context, parentID string, deletecnt int) error {
	return config.MYSQLDB.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", parentID).Update("child_count", gorm.Expr("child_count + ?", deletecnt)).Error
}

func ReduceVideoCommentCount(ctx context.Context, videoID string, deletecnt int) error {
	return config.MYSQLDB.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count - ?", deletecnt)).Error
}

func ReduceChildCommentCount(ctx context.Context, parentID string, deletecnt int) error {
	return config.MYSQLDB.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", parentID).Update("child_count", gorm.Expr("child_count - ?", deletecnt)).Error
}

func FindRootCommentByVideoID(ctx context.Context, videoID string, offset int, pagesize int, comment *[]model.Comment) error {
	return config.MYSQLDB.WithContext(ctx).Where("video_id = ? AND root_id = ''", videoID).Offset(offset).Limit(pagesize).Find(comment).Error
}
func FindCommentByRootCommentID(ctx context.Context, rootCommentID string, offset int, pagesize int, comment *[]model.Comment) error {
	return config.MYSQLDB.WithContext(ctx).Where("root_id = ?", rootCommentID).Offset(offset).Limit(pagesize).Find(comment).Error
}

func DeleteComment(ctx context.Context, commentID string, userID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND id = ?", userID, commentID).Delete(&model.Comment{}).Error
}

func FindVideoIDByCommentID(ctx context.Context, commentID string) (videoID string, err error) {
	err = config.MYSQLDB.WithContext(ctx).Table("comments").
		Select("video_id").Where("id = ?", commentID).Scan(&videoID).Error
	return videoID, err
}

func FindParentIDByCommentID(ctx context.Context, commentID string) (parentID string, err error) {
	err = config.MYSQLDB.WithContext(ctx).Table("comments").
		Select("parent_id").Where("id = ?", commentID).Scan(&parentID).Error
	return parentID, err
}

func FindRootIDByCommentID(ctx context.Context, commentID string) (rootID string, err error) {
	err = config.MYSQLDB.WithContext(ctx).Table("comments").
		Select("root_id").Where("id = ?", commentID).Scan(&rootID).Error
	return rootID, err
}

func DeleteCommentLike(ctx context.Context, commentID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("comment_id = ?", commentID).Delete(&model.Like{}).Error
}

func FindRootIDByParentID(ctx context.Context, parentID string) (rootID string, err error) {
	err = config.MYSQLDB.WithContext(ctx).Table("comments").
		Select("root_id").Where("id = ?", parentID).Scan(&rootID).Error
	if rootID == "" {
		rootID = parentID
	}
	return rootID, err
}

func DeleteChildCommentByRootID(ctx context.Context, rootID string) (int, error) {
	rsl := config.MYSQLDB.WithContext(ctx).Where("root_id = ?", rootID).Delete(&model.Comment{})
	return int(rsl.RowsAffected), rsl.Error
}

func DeleteCommentLikeByParentID(ctx context.Context, parentID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("comment_id IN (SELECT id FROM comments WHERE parent_id = ?)", parentID).Delete(&model.Like{}).Error
}

func DeleteCommentLikeByRootID(ctx context.Context, rootID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("comment_id IN (SELECT id FROM comments WHERE root_id = ?)", rootID).Delete(&model.Like{}).Error
}
