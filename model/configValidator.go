package model

// constant for logger code, it needs to match log code (logConfig)in configuration
const (
	LOGRUS string = "logrus"
	ZAP    string = "zap"
)

// data service code. Need to map to the data service code (DataConfig) in the configuration yaml file.
const (
	USER_DATA   string = "userData"
	CACHE_DATA  string = "cacheData"
	TX_DATA     string = "txData"
	COURSE_DATA string = "courseData"
	WORKWX_DATA string = "workWxData"
)
const (
	SQLDB           string = "sqldb"
	COUCHDB         string = "couch"
	CACHE_GRPC      string = "cacheGrpc"
	USER_GRPC       string = "userGrpc"
	FileSystem_GRPC string = "fileSystemGrpc"
	TBK_GRPC        string = "tbkGrpc"
	MONGO           string = "mongodb"
	GRPC            string = "grpc"
)

// use case code. Need to map to the use case code (UseCaseConfig) in the configuration yaml file.
// Client app use those to retrieve use case from the container
const (
	REGISTRATION string = "registration"
	LIST_USER    string = "listUser"
	LIST_COURSE  string = "listCourse"
	WORKWX       string = "workWx"
	OfficialWx   string = "officialWx"
	Api          string = "api"
)
const (
	Prometheus string = "prometheus"
)
