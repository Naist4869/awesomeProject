package usermodel

import (
	"time"
)

const (
	DBUserVersion = 1
	DBUserKey     = "users"
)

// 用户状态
const (
	_         = iota
	NotActive // 未激活
	Normal    // 正常(激活)
	Freeze    // 冻结
)

type RegisterArgument struct {
	Phone    string `json:"phone"`    // 手机号
	NickName string `json:"nickName"` // 昵称
	PID      int64  `json:"pid"`      // 上级ID
}

func (a RegisterArgument) Validate() error {
	if a.Phone == "" {
		return ErrUserPhoneEmpty
	}
	if a.NickName == "" {
		return ErrUserNickNameEmpty
	}
	return nil
}
func (a RegisterArgument) NewUserData(id int64) *User {
	return &User{
		ID:       id,
		Phone:    a.Phone,
		NickName: a.NickName,
		PID:      a.PID,
		Status:   Normal,
		AddTime:  time.Now(),
		Deleted:  false,
		Meta: DbMeta{
			Version: DBUserVersion,
		},
	}
}
