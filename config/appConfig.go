package config

import "go.uber.org/zap/zapcore"

const (
	filePath           = "config/config.yaml"
	customFilePath     = "config/config.custom.yaml"
	productionFilePath = "config/config.production.yaml"
)

// AppConfig represents the application config
type AppConfig struct {
	SQLConfig       DataStoreConfig `yaml:"sqlConfig"`
	CouchdbConfig   DataStoreConfig `yaml:"couchdbConfig"`
	CacheGrpcConfig DataStoreConfig `yaml:"cacheGrpcConfig"`
	UserGrpcConfig  DataStoreConfig `yaml:"userGrpcConfig"`
	MongodbConfig   DataStoreConfig `yaml:"mongodbConfig"`
	ZapConfig       LogConfig       `yaml:"zapConfig"`
	LorusConfig     LogConfig       `yaml:"logrusConfig"`
	Log             LogConfig       `yaml:"logConfig"`
	UseCase         UseCaseConfig   `yaml:"useCaseConfig"`
}

// UseCaseConfig represents different use cases
type UseCaseConfig struct {
	Registration RegistrationConfig `yaml:"registration"`
	ListUser     ListUserConfig     `yaml:"listUser"`
	ListCourse   ListCourseConfig   `yaml:"listCourse"`
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
