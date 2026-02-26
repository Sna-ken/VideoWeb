package e

import "errors"

var (
	ErrHasedPassword    = errors.New("Hashed Password failed")
	ErrDB               = errors.New("database internal error")
	ErrUserHasExisted   = errors.New("user has alredy existed")
	ErrUserNotFound     = errors.New("user not found")
	ErrWrongPassword    = errors.New("wrong password")
	ErrGenerateToken    = errors.New("generate token failed")
	ErrUserIDNotFound   = errors.New("user ID not found")
	ErrFileRequired     = errors.New("file is empty")
	ErrFileSaveFailed   = errors.New("save file failed")
	ErrFileDeleteFailed = errors.New("delete file failed")
)
