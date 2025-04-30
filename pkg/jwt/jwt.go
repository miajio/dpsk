package jwt

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTUtil JWT工具结构体
type JWTUtil struct {
	secretKey []byte // 密钥
}

// NewJWTUtil 创建一个新的JWTUtil实例
func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken 生成JWT token
// data: 任意结构体数据
// expiresAt: 过期时间
func (j *JWTUtil) GenerateToken(data interface{}, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"data": data,
		"exp":  expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken 验证token是否有效且未过期
func (j *JWTUtil) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ErrExpiredToken
		}
		return ErrInvalidToken
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

// ParseToken 解析token并将数据写入目标结构体
func (j *JWTUtil) ParseToken(tokenString string, target interface{}) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ErrExpiredToken
		}
		return ErrInvalidToken
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrInvalidToken
	}

	data, ok := claims["data"]
	if !ok {
		return ErrInvalidToken
	}

	// 将data转换为map[string]interface{}，然后使用json转换到目标结构体
	// 这种方式可以处理任意结构体
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return ErrInvalidToken
	}

	// 使用第三方库如 mapstructure 可以更方便地转换
	// 这里使用json的marshal/unmarshal作为示例
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, target)
}
