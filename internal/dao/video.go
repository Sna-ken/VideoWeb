package dao

import (
	"context"
	"time"

	"github.com/Sna-ken/videoweb/config"
)

func CreateVideo(ctx context.Context, video *Video) error {
	return config.MYSQLDB.WithContext(ctx).Create(video).Error
}

func UploadVideoCover(ctx context.Context, userID string, videoID string, coverURL string) error {
	return config.MYSQLDB.WithContext(ctx).Model(&User{}).Where("user_id = ? AND video_id = ?", userID, videoID).Updates(map[string]interface{}{
		"cover_url": coverURL,
		"update_at": time.Now(),
	}).Error
}

func FindVideoByUserID(ctx context.Context, video *[]Video, userID string, offset int, pagesize int) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(pagesize).Find(video).Error
}

func OrderPopular(ctx context.Context, video *[]Video, offset int, pagesize int) error {
	return config.MYSQLDB.WithContext(ctx).Order("visit_count * 2 + like_count * 3 + comment_count * 5 desc").
		Offset(offset).Limit(pagesize).Find(video).Error
}

func SearchByKeyword(ctx context.Context, video *[]Video, offset int, pagesize int, keyword string) error {
	return config.MYSQLDB.WithContext(ctx).Where("title LIKE ? OR description LIKE ? ", "%"+keyword+"%", "%"+keyword+"%").
		Offset(offset).Limit(pagesize).Find(video).Error
}

func SearchByName(ctx context.Context, video *[]Video, offset int, pagesize int, name string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_name LIKE ?", "%"+name+"%").
		Offset(offset).Limit(pagesize).Find(video).Error
}

func FindUsernameByID(ctx context.Context, userID string) (username string, err error) {
	err = config.MYSQLDB.WithContext(ctx).Table("users").
		Select("username").Where("id = ?", userID).Scan(&username).Error

	return username, err
}
