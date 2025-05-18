package repository

import (
	"github.com/miajio/dpsk/internal/errors"
	"github.com/miajio/dpsk/internal/model"
	"github.com/miajio/dpsk/pkg/bcrypt"
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
	return &userRepository{db: db}
}

func (r *userRepository) Register(account, nickname, password string) error {
	// 校验账号是否存在
	var userModel model.UserModel
	if err := r.db.Where("account = ?", account).First(&userModel).Error; err == nil {
		return errors.ErrAccountExisted
	}

	if hashPassword, err := bcrypt.Gen(password); err != nil {
		return err
	} else {
		password = hashPassword
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
		return nil, errors.ErrAccountNotExisted
	}
	if !bcrypt.Check(password, userModel.Password) {
		return nil, errors.ErrPassword
	}
	return &userModel, nil
}
