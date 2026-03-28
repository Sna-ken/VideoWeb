package dao

import (
	"context"

	"github.com/Sna-ken/videoweb/config"
	"gorm.io/gorm"
)

func FindLikeByUserIDAndVideoID(ctx context.Context, like *Like, userID string, videoID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).First(like).Error
}

func FindLikeByUserIDAndCommentID(ctx context.Context, like *Like, userID string, commentID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).First(like).Error
}

func CreateLike(ctx context.Context, like *Like) error {
	return config.MYSQLDB.WithContext(ctx).Create(like).Error
}

func CancelVideoLike(ctx context.Context, userID string, videoID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&Like{}).Error
}

func CancelCommentLike(ctx context.Context, userID string, commentID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&Like{}).Error
}

func UpdateVideoLike(ctx context.Context, videoID string, statu string) error {
	if statu == "1" {
		return config.MYSQLDB.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).Update("like_count", gorm.Expr("like_count + ?", 1)).Error
	}
	if statu == "0" {
		return config.MYSQLDB.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).Update("like_count", gorm.Expr("like_count - ?", 1)).Error
	}
	return nil
}

func UpdateCommentLike(ctx context.Context, CommetID string, statu string) error {
	if statu == "1" {
		return config.MYSQLDB.WithContext(ctx).Model(&Comment{}).Where("id = ?", CommetID).Update("like_count", gorm.Expr("like_count + ?", 1)).Error
	}
	if statu == "0" {
		return config.MYSQLDB.WithContext(ctx).Model(&Comment{}).Where("id = ?", CommetID).Update("like_count", gorm.Expr("like_count - ?", 1)).Error
	}
	return nil
}

func FindLikedVideoByUserID(ctx context.Context, userID string, offset int, pagesize int, video *[]Video) error {
	rsl := config.MYSQLDB.WithContext(ctx).Model(&Video{}).Select("videos.*").Joins("INNER JOIN likes ON likes.video_id = videos.id").
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

func CreateComment(ctx context.Context, comment *Comment) error {
	return config.MYSQLDB.WithContext(ctx).Create(comment).Error
}

func AddCommentCount(ctx context.Context, videoID string) error {
	return config.MYSQLDB.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}

func ReduceCommentCount(ctx context.Context, videoID string) error {
	return config.MYSQLDB.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
}

func FindCommentByVideoID(ctx context.Context, videoID string, offset int, pagesize int, comment *[]Comment) error {
	return config.MYSQLDB.WithContext(ctx).Where("video_id = ?", videoID).Offset(offset).Limit(pagesize).Find(comment).Error
}

func DeleteComment(ctx context.Context, videoID string, userID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND id = ?", userID, videoID).Delete(&Comment{}).Error
}

func FindVideoIDByCommentID(ctx context.Context, commentID string) (videoID string, err error) {
	err = config.MYSQLDB.WithContext(ctx).Table("comments").
		Select("video_id").Where("id = ?", commentID).Scan(&videoID).Error
	return videoID, err
}
