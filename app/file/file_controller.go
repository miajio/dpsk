package file

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/dto"
	"github.com/miajio/dpsk/pkg/router"
	"github.com/miajio/dpsk/service/file"
)

type FileControllerTemp interface {
	UploadFile(ctx *gin.Context)
	UploadChunk(ctx *gin.Context)
	GetFile(ctx *gin.Context)
	ListFiles(ctx *gin.Context)
	DownloadFile(ctx *gin.Context)
	StreamDownloadFile(ctx *gin.Context)
	DeleteFile(ctx *gin.Context)
	DisableFile(ctx *gin.Context)
}

type FileControllerImpl struct{}

var FileController router.Controller = &FileControllerImpl{}

func (c *FileControllerImpl) Option(e *gin.Engine) {
	fileGroup := e.Group("/api/files")
	{
		fileGroup.POST("/upload", c.UploadFile)
		fileGroup.GET("/:id", c.GetFile)
		fileGroup.GET("", c.ListFiles)
		fileGroup.GET("/download/:id", c.DownloadFile)
		fileGroup.GET("/stream/:id", c.StreamDownloadFile)
		fileGroup.DELETE("", c.DeleteFile)
		fileGroup.PUT("/disable", c.DisableFile)
		fileGroup.POST("/upload/chunk", c.UploadChunk)
		fileGroup.POST("/upload/complete", c.CompleteChunkUpload)
	}
}

func (c *FileControllerImpl) UploadFile(ctx *gin.Context) {
	var req dto.UploadFileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	response, err := file.FileService.UploadFile(fileHeader, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *FileControllerImpl) GetFile(ctx *gin.Context) {
	var req dto.GetFileRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := file.FileService.GetFile(req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, file)
}

func (c *FileControllerImpl) ListFiles(ctx *gin.Context) {
	var req dto.ListFilesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files, err := file.FileService.ListFiles(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, files)
}

func (c *FileControllerImpl) DownloadFile(ctx *gin.Context) {
	var req dto.DownloadFileRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filePath, err := file.FileService.DownloadFile(req.ID)
	if err != nil {
		if errors.Is(err, errors.New("file too large, please use stream download")) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.File(filePath)
}

func (c *FileControllerImpl) StreamDownloadFile(ctx *gin.Context) {
	var req dto.DownloadFileRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filePath, err := file.FileService.StreamDownloadFile(req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", string(fileInfo.Size()))

	ctx.File(filePath)
}

func (c *FileControllerImpl) DeleteFile(ctx *gin.Context) {
	var req dto.DeleteFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := file.FileService.DeleteFile(req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}

func (c *FileControllerImpl) DisableFile(ctx *gin.Context) {
	var req dto.DisableFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := file.FileService.DisableFile(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "file disabled successfully"})
}

func (c *FileControllerImpl) UploadChunk(ctx *gin.Context) {
	var req dto.ChunkUploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := ctx.FormFile("chunk")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "chunk file is required"})
		return
	}

	response, err := file.FileService.UploadChunk(fileHeader, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *FileControllerImpl) CompleteChunkUpload(ctx *gin.Context) {
	var req dto.CompleteChunkUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := file.FileService.CompleteChunkUpload(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
