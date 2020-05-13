package mongo

import (
	"github.com/Naist4869/awesomeProject/model/wxmodel"
	"github.com/Naist4869/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkWxClient struct {
	Client      *mongo.Client
	DB          string                       // 数据库名
	collections map[string]*mongo.Collection // 数据表map
	Logger      *log.Logger
}

func NewWorkWxClient(client *mongo.Client, DB string, logger *log.Logger) *WorkWxClient {
	return &WorkWxClient{Client: client, DB: DB, Logger: logger}
}

func (wc *WorkWxClient) Insert(u *wxmodel.UserInfo) (err error) {
	panic("implement me")
}
func (wc *WorkWxClient) Init() error {
	return nil
}
