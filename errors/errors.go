package errors

import (
	"fmt"
)

// ErrPrintType 错误打印类型
type ErrPrintType string

const (
	ErrPrintTypeConsole ErrPrintType = "console" // 控制台
	ErrPrintTypeJson    ErrPrintType = "json"    // 文件
)

var (
	ErrPrintTypeDefault = ErrPrintTypeConsole
)

type err struct {
	code int    // 错误码
	msg  string // 错误信息
}

// Error 当错误码为0时直接返回错误信息, 否则判断错误输出类型进行返回错误信息
func (e *err) Error() string {
	if e.code == 0 {
		return e.msg
	}
	switch ErrPrintTypeDefault {
	case ErrPrintTypeJson:
		return fmt.Sprintf("{\"code\":%d,\"message\":\"%s\"}", e.code, e.msg)
	default:
		return fmt.Sprintf("error code: %d message: %s", e.code, e.msg)
	}
}

// New 创建错误
func New(msg string) error {
	return &err{code: 0, msg: msg}
}

// NewF 创建错误
func NewF(format string, a ...any) error {
	return &err{code: 0, msg: fmt.Sprintf(format, a...)}
}

// NewCodeError 创建错误
func NewCodeError(code int, msg string) error {
	return &err{code: code, msg: msg}
}

// NewCodeErrorF 创建错误
func NewCodeErrorF(code int, format string, a ...any) error {
	return NewCodeError(code, fmt.Sprintf(format, a...))
}
