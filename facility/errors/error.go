package errors

import (
	pkgerror "github.com/pkg/errors"
)

type Error struct {
	raw  error
	code string
}

func (e *Error) Error() string {
	return e.raw.Error()
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Raw() error {
	return e.raw
}

type readable interface {
	Code() string
	Raw() error
	Error() string
}

func Readable(err error) (readable, bool) {
	r, ok := err.(readable)
	return r, ok
}

// NewWithCodef use error.Error as result, will contains error code
// to distinguish differences between errors
func NewWithCodef(code, format string, args ...interface{}) error {
	return &Error{
		Newf(nil, format, args...),
		code,
	}
}

func WrapWithCodef(err error, code, format string, args ...interface{}) error {
	return &Error{
		Newf(err, format, args...),
		code,
	}
}

// Newf is an adapt method of github.com/pkg/errors.Wrapf
// in order to avoid confusing of which method in github.com/pkg/errors can be used to
// fit for the error logging and responsing
func Newf(err error, format string, args ...interface{}) error {
	if err != nil {
		return pkgerror.Wrapf(err, format, args)
	}
	return pkgerror.Errorf(format, args)
}

func Panicf(err error, format string, args ...interface{}) {
	panic(pkgerror.Wrapf(err, format, args))
}
