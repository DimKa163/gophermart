package gophermart

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/user/application/login"
	"github.com/DimKa163/gophermart/internal/user/application/register"
	"github.com/DimKa163/gophermart/internal/user/infrastructure/persistence"
	"github.com/DimKa163/gophermart/internal/user/interfaces/rest"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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
	pg, err := pgxpool.New(context.Background(), conf.Database)
	if err != nil {
		return nil, err
	}
	if err = persistence.Migrate(pg); err != nil {
		return nil, err
	}
	router := gin.New()
	router.Use(gin.Recovery())
	uow := persistence.NewUnitOfWork(pg)
	jwtBuilder := auth.NewJWTBuilder(auth.JWTBuilderConfig{
		TokenExpiration: time.Minute * 5,
		SecretKey:       []byte(conf.Secret),
	})
	registerHandler := register.New(uow)
	loginHandler := login.New(uow, jwtBuilder)
	return &Server{
		Engine: router,
		Server: &http.Server{
			Addr:    conf.Addr,
			Handler: router.Handler(),
		},
		ServiceContainer: &ServiceContainer{
			userApi: rest.NewUserApi(registerHandler,
				loginHandler),
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
