package model

// UserModel 用户模型
type UserModel struct {
	Model
	Account  string `gorm:"column:account;type:varchar(32);uniqueIndex;not null"` // 账号
	Nickname string `gorm:"column:nickName;type:varchar(32);not null"`            // 昵称
	Password string `gorm:"column:password;type:varchar(256);not null"`           // 密码
}

func (m *UserModel) TableName() string {
	return "user_info"
}
