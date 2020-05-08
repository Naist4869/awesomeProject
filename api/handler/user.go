package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Naist4869/awesomeProject/usecase"

	"github.com/Naist4869/awesomeProject/model/usermodel"
)

type UserController struct {
	Registration usecase.IUserRegistration
	View         IView
}

func NewUserController(registration usecase.IUserRegistration, view IView) *UserController {
	return &UserController{Registration: registration, View: view}
}

func (c UserController) register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			c.View.NewServerErrorResponse(ErrMethodNotAllowed, w)
			return
		}
		user := &usermodel.RegisterArgument{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			c.View.NewMarshalErrorResponse(err, w)
			return
		}
		if err := user.Validate(); err != nil {
			c.View.NewArgumentErrorResponse(err, w)
			return
		}
		if err := c.Registration.RegisterUser(*user); err != nil {
			c.View.NewServerErrorResponse(err, w)
			return
		}
		c.View.NewOKResponse("注册成功", w)
	})

}
func (c UserController) MakeUserHandler(r *http.ServeMux) {
	r.Handle("/api/vi/user/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.View.NewOKResponse("pong", w)
	}))
	r.Handle("/api/v1/user/register", c.register())
}

type IView interface {
	NewOKResponse(data interface{}, w http.ResponseWriter)
	NewMarshalErrorResponse(err error, w http.ResponseWriter)
	NewArgumentErrorResponse(err error, w http.ResponseWriter)
	NewServerErrorResponse(err error, w http.ResponseWriter)
	NewAuthError(operation string, w http.ResponseWriter)
}
