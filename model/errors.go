package model

import "errors"

var (
	ErrUserPhoneEmpty    = errors.New("手机号不能为空")
	ErrUserNickNameEmpty = errors.New("昵称不能为空")
)
