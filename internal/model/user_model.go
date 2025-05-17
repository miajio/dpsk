package model

// UserModel 用户模型
type UserModel struct {
	Model
	Account  string `gorm:"type:varchar(32);uniqueIndex;not null"` // 账号
	Nickname string `gorm:"type:varchar(32);not null"`             // 昵称
	Password string `gorm:"type:varchar(256);not null"`            // 密码
}

func (m *UserModel) TableName() string {
	return "user_info"
}
