package usecase

import (
	"github.com/Naist4869/awesomeProject/model"
	"github.com/Naist4869/awesomeProject/model/usermodel"
)

type IUserRegistration interface {
	// 注册用户
	RegisterUser(u *usermodel.User) (err error)
	// 注销用户
	UnregisterUser(userID int64) error
	// 修改用户
	ModifyUser(u *usermodel.User) error
	// 修改后注销
	ModifyAndUnregister(u *usermodel.User) error
}

type IUserQuery interface {
	// 查找用户
	QueryUser(t model.Table) (u []*usermodel.User, count int, err error)
}

type IUserTeam interface {
	// 获取团队树
	GetTeamTree(userID int64) (interface{}, error)
}
