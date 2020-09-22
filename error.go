package errors

import (
	"fmt"
	"runtime"
)

type fundamental struct {
	message string
}

func (f fundamental) Error() string {
	return f.message
}

func New(msg string) error {
	return fundamental{message: msg}
}

func Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

type cause interface {
	Cause(err error) error
}

type withMessage struct {
	error
	message string
}

func (w withMessage) Error() string {
	return w.message + ":" + w.error.Error()
}

func (w withMessage) Cause(err error) error {
	return w.error
}

type withStack struct {
	error
	stack string
}

func (w withStack) Error() string {
	return  w.error.Error()+w.stack
}

func (w withStack) Cause(err error) error {
	return w.error
}

func WithMessage(err error, msg string) error {
	return &withMessage{
		error:   err,
		message: msg,
	}
}

func WithStack(err error) error {
	return &withStack{
		error: err,
		stack: caller(),
	}
}

func Wrap(err error, message string) error {
	msg := &withMessage{
		error:   err,
		message: message,
	}
	stack := &withStack{
		error: msg,
		stack: caller(),
	}
	return stack
}

func Wrapf(err error, message string, args ...interface{}) error {
	msg := &withMessage{
		error:   err,
		message: fmt.Sprintf(message, args...),
	}
	stack := &withStack{
		error: msg,
		stack: caller(),
	}
	return stack
}

func Cause(err error) error {
	for {
		if e, ok := err.(cause); ok {
			err = e.Cause(err)
		} else {
			break
		}
	}
	return err
}

func caller() string {
	pc, file, line, _ := runtime.Caller(2)
	fun := runtime.FuncForPC(pc)
	if fun == nil {
		return fmt.Sprintf("\n%s:%d", file, line)
	}
	return fmt.Sprintf("\n%s:%d %s: ", file, line, fun.Name())
}
