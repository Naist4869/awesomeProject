package usecase

import (
	"github.com/Naist4869/awesomeProject/api"
	"github.com/Naist4869/awesomeProject/model"
	"github.com/Naist4869/awesomeProject/model/apimodel"
	"github.com/Naist4869/awesomeProject/model/usermodel"
	"github.com/Naist4869/awesomeProject/model/wxmodel"
)

type IUserRegistration interface {
	// 注册用户
	RegisterUser(u usermodel.RegisterArgument) (err error)
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

type IWorkWx interface {
	//读取成员
	GetUser(userName string) (*wxmodel.UserInfo, error)
	//获取部门成员详情
	UserListByDept(deptID int64, fetchChild bool) ([]*wxmodel.UserInfo, error)
	//获取全量组织架构
	ListAllDepts() ([]*wxmodel.DeptInfo, error)
	//获取指定部门及其下的子部门
	ListDepts(id int64) ([]*wxmodel.DeptInfo, error)
	// 获取群聊会话
	GetAppchat(chatid string) (*wxmodel.ChatInfo, error)
	// SendTextMessage 发送文本消息
	SendTextMessage(recipient *wxmodel.Recipient, content string, isSafe bool) error
	// 发送图片消息
	SendImageMessage(recipient *wxmodel.Recipient, mediaID string, isSafe bool) error
	// 发送语音消息
	SendVoiceMessage(recipient *wxmodel.Recipient, mediaID string, isSafe bool) error
	// 发送视频消息
	SendVideoMessage(recipient *wxmodel.Recipient, mediaID string, description string, title string, isSafe bool) error
	// 发送文件消息
	SendFileMessage(recipient *wxmodel.Recipient, mediaID string, isSafe bool) error
	// 发送文本卡片消息
	SendTextCardMessage(recipient *wxmodel.Recipient, title string, description string, url string, buttonText string, isSafe bool) error
	//发送图文消息
	SendNewsMessage(recipient *wxmodel.Recipient, title string, description string, url string, picURL string, isSafe bool) error
	//发送 mpnews 类型的图文消息
	SendMPNewsMessage(recipient *wxmodel.Recipient, title string, thumbMediaID string, author string, sourceContentURL string, content string, digest string, isSafe bool) error
	//发送 Markdown 消息
	SendMarkdownMessage(recipient *wxmodel.Recipient, content string, isSafe bool) error
}
type IOfficialWx interface {
	ReplyMessage(xmlMsg []byte) (reply string, err error)
}
type IAPI interface {
	Handle(arg *apimodel.Argument, api api.Api) (interface{}, int, error)
}
