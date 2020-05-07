package dataservice

import (
	"context"

	"github.com/Naist4869/awesomeProject/model/usermodel"
)

type IUserDataService interface {
	// 删除操作
	Remove(userID int64) (err error)
	// 插入操作
	Insert(u *usermodel.User) (err error)
	// 查询操作
	FindByID(ctx context.Context, userID int64) (user *usermodel.User, err error)
	// 查询操作
	FindByIDs(userIDs []int64) (users []*usermodel.User, err error)
	// 查询操作
	FindByPID(userID int64) (user *usermodel.User, count int64, err error)
	// 查询操作
	FindByPhone(phone string) (user *usermodel.User, err error)
	// 更新操作
	Update(user *usermodel.User) (err error)
}
