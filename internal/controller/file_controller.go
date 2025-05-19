package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/miajio/dpsk/internal/cache"
	"github.com/miajio/dpsk/internal/dto"
	"github.com/miajio/dpsk/internal/errors"
	"github.com/miajio/dpsk/internal/model"
	"github.com/miajio/dpsk/internal/service"
)

type FileController interface {
	Upload(ctx *gin.Context)
}

type fileController struct {
	service service.FileService
}

// NewFileController 创建文件控制器
func NewFileController(service service.FileService) FileController {
	return &fileController{service: service}
}

// Upload 上传文件
func (c *fileController) Upload(ctx *gin.Context) {
	var fileUpload dto.FileUploadRequest
	if err := ctx.ShouldBindWith(&fileUpload, binding.FormMultipart); err != nil {
		ctx.JSON(400, dto.NewBaseResponse(400, "参数错误", nil, err))
		return
	}

	var userId int64
	// 如果业务类型不是用户头像，则需要验证登录用户
	if fileUpload.BusinessType != "USER_HEADER" {
		// 获取登录用户信息
		loginUser, err := service.GetLoginUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dto.NewBaseResponse(http.StatusUnauthorized, "未授权", nil, err))
			return
		}
		userId = loginUser.Id
	}

	hash, err := fileUpload.GetHash()
	if err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusBadRequest, "文件上传失败", nil, err))
		return
	}

	// 判断文件是否存在
	exists, err := c.service.IsFileExists(hash)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "文件上传失败", nil, err))
		return
	} else if exists {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "文件上传失败", nil, errors.ErrFileExists))
		return
	}

	file := fileUpload.File
	suffix := filepath.Ext(file.Filename)
	if suffix == "" {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusBadRequest, "文件上传失败", nil, errors.ErrFileSuffixNotFound))
		return
	}
	fileName := uuid.New().String() + suffix

	savePath := filepath.Join(cache.File.SavePath, fileName)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "文件上传失败", nil, err))
		return
	}

	fileModel := &model.FileModel{
		Name:         filepath.Base(file.Filename),                 // 文件名称
		Path:         filepath.Join(cache.File.FilePath, fileName), // 文件存储路径
		Url:          cache.File.Url + fileName,                    // 文件访问地址
		Size:         file.Size,                                    // 文件大小
		Hash:         hash,                                         // 文件哈希值
		Extension:    suffix,                                       // 文件后缀
		BusinessType: fileUpload.BusinessType,                      // 业务类型
		Model: model.Model{
			CreateUserId: userId,
		},
	}
	if err := c.service.CreateFile(fileModel); err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "文件上传失败", nil, err))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "文件上传成功", fileModel, nil))
}
