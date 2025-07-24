package rest

import (
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"github.com/gin-gonic/gin"
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
	repository.Unit
}

func NewUserApi() UserApi {
	return &userApi{}
}

func (u userApi) Register(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u userApi) Login(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u userApi) AddOrder(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u userApi) GetOrders(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u userApi) GetBalance(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u userApi) Withdraw(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (u userApi) GetWithdrawals(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}
