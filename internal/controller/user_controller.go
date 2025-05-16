package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miajio/dpsk/internal/dto"
	"github.com/miajio/dpsk/internal/service"
)

type UserController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(service service.UserService) UserController {
	return &userController{service: service}
}

func (c *userController) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewBaseResponse(http.StatusBadRequest, "参数错误", nil, err))
		return
	}
	userModel, err := c.service.Login(loginRequest.Account, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "登录失败", nil, err))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "登录成功", dto.LoginResponse{
		Account:  userModel.Account,
		Nickname: userModel.Nickname,
		Token:    "",
	}, nil))
}

func (c *userController) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusBadRequest, "参数错误", nil, err))
		return
	}
	if err := c.service.Register(registerRequest.Account, registerRequest.Nickname, registerRequest.Password); err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "注册失败", nil, err))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "注册成功", nil, nil))
}
