package model

import "errors"

var (
	ErrUserPhoneEmpty    = errors.New("手机号不能为空")
	ErrUserNickNameEmpty = errors.New("昵称不能为空")
	ErrIDInvalid         = errors.New("id非法")
	ErrTimeOut           = errors.New("超时")
	ErrAPIKey            = errors.New("api key非法")
	ErrMarshal           = errors.New("参数序列化错误")
)
