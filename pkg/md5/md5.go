package md5

import (
	"crypto/md5"
	"fmt"
)

// Md5 md5加密
func Md5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	hash2 := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash2)
}
