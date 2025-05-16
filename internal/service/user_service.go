package service

import (
	"github.com/miajio/dpsk/internal/model"
	"github.com/miajio/dpsk/internal/repository"
)

type UserService interface {
	Login(account, password string) (*model.UserModel, error)
	Register(account, nickname, password string) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Login(account, password string) (*model.UserModel, error) {
	return s.repo.Login(account, password)
}

func (s *userService) Register(account, nickname, password string) error {
	return s.repo.Register(account, nickname, password)
}
