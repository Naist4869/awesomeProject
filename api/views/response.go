package views

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Naist4869/awesomeProject/dataservice"

	"github.com/Naist4869/awesomeProject/api/handler"

	"github.com/Naist4869/awesomeProject/model/usermodel"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

const (
	// MarshalJSONError 解析json参数错误
	MarshalJSONError = iota + 1
	// ArgumentError 参数错误
	ArgumentError
	// ServerError 服务器错误
	ServerError
	// NeedLoginError 需要重新登录
	NeedLoginError
	// AuthError 权限不足
	AuthError
)

// Response 请求应答
type Response struct {
	Data     interface{} `json:"data"`     // 数据
	Error    int         `json:"error"`    // 错误码
	ErrorMsg string      `json:"errorMsg"` // 错误信息
}

type View struct{}

var ErrHTTPStatusMap = map[string]int{
	usermodel.ErrUserPhoneEmpty.Error():    http.StatusBadRequest,
	usermodel.ErrUserNickNameEmpty.Error(): http.StatusBadRequest,
	handler.ErrMethodNotAllowed.Error():    http.StatusBadRequest,
	dataservice.ErrInsertFailed.Error():    http.StatusInternalServerError,
	dataservice.ErrPhoneExists.Error():     http.StatusConflict,
	// pkg.ErrNotFound.Error():     http.StatusNotFound,
	// pkg.ErrInvalidSlug.Error():  http.StatusBadRequest,
	// pkg.ErrExists.Error():       http.StatusConflict,
	// pkg.ErrNoContent.Error():    http.StatusNotFound,
	// pkg.ErrDatabase.Error():     http.StatusInternalServerError,
	// pkg.ErrUnauthorized.Error(): http.StatusUnauthorized,
	// pkg.ErrForbidden.Error():    http.StatusForbidden,
	// ErrMethodNotAllowed.Error(): http.StatusMethodNotAllowed,
	// ErrInvalidToken.Error():     http.StatusBadRequest,
	// ErrUserExists.Error():       http.StatusConflict,
}

func (View) NewOKResponse(data interface{}, w http.ResponseWriter) {
	writeContentType(w)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&Response{
		Data: data,
	})
}

func writeContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = jsonContentType
	}
}
func (View) NewMarshalErrorResponse(err error, w http.ResponseWriter) {
	msg := err.Error()
	code, ok := ErrHTTPStatusMap[msg]
	if !ok {
		code = http.StatusInternalServerError
	}
	writeContentType(w)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&Response{
		Error:    MarshalJSONError,
		ErrorMsg: err.Error(),
	})
}
func (View) NewArgumentErrorResponse(err error, w http.ResponseWriter) {
	msg := err.Error()
	code, ok := ErrHTTPStatusMap[msg]
	if !ok {
		code = http.StatusInternalServerError
	}
	writeContentType(w)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&Response{
		Error:    ArgumentError,
		ErrorMsg: err.Error(),
	})
}

func (View) NewServerErrorResponse(err error, w http.ResponseWriter) {
	msg := err.Error()
	code, ok := ErrHTTPStatusMap[msg]
	if !ok {
		code = http.StatusInternalServerError
	}
	writeContentType(w)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&Response{
		Error:    ServerError,
		ErrorMsg: err.Error(),
	})
}

func (View) NewAuthError(operation string, w http.ResponseWriter) {
	writeContentType(w)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&Response{
		Error:    AuthError,
		ErrorMsg: fmt.Sprintf("当前用户没有[%s]的权限", operation),
	})
}
