package e

import "errors"

var (
	ErrHasedPassword          = errors.New("Hashed Password failed")
	ErrDB                     = errors.New("database internal error")
	ErrUserHasExisted         = errors.New("user has alredy existed")
	ErrUserNotFound           = errors.New("user not found")
	ErrWrongPassword          = errors.New("wrong password")
	ErrGenerateToken          = errors.New("generate token failed")
	ErrUserIDNotFound         = errors.New("user ID not found")
	ErrFileRequired           = errors.New("file is empty")
	ErrFileSaveFailed         = errors.New("save file failed")
	ErrFileDeleteFailed       = errors.New("delete file failed")
	ErrFileOpenFailed         = errors.New("open file failed")
	ErrVideoOpenFailed        = errors.New("open video file failed")
	ErrVideoRequired          = errors.New("video file is empty")
	ErrUpdateLikeFailed       = errors.New("update like failed")
	ErrOperationRepeated      = errors.New("Operation repeated")
	ErrIDAndNameInconsistent  = errors.New("UserID and username are inconsistent")
	ErrNoPermissionOrNotFound = errors.New("Permission denied or record not found")
	ErrCanNotSelfFollow       = errors.New("Can't follow youself")
	ErrLikeNotexisted         = errors.New("like not existed")
)
