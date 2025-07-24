package gophermart

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/interfaces/rest"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ServiceContainer struct {
	userApi rest.UserApi
}
type Server struct {
	*gin.Engine
	*http.Server
	*ServiceContainer
}

func New(conf Config) (*Server, error) {
	router := gin.New()
	router.Use(gin.Recovery())
	return &Server{
		Engine: router,
		ServiceContainer: &ServiceContainer{
			userApi: rest.NewUserApi(),
		},
	}, nil
}

func (s *Server) Map() {
	userGroup := s.Group("api/user")
	{
		userApi := s.userApi
		userGroup.GET("/orders", userApi.GetOrders)
		userGroup.GET("/withdrawals", userApi.GetWithdrawals)
		userGroup.POST("/register", userApi.Register)
		userGroup.POST("/login", userApi.Login)
		userGroup.POST("/orders", userApi.AddOrder)
		balanceGroup := userGroup.Group("/balance")
		{
			balanceGroup.GET("/", userApi.GetBalance)
			balanceGroup.POST("/withdraw", userApi.Withdraw)
		}

	}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = s.Server.Shutdown(timeoutCtx)
	}()
	return s.ListenAndServe()
}
