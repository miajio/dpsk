package dto

import "github.com/miajio/dpsk/models"

type FileResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Extension     string `json:"extension"`
	RelativePath  string `json:"relative_path"`
	AbsoluteURL   string `json:"absolute_url"`
	Hash          string `json:"hash"`
	Size          int64  `json:"size"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Status        string `json:"status"`
	IsDeleted     string `json:"isDeleted"`
	DisableReason string `json:"disable_reason,omitempty"`
}

func ToFileResponse(file *models.SysFile) FileResponse {
	return FileResponse{
		ID:            file.ID,
		Name:          file.Name,
		Type:          file.Type,
		Extension:     file.Extension,
		RelativePath:  file.RelativePath,
		AbsoluteURL:   file.AbsoluteURL,
		Hash:          file.Hash,
		Size:          file.Size,
		CreatedAt:     file.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     file.UpdatedAt.Format("2006-01-02 15:04:05"),
		Status:        string(file.Status),
		IsDeleted:     string(file.IsDeleted),
		DisableReason: file.DisableReason,
	}
}

type FileListResponse struct {
	Total int64          `json:"total"`
	Files []FileResponse `json:"files"`
}

type UploadResponse struct {
	ID       int64  `json:"id"`
	Chunked  bool   `json:"chunked"`
	Message  string `json:"message,omitempty"`
	UploadID string `json:"upload_id"`
}

type DownloadResponse struct {
	URL string `json:"url,omitempty"`
}
