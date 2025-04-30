package jwt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/miajio/dpsk/pkg/jwt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func TestJWT(t *testing.T) {
	// 初始化JWT工具
	jwtUtil := jwt.NewJWTUtil("your-secret-key")

	// 示例用户数据
	user := User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
	}

	// 生成token
	token, err := jwtUtil.GenerateToken(user, time.Now().Add(24*time.Hour))
	if err != nil {
		fmt.Println("生成token失败:", err)
		return
	}
	fmt.Println("生成的token:", token)

	// 验证token
	err = jwtUtil.ValidateToken(token)
	if err != nil {
		fmt.Println("token验证失败:", err)
	} else {
		fmt.Println("token验证成功")
	}

	// 解析token
	var decodedUser User
	err = jwtUtil.ParseToken(token, &decodedUser)
	if err != nil {
		fmt.Println("解析token失败:", err)
		return
	}
	fmt.Printf("解析出的用户数据: %+v\n", decodedUser)

	// 测试过期token
	expiredToken, _ := jwtUtil.GenerateToken(user, time.Now().Add(-1*time.Hour))
	err = jwtUtil.ValidateToken(expiredToken)
	if err == jwt.ErrExpiredToken {
		fmt.Println("token已过期(符合预期)")
	}
}
