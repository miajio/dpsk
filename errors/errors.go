package errors

import (
	"encoding/json"
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

type CodeError struct {
	Code    int    `json:"code,omitempty"` // 错误码
	Message string `json:"message"`        // 错误信息
}

// Error 当错误码为0时直接返回错误信息, 否则判断错误输出类型进行返回错误信息
func (e *CodeError) Error() string {
	if e.Code == 0 {
		return e.Message
	}
	switch ErrPrintTypeDefault {
	case ErrPrintTypeJson:
		bytes, _ := json.Marshal(e)
		return string(bytes)
	default:
		return fmt.Sprintf("error code: %d message: %s", e.Code, e.Message)
	}
}

// New 创建错误
func New(msg string) error {
	return &CodeError{Code: 0, Message: msg}
}

// NewF 创建错误
func NewF(format string, a ...any) error {
	return &CodeError{Code: 0, Message: fmt.Sprintf(format, a...)}
}

// NewCodeError 创建错误
func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Message: msg}
}

// NewCodeErrorF 创建错误
func NewCodeErrorF(code int, format string, a ...any) error {
	return NewCodeError(code, fmt.Sprintf(format, a...))
}

// ReadCodeError 读取CodeError
func ReadCodeError(err error) *CodeError {
	if err == nil {
		return nil
	}
	if codeErr, ok := err.(*CodeError); ok {
		return codeErr
	}
	return nil
}
