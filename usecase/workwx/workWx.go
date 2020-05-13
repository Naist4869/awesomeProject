package workwx

import (
	"github.com/Naist4869/awesomeProject/dataservice"
	"github.com/Naist4869/awesomeProject/model/wxmodel"
)

type UseCase struct {
	WorkWxDataService dataservice.IWorkWxDataService
}

func (u *UseCase) GetUser(userName string) (*wxmodel.UserInfo, error) {
	panic("implement me")
}

func (u *UseCase) UserListByDept(deptID int64, fetchChild bool) ([]*wxmodel.UserInfo, error) {
	panic("implement me")
}

func (u *UseCase) ListAllDepts() ([]*wxmodel.DeptInfo, error) {
	panic("implement me")
}

func (u *UseCase) ListDepts(id int64) ([]*wxmodel.DeptInfo, error) {
	panic("implement me")
}

func (u *UseCase) GetAppchat(chatid string) (*wxmodel.ChatInfo, error) {
	panic("implement me")
}

func (u *UseCase) SendTextMessage(recipient *wxmodel.Recipient, content string, isSafe bool) error {
	panic("implement me")
}

func (u *UseCase) SendImageMessage(recipient *wxmodel.Recipient, mediaID string, isSafe bool) error {
	panic("implement me")
}

func (u *UseCase) SendVoiceMessage(recipient *wxmodel.Recipient, mediaID string, isSafe bool) error {
	panic("implement me")
}

func (u *UseCase) SendVideoMessage(recipient *wxmodel.Recipient, mediaID string, description string, title string, isSafe bool) error {
	panic("implement me")
}

func (u *UseCase) SendFileMessage(recipient *wxmodel.Recipient, mediaID string, isSafe bool) error {
	panic("implement me")
}

func (u *UseCase) SendTextCardMessage(recipient *wxmodel.Recipient, title string, description string, url string, buttonText string, isSafe bool) error {
	panic("implement me")
}

func (u *UseCase) SendNewsMessage(recipient *wxmodel.Recipient, title string, description string, url string, picURL string, isSafe bool) error {
	panic("implement me")
}

func (u *UseCase) SendMPNewsMessage(recipient *wxmodel.Recipient, title string, thumbMediaID string, author string, sourceContentURL string, content string, digest string, isSafe bool) error {
	panic("implement me")
}

func (u *UseCase) SendMarkdownMessage(recipient *wxmodel.Recipient, content string, isSafe bool) error {
	panic("implement me")
}
