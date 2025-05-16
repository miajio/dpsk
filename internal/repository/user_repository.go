package repository

import (
	"github.com/miajio/dpsk/internal/errors"
	"github.com/miajio/dpsk/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(account, nickname, password string) error
	Login(account, password string) (*model.UserModel, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据库操作对象
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Register(account, nickname, password string) error {
	// 校验账号是否存在
	var userModel model.UserModel
	if err := r.db.Where("account = ?", account).First(&userModel).Error; err == nil {
		return errors.ErrAccountExisted
	}
	return r.db.Create(&model.UserModel{
		Account:  account,
		Nickname: nickname,
		Password: password,
	}).Error
}

func (r *userRepository) Login(account, password string) (*model.UserModel, error) {
	var userModel model.UserModel
	if err := r.db.Where("account = ?", account).First(&userModel).Error; err != nil {
		return nil, errors.New("账号不存在")
	}
	if userModel.Password != password {
		return nil, errors.New("密码错误")
	}
	return &userModel, nil
}
