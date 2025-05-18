package errors

import (
	"errors"
	"fmt"
)

var (
	ErrAccountExisted        = errors.New("账号已存在")
	ErrAccountNotExisted     = errors.New("账号不存在")
	ErrPassword              = errors.New("密码错误")
	ErrNoToken               = errors.New("未提供认证令牌")
	ErrInvalidToken          = errors.New("认证令牌无效")
	ErrTokenParseFailed      = errors.New("令牌解析失败")
	ErrTokenRefreshFailed    = errors.New("令牌续期失败")
	ErrFileNotFound          = errors.New("文件不存在")
	ErrFileOpenFailed        = errors.New("文件打开失败")
	ErrHashCalculationFailed = errors.New("哈希计算失败")
	ErrFileExists            = errors.New("文件已存在")
	ErrFileSuffixNotFound    = errors.New("文件后缀不存在")
)

func New(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}
