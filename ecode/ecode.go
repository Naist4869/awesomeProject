package ecode

import (
	"fmt"
	"strconv"
	"sync/atomic"
)

func init() {
}

var (
	_messages atomic.Value         // NOTE: stored map[string]map[int]string
	_codes    = map[int]struct{}{} // register codes.
)

// Register register ecode message map.
func Register(cm map[int]string) {
	_messages.Store(cm)
}

// New new a ecode.Codes by int value.
// NOTE: ecode must unique in global, the New will check repeat and then panic.
func New(e int) Code {
	if e <= 0 {
		panic("business ecode must greater than zero")
	}
	return add(e)
}
func add(e int) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Int(e)
}

type Codes interface {
	Error() string
	Code() int
	Message() string
	Details() []interface{}
}

type Code int

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

func (e Code) Code() int {
	return int(e)
}

func (e Code) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[e.Code()]; ok {
			return msg
		}
	}
	return e.Error()
}

// Int parse code int to error.
func Int(i int) Code { return Code(i) }

// String parse code string to error.
func String(e string) Code {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return ServerErr
	}
	return Code(i)
}

// Details return details.
func (e Code) Details() []interface{} { return nil }
