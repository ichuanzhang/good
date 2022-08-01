package errcode

import (
	"fmt"
	"runtime"
	//"github.com/pkg/errors"
)

// Error 结构体
type Error struct {
	status int
	code   int
	msg    string
	err    error
	stack  *stack
}

type stack struct {
	file     string
	function string
	line     int
}

// Error 实现 interface error
func (e *Error) Error() string {
	return fmt.Sprintf("status=%d, code=%d, msg=%s, err=%v, stack=%v", e.status, e.code, e.msg, e.err, e.stack)
}

// Wrap 包装错误
func (e *Error) Wrap(err error) *Error {
	return &Error{
		status: e.status,
		code:   e.code,
		msg:    e.msg,
		err:    err,
	}
}

// Unwrap 实现 interface Unwrap
func (e *Error) Unwrap() error {
	return e.err
}

// WithStack 添加调用者堆栈
func (e *Error) WithStack() *Error {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}
	function := runtime.FuncForPC(pc).Name()
	s := stack{
		file:     file,
		line:     line,
		function: function,
	}
	return &Error{
		status: e.status,
		code:   e.code,
		msg:    e.msg,
		err:    e.err,
		stack:  &s,
	}
}

// AddMsg 追加 msg
// 返回一个新的错误，不会改变原始错误
func (e *Error) AddMsg(format string, args ...interface{}) *Error {
	var msg string
	if format != "" && args != nil {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	return &Error{
		status: e.status,
		code:   e.code,
		msg:    e.msg + " " + msg,
		err:    e.err,
	}
}

// SetMsg 设置 msg
// 返回一个新的错误，不会改变原始错误
func (e *Error) SetMsg(format string, args ...interface{}) *Error {
	var msg string
	if format != "" && args != nil {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	return &Error{
		status: e.status,
		code:   e.code,
		msg:    msg,
		err:    e.err,
	}
}

// FillMsg 填充 msg
// 返回一个新的错误，不会改变原始错误
func (e *Error) FillMsg(args ...interface{}) *Error {
	return &Error{
		status: e.status,
		code:   e.code,
		msg:    fmt.Sprintf(e.msg, args...),
		err:    e.err,
	}
}

// IsNotFound 判断是否为 404 错误
func IsNotFound(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.status != StatusNotFound {
		return false
	}
	return true
}

// IsForbidden 判断是否为 403 错误
func IsForbidden(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.status != StatusForbidden {
		return false
	}
	return true
}
