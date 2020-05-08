package userdataservicefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/datastorefactory"
	"github.com/Naist4869/awesomeProject/dataservice"
	"github.com/Naist4869/awesomeProject/dataservice/userdao/mongo"
	"github.com/Naist4869/log"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func mongoDataServiceFactory(c container.Container, dataConfig *config.DataConfig) (dataservice.IUserDataService, error) {
	dsc := dataConfig.DataStoreConfig
	dsi, err := datastorefactory.GetDataStoreFb(dsc.Code)(c, dsc)
	if err != nil {
		return nil, err
	}
	client := dsi.(*mongodriver.Client)
	userClient := mongo.NewUserClient(client, dsc.DB, log.BaseLogger.With(zap.String("组件", "mongo")))
	if err = userClient.Init(); err != nil {
		return nil, err
	}
	return userClient, nil
}
