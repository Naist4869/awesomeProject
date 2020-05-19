package usecase

import "errors"

var (
	ErrParseMessage             = errors.New("解析消息失败")
	ErrMessageTypeAssertionFail = errors.New("消息类型断言失败")
)
