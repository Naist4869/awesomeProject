package usermodel

import "time"

const (
	IDField           = "_id"
	PhoneField        = "phone"
	NickNameField     = "nickName"
	PidField          = "pid"
	StatusField       = "status"
	ActivateTimeField = "activateTime"
	OperateTimeField  = "operateTime"
	DeletedField      = "deleted"
)

type User struct {
	ID           int64       `bson:"_id" json:"id"`                                        // id
	Phone        string      `bson:"phone" json:"phone"`                                   // 手机号
	NickName     string      `bson:"nickname" json:"nickName"`                             // 昵称
	PID          int64       `bson:"pid" json:"pid"`                                       // 上级ID
	Status       int8        `bson:"status" json:"status"`                                 // 用户状态 1未激活 2 正常(解冻) 3冻结
	AddTime      time.Time   `bson:"addTime" json:"addTime"`                               // 注册添加时间
	ActivateTime time.Time   `bson:"activateTime,omitempty" json:"activateTime,omitempty"` // 做为普通用户激活时间(认证激活)
	OperateTime  time.Time   `bson:"operateTime,omitempty" json:"operateTime,omitempty"`   // 操作解冻(冻结)时间
	AgentTrait   AgentTrait  `bson:"agentTrait" json:"agentTrait"`                         // 代理特征
	WeChatTrait  WeChatTrait `bson:"weChatTrait" json:"weChatTrait"`                       // 微信用户特征
	Deleted      bool        `bson:"deleted" json:"-"`                                     // 是否被删除  false没删除 true已删除
	Meta         DbMeta      `bson:"meta" json:"meta"`                                     // 版本
}

// AgentTrait 代理特征
type AgentTrait struct {
	OpenTime    time.Time `bson:"openTime,omitempty" json:"openTime,omitempty"`       // 开通时间
	OperateTime time.Time `bson:"operateTime,omitempty" json:"operateTime,omitempty"` // 操作解冻(冻结)时间
	Team        Team      `bson:"team" json:"team"`                                   // 团队
	Status      int8      `bson:"status" json:"status"`                               // 代理状态 1未开通 2 正常(解冻) 3冻结
}

// WeChatTrait 微信用户特征
type WeChatTrait struct {
	OpenTime    time.Time `bson:"openTime,omitempty" json:"openTime,omitempty"`       // 开通时间
	OperateTime time.Time `bson:"operateTime,omitempty" json:"operateTime,omitempty"` // 操作解冻(冻结)时间
	OpenID      string    `bson:"openid,omitempty" json:"openid,omitempty"`           // 授权用户唯一标识
	NickName    string    `bson:"nickname,omitempty" json:"nickname,omitempty"`       // 普通用户昵称
	Sex         uint32    `bson:"sex,omitempty" json:"sex,omitempty"`                 // 普通用户性别，1为男性，2为女性
	Province    string    `bson:"province,omitempty" json:"province,omitempty"`       // 普通用户个人资料填写的省份
	City        string    `bson:"city,omitempty" json:"city,omitempty"`               // 普通用户个人资料填写的城市
	Country     string    `bson:"country,omitempty" json:"country,omitempty"`         // 国家，如中国为CN
	HeadImgURL  string    `bson:"headImgURL,omitempty" json:"headImgURL,omitempty"`   // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
	// Privilege  string `json:"privilege"`
	Privilege []string `bson:"privilege,omitempty" json:"privilege,omitempty"` // 用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
	UnionID   string   `bson:"unionID,omitempty" json:"unionID,omitempty"`     // 普通用户的标识，对当前开发者帐号唯一
	ErrCode   uint     `bson:"errCode,omitempty" json:"errCode,omitempty"`
	ErrMsg    string   `bson:"errMsg,omitempty" json:"errMsg,omitempty"`
	Status    int8     `bson:"status" json:"status"` // 微信用户状态 1未开通 2 正常(解冻) 3冻结
}

type Team struct {
	children []User
}

// DbMeta 数据库元数据信息
type DbMeta struct {
	Version int `bson:"version"` // 版本
}
