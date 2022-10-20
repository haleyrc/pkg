package errors

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
)

type Coder interface {
	Code() int
}

type Messager interface {
	Message() string
}

func ErrorCode(err error) int {
	root := errors.Unwrap(err)
	return errorCode(root)
}

func errorCode(err error) int {
	if err, ok := err.(Coder); ok {
		return err.Code()
	}
	return http.StatusInternalServerError
}

func ErrorMessage(err error) string {
	root := errors.Unwrap(err)
	if root == nil {
		root = err
	}
	return errorMessage(root)
}

func errorMessage(err error) string {
	if err, ok := err.(Messager); ok {
		return err.Message()
	}
	return "Something went wrong."
}

func NewRawError(code int, msg string) *RawError {
	return &RawError{
		code: code,
		msg:  msg,
	}
}

type RawError struct {
	code int
	msg  string
}

func (re *RawError) Code() int { return re.code }

func (re *RawError) Error() string { return fmt.Sprintf("error: %s", re.msg) }

func (re *RawError) Message() string { return re.msg }

type wrappedError struct {
	file   string
	lineno int
	name   string
	cause  error
}

func (we wrappedError) Code() int { return ErrorCode(we.cause) }

func (we wrappedError) Error() string { return fmt.Sprintf("%s: %v", we.name, we.cause) }

func (we wrappedError) Message() string { return ErrorMessage(we.cause) }

func (we wrappedError) Unwrap() error { return we.cause }

func WrapError(err error) error {
	pc, file, lineno, ok := runtime.Caller(1)

	we := wrappedError{
		file:   filepath.Base(file),
		lineno: lineno,
		cause:  err,
	}

	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		we.name = filepath.Base(details.Name())
	}

	return we
}
