package service

import (
	"github.com/miajio/dpsk/internal/model"
	"github.com/miajio/dpsk/internal/repository"
)

type FileService interface {
	CreateFile(file *model.FileModel) error
	IsFileExists(hash string) (bool, error)
}

type fileService struct {
	repo repository.FileRepository
}

// NewFileService 创建文件服务
func NewFileService(repo repository.FileRepository) FileService {
	return &fileService{repo: repo}
}

// CreateFile 创建文件
func (s *fileService) CreateFile(file *model.FileModel) error {
	return s.repo.CreateFile(file)
}

// IsFileExists 判断文件是否存在
func (s *fileService) IsFileExists(hash string) (bool, error) {
	file, err := s.repo.GetFileByHash(hash)
	if err != nil {
		return false, err
	}
	if file == nil {
		return false, nil
	}
	return true, nil
}
