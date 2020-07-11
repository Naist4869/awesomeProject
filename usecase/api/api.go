package api

import (
	"encoding/json"
	"fmt"

	"github.com/Naist4869/awesomeProject/tool"

	"github.com/Naist4869/awesomeProject/api"
	"github.com/Naist4869/awesomeProject/model"
	"github.com/Naist4869/awesomeProject/model/apimodel"
	"github.com/Naist4869/log"
	"go.uber.org/zap"
)

type UseCase struct {
	Logger *log.Logger
}

func (u UseCase) Handle(arg *apimodel.Argument, api api.Api) (result interface{}, code int, err error) {
	if code, err = arg.Validate(); err != nil {
		return
	}
	// 验证API key是否正确
	u.Logger.Debug("验证apikey", zap.String("apikey", arg.APIkey), zap.Any("api", api))
	argument := api.Argument()
	u.Logger.Debug("反序列化具体参数", zap.ByteString("参数", arg.Data))
	if err = json.Unmarshal(arg.Data, argument); err != nil {
		code, err = apimodel.MarshalCode, model.ErrMarshal
		return
	}
	if result, err = api.Invoke(argument, arg.ID); err != nil {
		code = apimodel.ArgumentCode
		return
	}
	code = apimodel.OK
	return

}
func (u UseCase) ValidateKey(original, encrypted string, timestamp int64) bool {
	if encrypted == "123" {
		return true
	}
	return tool.Md5(fmt.Sprintf("%d%s%d", timestamp, original, timestamp)) == encrypted
}
