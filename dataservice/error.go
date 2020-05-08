package dataservice

import (
	"errors"
	"fmt"
)

var (
	ErrPhoneExists  = errors.New("手机号码已存在")
	ErrInsertFailed = errors.New("保存用户信息失败")
)

// ErrIDNotFound 指定ID的用户不存在错误
type ErrIDNotFound struct {
	id int64
}

func NewErrIDNotFound(id int64) ErrIDNotFound {
	return ErrIDNotFound{id: id}
}
func (e ErrIDNotFound) Error() string {
	return fmt.Sprintf("未找到ID为%d的用户", e.id)
}

// ErrPhoneNotFound 指定手机的用户不存在错误
type ErrPhoneNotFound struct {
	phone string
}

func NewErrPhoneNotFound(phone string) ErrPhoneNotFound {
	return ErrPhoneNotFound{phone: phone}
}
func (e ErrPhoneNotFound) Error() string {
	return fmt.Sprintf("未找到手机为:%s 的用户", e.phone)
}

func IsErrPhoneNotFound(err error) bool {
	_, ok := err.(ErrPhoneNotFound)
	return ok
}
