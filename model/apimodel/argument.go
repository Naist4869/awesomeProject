package apimodel

import (
	"encoding/json"
	"time"

	"github.com/Naist4869/awesomeProject/model"
)

const (
	// 错误代码
	OK            = 0 // 正确
	IDInvalidCode = 1 // ID非法
	TimeOutCode   = 2 // 验证超时
	ArgumentCode  = 3 // 参数错误
	ApiKeyCode    = 4 // 方法key错误
	AuthCode      = 5 // 验证错误
	MarshalCode   = 6 // 参数序列化错误

	readmeDocPath = "data/api.html"
)

//Argument 请求参数
type Argument struct {
	ID     int64           `json:"id"`     //商户ID
	Time   int64           `json:"time"`   //调用发起时间,unix epoch 精确到秒
	Key    string          `json:"key"`    //加密之后的key
	Data   json.RawMessage `json:"data"`   //调用参数
	APIkey string          `json:"apiKey"` //调用API的key
}

func (a Argument) Validate() (int, error) {
	if a.ID < 0 {
		return IDInvalidCode, model.ErrIDInvalid
	}
	now := time.Now().Unix()
	if a.Time < now-20 || a.Time > now+20 {
		return TimeOutCode, model.ErrTimeOut
	}
	return OK, nil
}
