package file

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/miajio/dpsk/comm/ctx"
	"github.com/miajio/dpsk/dto"
	"github.com/miajio/dpsk/models"
	"github.com/miajio/dpsk/repository/file"
)

type FileServiceTemp interface {
	UploadFile(fileHeader *multipart.FileHeader, req dto.UploadFileRequest) (*dto.UploadResponse, error)
	GetFile(id int64) (*dto.FileResponse, error)
	ListFiles(req dto.ListFilesRequest) (*dto.FileListResponse, error)
	DownloadFile(id int64) (string, error)
	StreamDownloadFile(id int64) (string, error)
	DeleteFile(id int64) error
	DisableFile(req dto.DisableFileRequest) error

	UploadChunk(fileHeader *multipart.FileHeader, req dto.ChunkUploadRequest) (*dto.UploadResponse, error)
	CompleteChunkUpload(req dto.CompleteChunkUploadRequest) (*dto.UploadResponse, error)
}

type FileServiceImpl struct {
}

var FileService FileServiceTemp = (*FileServiceImpl)(nil)

func (s *FileServiceImpl) UploadFile(fileHeader *multipart.FileHeader, req dto.UploadFileRequest) (*dto.UploadResponse, error) {
	// 检查文件大小
	if fileHeader.Size > ctx.File.MinChunkSize {
		return nil, errors.New("file too large, please use chunk upload")
	}

	// 打开文件
	fileObj, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer fileObj.Close()

	// 计算文件hash
	hash := md5.New()
	if _, err := io.Copy(hash, fileObj); err != nil {
		return nil, err
	}
	fileHash := fmt.Sprintf("%x", hash.Sum(nil))

	// 检查文件是否已存在
	existingFile, err := file.FileRepository.FindByHashAndType(fileHash, req.FileType)
	if err == nil {
		// 文件已存在
		if existingFile.Status == models.FileStatusDisable {
			return nil, errors.New("file is disabled and cannot be uploaded")
		}

		if existingFile.Status == models.FileStatusNormal && existingFile.IsDeleted == models.DeleteNormal {
			return &dto.UploadResponse{ID: existingFile.ID, Chunked: false}, nil
		}

		// 更新状态为NORMAL
		existingFile.Status = models.FileStatusNormal
		existingFile.IsDeleted = models.DeleteNormal
		if err := file.FileRepository.Update(existingFile); err != nil {
			return nil, err
		}
		return &dto.UploadResponse{ID: existingFile.ID, Chunked: false}, nil
	}

	// 重新读取文件内容
	if _, err := fileObj.Seek(0, 0); err != nil {
		return nil, err
	}

	// 创建上传目录
	uploadPath := ctx.File.UploadDir
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, err
	}

	// 生成文件名和路径
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), fileHash, ext)
	relativePath := filepath.Join(req.FileType, fileName)
	absolutePath := filepath.Join(uploadPath, relativePath)

	// 保存文件
	out, err := os.Create(absolutePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, fileObj); err != nil {
		return nil, err
	}

	// 创建文件记录
	newFile := &models.SysFile{
		Name:         fileHeader.Filename,
		Type:         req.FileType,
		Extension:    strings.TrimPrefix(ext, "."),
		RelativePath: relativePath,
		AbsoluteURL:  fmt.Sprintf("%s/%s", ctx.File.BaseUrl, relativePath),
		Hash:         fileHash,
		Size:         fileHeader.Size,
		Status:       models.FileStatusNormal,
		IsDeleted:    models.DeleteNormal,
	}

	if err := file.FileRepository.Create(newFile); err != nil {
		return nil, err
	}

	return &dto.UploadResponse{ID: newFile.ID, Chunked: false}, nil
}

func (s *FileServiceImpl) GetFile(id int64) (*dto.FileResponse, error) {
	file, err := file.FileRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	response := dto.ToFileResponse(file)
	return &response, nil
}

