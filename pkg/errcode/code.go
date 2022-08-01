package errcode

// 状态码
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusServiceUnavailable  = 503
)

// 错误码
const (
	CodeOK                  = 1
	CodeInvalidParam        = 10001
	CodeForbidden           = 10002
	CodeNotFound            = 10003
	CodeInternalServerError = 10004
	CodeUnavailable         = 10005
)

// 错误
var (
	OK                = register(StatusOK, CodeOK, "OK")
	ErrParam          = register(StatusBadRequest, CodeInvalidParam, "参数错误")
	ErrForbidden      = register(StatusForbidden, CodeForbidden, "无权访问")
	ErrNotFound       = register(StatusNotFound, CodeNotFound, "内容不存在")
	ErrInternalServer = register(StatusInternalServerError, CodeInternalServerError, "服务器内部错误")
	ErrUnavailable    = register(StatusServiceUnavailable, CodeUnavailable, "服务暂不可用，请稍后重试")
)

// codeMap 全部错误码
var codeMap = make(map[int]*Error)

// register 注册错误
func register(status, code int, msg string) *Error {
	err := &Error{
		status: status,
		code:   code,
		msg:    msg,
	}
	if _, ok := codeMap[code]; ok {
		panic("error code register failed, code exist")
	}
	codeMap[code] = err
	return err
}
