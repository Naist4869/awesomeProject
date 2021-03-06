package dataservice

import (
	"context"

	"github.com/Naist4869/awesomeProject/model/officialmodel"
	"github.com/Naist4869/awesomeProject/model/wxmodel"

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

type IWorkWxDataService interface {
	Insert(u *wxmodel.UserInfo) (err error)
}

// 微信公众号rpc服务
type IOfficialWxRpcService interface {
	IFileSystemRpcService
	ITBKRpcService
}

// 文件系统rpc服务
type IFileSystemRpcService interface {
	MediaIDGet(ctx context.Context, req officialmodel.MediaIDReq, args ...interface{}) (resp officialmodel.MediaIDResp, err error)
}

// 淘宝客rpc服务
type ITBKRpcService interface {
	TitleConvertTBKey(ctx context.Context, req officialmodel.TitleConvertTBKeyReq, args ...interface{}) (resp officialmodel.TitleConvertTBKeyResp, err error)
	KeyConvertKey(ctx context.Context, req officialmodel.KeyConvertKeyReq, args ...interface{}) (resp officialmodel.KeyConvertKeyResp, err error)
}
