package dao

import (
	"context"

	"github.com/Sna-ken/videoweb/config"
)

func FindAvatarByUserID(ctx context.Context, userID string) (avatarURL string, err error) {
	err = config.MYSQLDB.WithContext(ctx).Table("users").
		Select("avatar_url").Where("id = ?", userID).Scan(&avatarURL).Error
	return avatarURL, err
}

func FindSocialObject(ctx context.Context, object *SocialObject, userID string, ToUserID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND object_id = ?", userID, ToUserID).First(object).Error
}

func CreateSocialObject(ctx context.Context, object *SocialObject) error {
	return config.MYSQLDB.WithContext(ctx).Create(object).Error
}

func RemoveSocialObject(ctx context.Context, userID string, ToUserID string) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ? AND object_id = ?", userID, ToUserID).Delete(&SocialObject{}).Error
}

func FindFollowByUserID(ctx context.Context, userID string, offset int, pagesize int, object *[]SocialObject) error {
	return config.MYSQLDB.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(pagesize).Find(object).Error
}

func FindFollowerByUserID(ctx context.Context, userID string, offset int, pagesize int, object *[]SocialObject) error {
	return config.MYSQLDB.WithContext(ctx).Table("social_object").Select("users.id as user_id, users.username,users.avatar_url").
		Joins("JOIN users ON social_object.user_id = users.id").
		Where("social_object.object_id = ?", userID).Offset(offset).Limit(pagesize).Find(object).Error
}

func FindFriendByUserID(ctx context.Context, userID string, offset int, pagesize int, object *[]SocialObject) error {
	return config.MYSQLDB.WithContext(ctx).Table("social_objects as o1").Select("o1.*").
		Joins("JOIN social_objects as o2 ON o1.user_id = o2.object_id AND o1.object_id = o2.user_id").
		Where("o1.user_id = ?", userID).Offset(offset).Limit(pagesize).Find(object).Error
}
