package dto

type LoginResponse struct {
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
}
