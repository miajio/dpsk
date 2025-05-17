package dto

import "encoding/json"

type LoginUser struct {
	Id       int64  `json:"id"`       // 用户ID
	Account  string `json:"account"`  // 用户账号
	Nickname string `json:"nickname"` // 用户昵称
}

func (l *LoginUser) Marshal() string {
	b, _ := json.Marshal(l)
	return string(b)
}

// LoginResponse 登录响应
type LoginResponse struct {
	LoginUser `json:",inline"` // 登录用户信息
	Token     string           `json:"token"` // 登录Token
}
