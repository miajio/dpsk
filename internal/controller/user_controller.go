package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miajio/dpsk/internal/cache"
	"github.com/miajio/dpsk/internal/dto"
	"github.com/miajio/dpsk/internal/service"
	"go.uber.org/zap"
)

type UserController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Logout(ctx *gin.Context)
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

	loginUser := dto.LoginUser{
		Id:       userModel.Id,
		Nickname: userModel.Nickname,
		Account:  userModel.Account,
	}

	token := uuid.NewString()
	zap.S().Infof("执行登录操作,将输入写入redis, key: %s, value: %s", cache.JWT.Prefix+token, loginUser.Marshal())
	if err := cache.RedisClient.SetEx(ctx, cache.JWT.Prefix+token, loginUser.Marshal(), cache.JWT.Expires); err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "登录失败", nil, err))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "登录成功", dto.LoginResponse{
		LoginUser: dto.LoginUser{
			Id:       int64(userModel.Id),
			Account:  userModel.Account,
			Nickname: userModel.Nickname,
		},
		Token: token,
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

func (c *userController) Logout(ctx *gin.Context) {
	if err := cache.RedisClient.Del(ctx, cache.JWT.Prefix+ctx.GetHeader("Authorization")); err != nil {
		ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusInternalServerError, "退出失败", nil, err))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewBaseResponse(http.StatusOK, "退出成功", nil, nil))
}
