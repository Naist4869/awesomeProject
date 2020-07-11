package config

import (
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	filePath           = "config/config.yaml"
	customFilePath     = "config/config.custom.yaml"
	productionFilePath = "config/config.production.yaml"
)

// AppConfig represents the application config
type AppConfig struct {
	SQLConfig                 DataStoreConfig `yaml:"sqlConfig"`
	CouchdbConfig             DataStoreConfig `yaml:"couchdbConfig"`
	CacheGrpcConfig           DataStoreConfig `yaml:"cacheGrpcConfig"`
	UserGrpcConfig            DataStoreConfig `yaml:"userGrpcConfig"`
	MongodbConfig             DataStoreConfig `yaml:"mongodbConfig"`
	FileSystemRpcClientConfig RpcClientConfig `yaml:"fileSystemRpcClientConfig"`
	TBKRpcClientConfig        RpcClientConfig `yaml:"tbkRpcClientConfig"`
	ZapConfig                 LogConfig       `yaml:"zapConfig"`
	LorusConfig               LogConfig       `yaml:"logrusConfig"`
	Log                       LogConfig       `yaml:"logConfig"`
	PrometheusConfig          MetricsConfig   `yaml:"metricsConfig"`
	UseCase                   UseCaseConfig   `yaml:"useCaseConfig"`
}

// UseCaseConfig represents different use cases
type UseCaseConfig struct {
	Registration RegistrationConfig `yaml:"registration"`
	ListUser     ListUserConfig     `yaml:"listUser"`
	ListCourse   ListCourseConfig   `yaml:"listCourse"`
	WorkWx       WorkWxConfig       `yaml:"workWx"`
	OfficialWx   OfficialWxConfig   `yaml:"officialWx"`
	ApiUseCase   ApiUseCaseConfig   `yaml:"apiUseCase"`
}
type WorkWxConfig struct {
	Code             string     `yaml:"code"`
	CorpID           string     `yaml:"corpID"`
	CorpSecret       string     `yaml:"corpSecret"`
	AgentID          int64      `yaml:"agentID"`
	WorkWxDataConfig DataConfig `yaml:"workWxDataConfig"`
}

type ApiUseCaseConfig struct {
	Limit int `yaml:"limit"`
}
type OfficialWxConfig struct {
	Code         string `yaml:"code"`
	OriID        string `yaml:"oriID"`        //可选; 公众号的原始ID(微信公众号管理后台查看), 如果设置了值则该Server只能处理 ToUserName 为该值的公众号的消息(事件);
	AppID        string `yaml:"appID"`        // 可选; 公众号的AppId, 如果设置了值则安全模式时该Server只能处理 AppId 为该值的公众号的消息(事件);
	Token        string `yaml:"token"`        //     必须; 公众号用于验证签名的token;
	Base64AESKey string `yaml:"base64AESKey"` // 可选; aes加密解密key, 43字节长(base64编码, 去掉了尾部的'='), 安全模式必须设置;
	Secret       string `yaml:"secret"`       // 应用 secret

	FileSystemConfig RpcConfig `yaml:"fileSystemGrpcConfig"` // 文件管理系统Grpc配置
	TBKConfig        RpcConfig `yaml:"tbkGrpcConfig"`        // 淘宝客系统Grpc配置
}

// RegistrationConfig represents registration use case
type RegistrationConfig struct {
	Code           string     `yaml:"code"`
	UserDataConfig DataConfig `yaml:"userDataConfig"`
	TxDataConfig   DataConfig `yaml:"txDataConfig"`
}

// ListUserConfig represents list user use case
type ListUserConfig struct {
	Code            string     `yaml:"code"`
	UserDataConfig  DataConfig `yaml:"userDataConfig"`
	CacheDataConfig DataConfig `yaml:"cacheDataConfig"`
}

// ListCourseConfig represents list course use case
type ListCourseConfig struct {
	Code             string     `yaml:"code"`
	CourseDataConfig DataConfig `yaml:"courseDataConfig"`
}

// DataConfig represents data service
type DataConfig struct {
	Code            string          `yaml:"code"`
	DataStoreConfig DataStoreConfig `yaml:"dataStoreConfig"`
}

// DataConfig represents handlers for data store. It can be a database or a gRPC connection
type DataStoreConfig struct {
	Code string `yaml:"code"`
	// Only database has a driver name, for grpc it is "tcp" ( network) for server
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Pass       string `yaml:"pass"`
	DB         string `yaml:"db"`
}

type RpcConfig struct {
	Code   string          `yaml:"code"`
	Client RpcClientConfig `yaml:"client"`
}

type RpcClientConfig struct {
	Code        string        `yaml:"code"` // 目前只有grpc
	DialTimeOut time.Duration `yaml:"dialTimeOut"`
	Target      string        `yaml:"target"`
}

// LogConfig represents logger handler
// Logger has many parameters can be set or changed. Currently, only three are listed here. Can add more into it to
// fits your needs.
type LogConfig struct {
	// log library name
	Code string `yaml:"code"`

	// 最大大小,单位为MB
	MaxSize int `yaml:"maxSize"`
	// 最长时间,单位为天
	MaxAge int `yaml:"maxAge"`
	// 日志目录
	LogDir string `yaml:"logDir"`
	// 日志名
	Name string `yaml:"name"`
	// 是否输出至终端
	Debug    bool                     `yaml:"debug"`
	Console  bool                     `yaml:"console"`
	MinLevel map[string]zapcore.Level `yaml:"minLevel"`
}

type MetricsConfig struct {
	Code string `yaml:"code"`
}