func (s *FileServiceImpl) ListFiles(req dto.ListFilesRequest) (*dto.FileListResponse, error) {
	var status models.FileStatus
	var deleted models.DeletedStatus

	if req.Status != "" {
		status = models.FileStatus(req.Status)
	}

	if req.IsDeleted != "" {
		deleted = models.DeletedStatus(req.IsDeleted)
	}

	files, total, err := file.FileRepository.List(req.RelativePathPrefix, req.ID, status, deleted, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	response := &dto.FileListResponse{
		Total: total,
		Files: make([]dto.FileResponse, 0, len(files)),
	}

	for _, file := range files {
		response.Files = append(response.Files, dto.ToFileResponse(file))
	}

	return response, nil
}

func (s *FileServiceImpl) DownloadFile(id int64) (string, error) {
	file, err := file.FileRepository.FindById(id)
	if err != nil {
		return "", err
	}

	if file.Size > ctx.File.MinChunkSize {
		return "", errors.New("file too large, please use stream download")
	}

	return filepath.Join(ctx.File.UploadDir, file.RelativePath), nil
}

func (s *FileServiceImpl) StreamDownloadFile(id int64) (string, error) {
	file, err := file.FileRepository.FindById(id)
	if err != nil {
		return "", err
	}

	return filepath.Join(ctx.File.UploadDir, file.RelativePath), nil
}

func (s *FileServiceImpl) DeleteFile(id int64) error {
	return file.FileRepository.Delete(id)
}

func (s *FileServiceImpl) DisableFile(req dto.DisableFileRequest) error {
	return file.FileRepository.Disable(req.ID, req.Reason)
}

func (s *FileServiceImpl) UploadChunk(fileHeader *multipart.FileHeader, req dto.ChunkUploadRequest) (*dto.UploadResponse, error) {
	// 检查分片大小是否匹配
	if fileHeader.Size != req.ChunkSize {
		return nil, errors.New("chunk size mismatch")
	}

	// 生成或使用传入的uploadID
	uploadID := req.UploadID
	if uploadID == "" {
		uploadID = generateUploadID()
	}

	// 创建分片存储目录
	chunkDir := filepath.Join(ctx.File.UploadDir, "chunks", uploadID)
	if err := os.MkdirAll(chunkDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create chunk directory: %v", err)
	}

	// 保存分片文件
	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d.chunk", req.ChunkNumber))
	if err := saveUploadedFile(fileHeader, chunkPath); err != nil {
		return nil, fmt.Errorf("failed to save chunk: %v", err)
	}

	// 记录分片信息到数据库
	chunk := &models.SysChunkFile{
		UploadID:    uploadID,
		ChunkNumber: req.ChunkNumber,
		TotalChunks: req.TotalChunks,
		FilePath:    chunkPath,
		FileType:    req.FileType,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		ChunkSize:   req.ChunkSize,
		Hash:        req.Hash,
	}

	if err := file.FileRepository.CreateChunkUpload(chunk); err != nil {
		return nil, fmt.Errorf("failed to save chunk info: %v", err)
	}

	return &dto.UploadResponse{
		ID:       0,
		Chunked:  true,
		Message:  fmt.Sprintf("Chunk %d/%d uploaded successfully", req.ChunkNumber, req.TotalChunks),
		UploadID: uploadID,
	}, nil
}

func (s *FileServiceImpl) CompleteChunkUpload(req dto.CompleteChunkUploadRequest) (*dto.UploadResponse, error) {
	// 获取所有分片
	chunks, err := file.FileRepository.GetChunkUploads(req.UploadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks: %v", err)
	}

	// 检查是否所有分片都已上传
	if len(chunks) != req.TotalChunks {
		return nil, fmt.Errorf("incomplete upload: received %d of %d chunks", len(chunks), req.TotalChunks)
	}

	// 检查文件是否已存在
	existingFile, err := file.FileRepository.FindByHashAndType(req.Hash, req.FileType)
	if err == nil {
		// 文件已存在的处理逻辑（同普通上传）
		if existingFile.Status == models.FileStatusDisable {
			return nil, errors.New("file is disabled and cannot be uploaded")
		}

		if existingFile.Status == models.FileStatusNormal && existingFile.IsDeleted == models.DeleteNormal {
			return &dto.UploadResponse{ID: existingFile.ID, Chunked: false}, nil
		}

		existingFile.Status = models.FileStatusNormal
		existingFile.IsDeleted = models.DeleteNormal
		if err := file.FileRepository.Update(existingFile); err != nil {
			return nil, err
		}
		return &dto.UploadResponse{ID: existingFile.ID, Chunked: false}, nil
	}

	// 合并分片
	mergedFile, err := s.mergeChunks(chunks, req)
	if err != nil {
		return nil, fmt.Errorf("failed to merge chunks: %v", err)
	}

	// 标记分片上传完成
	if err := file.FileRepository.MarkChunkUploadComplete(req.UploadID); err != nil {
		return nil, fmt.Errorf("failed to mark upload as complete: %v", err)
	}

	return &dto.UploadResponse{
		ID:      mergedFile.ID,
		Chunked: true,
		Message: "File uploaded successfully",
	}, nil
}

func (s *FileServiceImpl) mergeChunks(chunks []*models.SysChunkFile, req dto.CompleteChunkUploadRequest) (*models.SysFile, error) {
	// 创建目标文件
	ext := filepath.Ext(req.FileName)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), req.Hash, ext)
	relativePath := filepath.Join(req.FileType, fileName)
	absolutePath := filepath.Join(ctx.File.UploadDir, relativePath)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(absolutePath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	outFile, err := os.Create(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create merged file: %v", err)
	}
	defer outFile.Close()

	// 按顺序合并所有分片
	for _, chunk := range chunks {
		chunkFile, err := os.Open(chunk.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open chunk %d: %v", chunk.ChunkNumber, err)
		}

		_, err = io.Copy(outFile, chunkFile)
		chunkFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to merge chunk %d: %v", chunk.ChunkNumber, err)
		}
	}

	// 验证文件大小
	fileInfo, err := outFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	if fileInfo.Size() != req.FileSize {
		return nil, fmt.Errorf("merged file size mismatch: expected %d, got %d", req.FileSize, fileInfo.Size())
	}

	// 创建文件记录
	newFile := &models.SysFile{
		Name:         req.FileName,
		Type:         req.FileType,
		Extension:    strings.TrimPrefix(ext, "."),
		RelativePath: relativePath,
		AbsoluteURL:  fmt.Sprintf("%s/%s", ctx.File.BaseUrl, relativePath),
		Hash:         req.Hash,
		Size:         req.FileSize,
		Status:       models.FileStatusNormal,
		IsDeleted:    models.DeleteNormal,
	}

	if err := file.FileRepository.Create(newFile); err != nil {
		return nil, fmt.Errorf("failed to create file record: %v", err)
	}

	// 清理分片文件（可选，可以在后台任务中执行）
	go func() {
		_ = os.RemoveAll(filepath.Join(ctx.File.UploadDir, "chunks", req.UploadID))
		_ = file.FileRepository.DeleteChunkUploads(req.UploadID)
	}()

	return newFile, nil
}

// 辅助函数
func generateUploadID() string {
	return uuid.New().String()
}

func saveUploadedFile(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
