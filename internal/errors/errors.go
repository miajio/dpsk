package errors

import (
	"errors"
	"fmt"
)

var (
	ErrAccountExisted = errors.New("账号已存在")
)

func New(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}
