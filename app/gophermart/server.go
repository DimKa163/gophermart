package gophermart

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/DimKa163/gophermart/internal/user/application/balance"
	"github.com/DimKa163/gophermart/internal/user/application/login"
	"github.com/DimKa163/gophermart/internal/user/application/order"
	"github.com/DimKa163/gophermart/internal/user/application/register"
	"github.com/DimKa163/gophermart/internal/user/application/withdrawal"
	"github.com/DimKa163/gophermart/internal/user/infrastructure/persistence"
	"github.com/DimKa163/gophermart/internal/user/interfaces/middleware"
	"github.com/DimKa163/gophermart/internal/user/interfaces/rest"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ServiceContainer struct {
	userApi    rest.UserApi
	jwtBuilder *auth.JWT
	pgPool     *pgxpool.Pool
}
type Server struct {
	Config
	*gin.Engine
	*http.Server
	*ServiceContainer
}

func New(conf Config) (*Server, error) {
	pg, err := pgxpool.New(context.Background(), conf.Database)
	if err != nil {
		return nil, err
	}
	uow := persistence.NewUnitOfWork(pg)
	jwtBuilder := auth.NewJWTBuilder(auth.JWTConfig{
		TokenExpiration: time.Minute * 5,
		SecretKey:       []byte(conf.Secret),
	})
	registerHandler := register.New(uow, jwtBuilder)
	loginHandler := login.New(uow, jwtBuilder)
	uploadOrderHandler := order.NewUploadOrderHandler(uow)
	orderQueryHandler := order.NewOrderQueryHandler(uow)
	balanceHandler := balance.NewBalanceQueryHandler(uow)
	withdrawHandler := balance.NewWithdrawHandler(uow)
	withdrawalQueryHandler := withdrawal.NewWithdrawalQueryHandler(uow)
	router := gin.New()

	return &Server{
		Config: conf,
		Engine: router,
		Server: &http.Server{
			Addr:    conf.Addr,
			Handler: router.Handler(),
		},
		ServiceContainer: &ServiceContainer{
			userApi: rest.NewUserApi(registerHandler,
				loginHandler,
				uploadOrderHandler,
				orderQueryHandler,
				balanceHandler,
				withdrawHandler,
				withdrawalQueryHandler),
			jwtBuilder: jwtBuilder,
			pgPool:     pg,
		},
	}, nil
}

func (s *Server) AddLogging() error {
	return logging.Initialize(s.LogLevel)
}

func (s *Server) Map() {
	s.Use(gin.Recovery())
	s.Use(middleware.Logging())
	userGroup := s.Group("api/user")
	{
		userApi := s.userApi
		userGroup.POST("/register", userApi.Register)
		userGroup.POST("/login", userApi.Login)
		userGroup.Use(middleware.Auth(s.jwtBuilder))
		{
			userGroup.GET("/orders", userApi.GetOrders)
			userGroup.GET("/withdrawals", userApi.GetWithdrawals)
			userGroup.POST("/orders", userApi.Upload)
			balanceGroup := userGroup.Group("/balance")
			{
				balanceGroup.GET("", userApi.GetBalance)
				balanceGroup.POST("/withdraw", userApi.Withdraw)
			}
		}
	}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := persistence.Migrate(s.pgPool); err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = s.Server.Shutdown(timeoutCtx)
	}()
	return s.ListenAndServe()
}
