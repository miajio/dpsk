package bcrypt

import "golang.org/x/crypto/bcrypt"

// Gen 生成加密字符串
func Gen(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// Check 检查加密字符串是否匹配
func Check(password, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(inputPassword), []byte(password))
	return err == nil
}
