package api

// 本文件定义了全部的对外 API Code，每个模块具有100个可用范围，即从X00->X99
//	用户为100->199
//	订单为200->299
//	钱包为300->399
//  计算为400->499
//	推送为500->599
//	风控为600->699

// 生成 gRPC 代码
//go:generate kratos tool protoc --grpc  --swagger api.proto
const (
	APIGetSecret = "200"
)

type Api interface {
	Invoke(argument interface{}, merchantID int64) (interface{}, error) //调用
	Argument() interface{}                                              //参数
	Name() string                                                       //名称
	Key() string                                                        //接口的key
}

type ApiProvider interface {
	HTTPApis() []Api
}

type api struct {
	name     string
	key      string
	argument func() interface{}
	fun      func(argument interface{}, merchantID int64) (interface{}, error)
}

func NewApi(name, key string, argument func() interface{}, fun func(argument interface{}, merchantID int64) (interface{}, error)) Api {
	return api{
		name:     name,
		key:      key,
		argument: argument,
		fun:      fun,
	}
}

func (a api) Invoke(argument interface{}, merchantID int64) (interface{}, error) {
	return a.fun(argument, merchantID)
}

func (a api) Argument() interface{} {
	return a.argument()
}

func (a api) Name() string {
	return a.name
}

func (a api) Key() string {
	return a.key
}
