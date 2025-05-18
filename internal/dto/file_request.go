package dto

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"

	"github.com/miajio/dpsk/internal/errors"
)

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	BusinessType string                `form:"business_type" binding:"required"` // 业务类型
	File         *multipart.FileHeader `form:"file" binding:"required"`          // 文件
}

func (f *FileUploadRequest) GetHash() (string, error) {
	if f.File == nil {
		return "", errors.ErrFileNotFound
	}
	file, err := f.File.Open()
	if err != nil {
		return "", errors.ErrFileOpenFailed
	}
	defer file.Close()
	// sha256
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", errors.ErrHashCalculationFailed
	}
	hashSum := hex.EncodeToString(hash.Sum(nil))
	return hashSum, nil
}
