package errors

import (
	"errors"
	"fmt"
)

var (
	ErrAccountExisted    = errors.New("账号已存在")
	ErrAccountNotExisted = errors.New("账号不存在")
	ErrPasswordError     = errors.New("密码错误")
)

func New(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}
