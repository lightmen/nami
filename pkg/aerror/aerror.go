package aerror

import (
	"fmt"

	"github.com/lightmen/nami/codes"
)

type Error interface {
	Code() int32
	Error() string
}

type aError struct {
	code int32
	msg  string
}

func New(code int32, msg string) Error {
	return &aError{code: int32(code), msg: msg}
}

func (e *aError) Code() int32 {
	return e.code
}

func (e *aError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.code, e.msg)
}

func Code(err error) int32 {
	if err == nil {
		return 0
	}

	ae, ok := err.(Error)
	if ok {
		return ae.Code()
	}

	return int32(codes.Unknown)
}
