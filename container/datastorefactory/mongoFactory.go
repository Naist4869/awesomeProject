package datastorefactory

import (
	"context"
	"fmt"
	"time"

	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/dataservice/userdao/mongo"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mongoFactory(c container.Container, dsc config.DataStoreConfig) (interface{}, error) {
	key := dsc.Code
	if value, found := c.Get(key); found {
		return value.(mongo.UserClient), nil
	}
	auth := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		Username:      dsc.User,
		Password:      dsc.Pass,
		AuthSource:    dsc.DB,
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongodriver.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/platform", dsc.Host, dsc.Port)).SetAuth(auth))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}
