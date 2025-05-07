package dto

import (
	"github.com/go-playground/validator/v10"
)

type UploadFileRequest struct {
	FileType string `form:"file_type" validate:"required"`
}

type ChunkUploadRequest struct {
	FileID      int64  `form:"file_id" validate:"omitempty"`
	UploadID    string `form:"upload_id" validate:"omitempty"`
	ChunkNumber int    `form:"chunk_number" validate:"required,min=1"`
	TotalChunks int    `form:"total_chunks" validate:"required,min=1"`
	FileType    string `form:"file_type" validate:"required"`
	FileName    string `form:"file_name" validate:"required"`
	FileSize    int64  `form:"file_size" validate:"required,min=1"`
	ChunkSize   int64  `form:"chunk_size" validate:"required,min=1"`
	Hash        string `form:"hash" validate:"required"`
}

type CompleteChunkUploadRequest struct {
	UploadID    string `json:"upload_id" validate:"required"`
	FileType    string `json:"file_type" validate:"required"`
	FileName    string `json:"file_name" validate:"required"`
	FileSize    int64  `json:"file_size" validate:"required,min=1"`
	TotalChunks int    `json:"total_chunks" validate:"required,min=1"`
	Hash        string `json:"hash" validate:"required"`
}

type GetFileRequest struct {
	ID int64 `uri:"id" validate:"required"`
}

type ListFilesRequest struct {
	RelativePathPrefix string `form:"relative_path_prefix"`
	ID                 int64  `form:"id"`
	Status             string `form:"status" validate:"omitempty,oneof=NORMAL DISABLE"`
	IsDeleted          string `form:"is_deleted" validate:"omitempty,oneof=NORMAL DELETED"`
	Page               int    `form:"page" validate:"min=1"`
	PageSize           int    `form:"page_size" validate:"min=1,max=100"`
}

type DownloadFileRequest struct {
	ID int64 `uri:"id" validate:"required"`
}

type DeleteFileRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type DisableFileRequest struct {
	ID     int64  `json:"id" validate:"required"`
	Reason string `json:"reason" validate:"required"`
}

var validate = validator.New()

func Validate(s interface{}) error {
	return validate.Struct(s)
}
