package usecasefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/dataservicefactory"
	"github.com/Naist4869/awesomeProject/dataservice"
)

func buildUserData(c container.Container, dc *config.DataConfig) (dataservice.IUserDataService, error) {
	dsi, err := dataservicefactory.GetDataServiceFb(dc.Code)(c, dc)
	if err != nil {
		return nil, err
	}
	udi := dsi.(dataservice.IUserDataService)
	return udi, nil
}
