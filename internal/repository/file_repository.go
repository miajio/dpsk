package repository

import (
	"github.com/miajio/dpsk/internal/model"
	"gorm.io/gorm"
)

type FileRepository interface {
	CreateFile(fileModel *model.FileModel) error         // CreateFile 创建文件
	GetFileByHash(hash string) (*model.FileModel, error) // GetFileByHash 根据文件哈希值获取文件
}

type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件仓库
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (f *fileRepository) CreateFile(file *model.FileModel) error {
	return f.db.Create(file).Error
}

func (f *fileRepository) GetFileByHash(hash string) (*model.FileModel, error) {
	var file model.FileModel
	err := f.db.Where("hash = ? and isDeleted != 'DELETED'", hash).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			// 其他错误
			return nil, err
		}
	}
	return &file, nil
}
