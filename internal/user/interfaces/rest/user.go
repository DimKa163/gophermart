package rest

import (
	"errors"
	"github.com/DimKa163/gophermart/internal/user/application/login"
	"github.com/DimKa163/gophermart/internal/user/application/register"
	"github.com/DimKa163/gophermart/internal/user/interfaces/contracts"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserApi interface {
	Register(context *gin.Context)
	Login(context *gin.Context)
	AddOrder(context *gin.Context)
	GetOrders(context *gin.Context)
	GetBalance(context *gin.Context)
	Withdraw(context *gin.Context)
	GetWithdrawals(context *gin.Context)
}

type userApi struct {
	registerHandler *register.RegisterHandler
	loginHandler    *login.LoginHandler
}

func NewUserApi(registerHandler *register.RegisterHandler, loginHandler *login.LoginHandler) UserApi {
	return &userApi{registerHandler: registerHandler, loginHandler: loginHandler}
}

func (u *userApi) Register(context *gin.Context) {
	var user contracts.User
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := u.registerHandler.Handle(context, &register.RegisterCommand{Login: user.Login, Password: user.Password})
	if err != nil {
		if errors.Is(err, register.ErrLoginAlreadyExists) {
			context.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"id": id})
}

func (u *userApi) Login(context *gin.Context) {
	var user contracts.User
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := u.loginHandler.Handle(context, &login.LoginCommand{Login: user.Login})
	if err != nil {
		if errors.Is(err, login.ErrInvalidPassword) {
			context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, login.ErrUserNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": token})
}

func (u *userApi) AddOrder(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *userApi) GetOrders(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *userApi) GetBalance(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *userApi) Withdraw(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *userApi) GetWithdrawals(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}
