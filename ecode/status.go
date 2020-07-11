package ecode

import (
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"

	"github.com/go-kratos/kratos/pkg/ecode/types"
	"github.com/golang/protobuf/ptypes"
)

// Error new status with code and message
func Error(code Code, message string) *Status {
	return &Status{s: &types.Status{Code: int32(code.Code()), Message: message}}
}

// Errorf new status with code and message
func Errorf(code Code, format string, args ...interface{}) *Status {
	return Error(code, fmt.Sprintf(format, args...))
}

var _ Codes = &Status{}

// Status statusError is an alias of a status proto
// implement ecode.Codes
type Status struct {
	s *types.Status
}

func (s Status) Error() string {
	return s.Message()
}

func (s Status) Code() int {
	return int(s.s.Code)
}

func (s Status) Message() string {
	if s.s.Message == "" {
		return strconv.Itoa(int(s.s.Code))
	}
	return s.s.Message
}

func (s Status) Details() (details []interface{}) {
	if s.s == nil || s.s.Details == nil || len(s.s.Details) == 0 {
		return
	}
	details = make([]interface{}, 0, len(s.s.Details))
	for _, any := range s.s.Details {
		detail := &ptypes.DynamicAny{}
		if err := ptypes.UnmarshalAny(any, detail); err != nil {
			details = append(details, err)
			continue
		}
		details = append(details, detail.Message)
	}
	return details
}

// WithDetails WithDetails
func (s *Status) WithDetails(pbs ...proto.Message) (*Status, error) {
	for _, pb := range pbs {
		anyMsg, err := ptypes.MarshalAny(pb)
		if err != nil {
			return s, err
		}
		s.s.Details = append(s.s.Details, anyMsg)
	}
	return s, nil
}
