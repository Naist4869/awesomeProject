package restserver

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"

	"github.com/Naist4869/awesomeProject/tool/metadata"

	"github.com/Naist4869/log"

	"github.com/Naist4869/awesomeProject/adapter/restserver/render"
)

const (
	abortIndex    int8 = math.MaxInt8 / 2
	defaultStatus      = http.StatusOK
)

type Context struct {
	context.Context

	Request        *http.Request
	ResponseWriter http.ResponseWriter
	status         int

	// 流量控制
	index    int8
	handlers []HandlerFunc

	//Keys 是一个键/值对,专门用于每个请求的上下文。
	Keys map[string]interface{}
	//Key的读写锁
	keysMutex sync.RWMutex

	// 错误
	Errors errorMsgs
	Params Params

	engine *Engine

	// 用于普罗米修斯统计
	RoutePath string
	method    string
}

func (c *Context) reset() {
	c.Context = nil
	c.index = -1
	c.handlers = nil
	c.Keys = nil
	c.Errors = c.Errors[0:0]
	c.method = ""
	c.RoutePath = ""
	c.Params = c.Params[0:0]
	c.status = defaultStatus
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in godoc.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Abort()
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	c.keysMutex.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
	c.keysMutex.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.keysMutex.RLock()
	value, exists = c.Keys[key]
	c.keysMutex.RUnlock()
	return
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}
func (c *Context) GetBytes(key string) (b []byte) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.([]byte)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.ResponseWriter.WriteHeader(code)
	if code > 0 && c.status != code {
		log.BaseLogger.Warn(fmt.Sprintf("[WARNING] Headers were already written. Wanted to override status code %d with %d", c.status, code))
	}
	c.status = code
}

// Render writes the response headers and calls render.Render to render data.
func (c *Context) Render(code int, r render.Render) {
	if code > 0 {
		c.Status(code)
	}
	r.WriteContentType(c.ResponseWriter)
	if !bodyAllowedForStatus(code) {
		return
	}

	if err := r.Render(c.ResponseWriter); err != nil {
		c.Error(err)
		return
	}
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

// Error attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
func (c *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypeServer,
		}
	}

	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

// XML serializes the given struct as XML into the response body.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(code int, obj interface{}) {
	c.Render(code, render.XML{Data: obj})
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface{}) {
	c.Render(code, render.JSON{Data: obj})
}

// ProtoBuf serializes the given struct as ProtoBuf into the response body.
func (c *Context) ProtoBuf(code int, obj interface{}) {
	c.Render(code, render.ProtoBuf{Data: obj})
}

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	c.Render(-1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.Request,
	})
}

func (c *Context) requestHeader(key string) string {
	return c.Request.Header.Get(key)
}

// RemoteIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
// Notice: metadata.RemoteIP take precedence over X-Forwarded-For and X-Real-Ip
func (c *Context) RemoteIP() (remoteIP string) {
	remoteIP = metadata.String(c, metadata.RemoteIP)
	if remoteIP != "" {
		return
	}

	remoteIP = c.Request.Header.Get("X-Forwarded-For")
	remoteIP = strings.TrimSpace(strings.Split(remoteIP, ",")[0])
	if remoteIP == "" {
		remoteIP = strings.TrimSpace(c.Request.Header.Get("X-Real-Ip"))
	}

	return
}

// Bytes writes some data into the body stream and updates the HTTP code.
func (c *Context) Bytes(code int, contentType string, data ...[]byte) {
	c.Render(code, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}

// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Render(code, render.String{Format: format, Data: values})
}
