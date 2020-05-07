package userregistration

import (
	"errors"

	"github.com/Naist4869/awesomeProject/model/usermodel"

	"github.com/Naist4869/awesomeProject/dataservice"
)

type UseCase struct {
	UserDataInterface dataservice.IUserDataService
}

func (i *UseCase) RegisterUser(u *usermodel.User) (err error) {
	if u.PID > 0 {
		var count int64
		_, count, err = i.UserDataInterface.FindByPID(u.PID)
		if err != nil {
			return
		}
		if count == 0 {
			err = errors.New("上级ID不存在或者身份被冻结")
			return
		}
	}
	err = i.UserDataInterface.Insert(u)
	return
}

func (i *UseCase) UnregisterUser(userID int64) error {
	panic("implement me")
}

func (i *UseCase) ModifyUser(u *usermodel.User) error {
	panic("implement me")
}

func (i *UseCase) ModifyAndUnregister(u *usermodel.User) error {
	panic("implement me")
}
