package restserver

import (
	"errors"
	"net/http"

	"github.com/Naist4869/awesomeProject/api"

	"github.com/Naist4869/awesomeProject/model/apimodel"
	"github.com/Naist4869/awesomeProject/usecase"
)

type ApiServer struct {
	apiUseCase usecase.IAPI
	limit      int
	apis       map[string]api.Api
	handler    IRouter // http句柄
}

func NewApiServer(apiUseCase usecase.IAPI, limit int, handler IRouter) *ApiServer {
	apiServer := &ApiServer{apiUseCase: apiUseCase, limit: limit, apis: make(map[string]api.Api, 10), handler: handler}
	apiServer.Http()
	return apiServer
}
func (s *ApiServer) Http() {
	s.handler.POST("/invoke", s.invoke)
}
func (s *ApiServer) invoke(ctx *Context) {
	arg := &apimodel.Argument{}
	if err := ctx.ShouldBindJSON(arg); err != nil {
		ctx.Error(err)
		return
	}

	if Iapi, exist := s.apis[arg.APIkey]; !exist {
		// error
		return
	} else {
		if result, _, err := s.apiUseCase.Handle(arg, Iapi); err != nil {
			ctx.Error(err)
			return
		} else {
			ctx.JSON(http.StatusOK, struct {
				Error     string
				ErrorCode int
				Data      interface{}
			}{Data: result})
		}
	}

}
func (s *ApiServer) Register(api api.Api) error {
	if _, exist := s.apis[api.Key()]; exist {
		return errors.New("已注册")
	}
	s.apis[api.Key()] = api
	return nil
}
