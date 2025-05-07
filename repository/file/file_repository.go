package file

import (
	"github.com/miajio/dpsk/comm/ctx"
	"github.com/miajio/dpsk/models"
)

type FileRepositoryTemp interface {
	Create(file *models.SysFile) error
	Update(file *models.SysFile) error
	FindById(id int64) (*models.SysFile, error)
	FindByHashAndType(hash, fileType string) (*models.SysFile, error)
	List(relativePathPrefix string, id int64, status models.FileStatus, isDelete models.DeletedStatus, page, pageSize int) ([]*models.SysFile, int64, error)
	Delete(id int64) error
	Disable(id int64, reason string) error

	CreateChunkUpload(chunk *models.SysChunkFile) error
	GetChunkUpload(uploadID string, chunkNumber int) (*models.SysChunkFile, error)
	GetChunkUploads(uploadID string) ([]*models.SysChunkFile, error)
	DeleteChunkUploads(uploadID string) error
	MarkChunkUploadComplete(uploadID string) error
}

type FileRepositoryImpl struct {
}

var FileRepository FileRepositoryTemp = (*FileRepositoryImpl)(nil)

// Create 创建文件
func (r *FileRepositoryImpl) Create(file *models.SysFile) error {
	return ctx.DB.Create(file).Error
}

// Update 更新文件
func (r *FileRepositoryImpl) Update(file *models.SysFile) error {
	return ctx.DB.Save(file).Error
}

// FindById 根据id查询文件
func (r *FileRepositoryImpl) FindById(id int64) (*models.SysFile, error) {
	var file models.SysFile
	err := ctx.DB.First(&file, id).Error
	return &file, err
}

// FindByHashAndType 根据hash和type查询文件
func (r *FileRepositoryImpl) FindByHashAndType(hash, fileType string) (*models.SysFile, error) {
	var file models.SysFile
	err := ctx.DB.Where("hash = ? AND type = ?", hash, fileType).First(&file).Error
	return &file, err
}

// List 根据条件查询文件列表
// relativePathPrefix: 文件相对路径前缀
// id: 文件id
// status: 文件状态
// isDelete: 是否删除
// page: 页码
// pageSize: 每页大小
// 返回值: 文件列表, 总数, 错误
func (r *FileRepositoryImpl) List(relativePathPrefix string, id int64, status models.FileStatus, isDelete models.DeletedStatus, page, pageSize int) ([]*models.SysFile, int64, error) {
	var files []*models.SysFile
	var total int64

	query := ctx.DB.Model(&models.SysFile{})

	if relativePathPrefix != "" {
		query = query.Where("relative_path LIKE ?", relativePathPrefix+"%")
	}

	if id != 0 {
		query = query.Where("id = ?", id)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if isDelete != "" {
		query = query.Where("is_deleted = ?", isDelete)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&files).Error

	return files, total, err
}

// Delete 删除文件
func (r *FileRepositoryImpl) Delete(id int64) error {
	return ctx.DB.Model(&models.SysFile{}).Where("id = ?", id).Update("is_deleted", models.DeletedDeleted).Error
}

// Disable 禁用文件
func (r *FileRepositoryImpl) Disable(id int64, reason string) error {
	return ctx.DB.Model(&models.SysFile{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         models.FileStatusDisable,
			"disable_reason": reason,
		}).Error
}

func (r *FileRepositoryImpl) CreateChunkUpload(chunk *models.SysChunkFile) error {
	return ctx.DB.Create(chunk).Error
}

func (r *FileRepositoryImpl) GetChunkUpload(uploadID string, chunkNumber int) (*models.SysChunkFile, error) {
	var chunk models.SysChunkFile
	err := ctx.DB.Where("upload_id = ? AND chunk_number = ?", uploadID, chunkNumber).First(&chunk).Error
	return &chunk, err
}

func (r *FileRepositoryImpl) GetChunkUploads(uploadID string) ([]*models.SysChunkFile, error) {
	var chunks []*models.SysChunkFile
	err := ctx.DB.Where("upload_id = ?", uploadID).Order("chunk_number ASC").Find(&chunks).Error
	return chunks, err
}

func (r *FileRepositoryImpl) DeleteChunkUploads(uploadID string) error {
	return ctx.DB.Where("upload_id = ?", uploadID).Delete(&models.SysChunkFile{}).Error
}

func (r *FileRepositoryImpl) MarkChunkUploadComplete(uploadID string) error {
	return ctx.DB.Model(&models.SysChunkFile{}).Where("upload_id = ?", uploadID).Update("is_complete", true).Error
}
