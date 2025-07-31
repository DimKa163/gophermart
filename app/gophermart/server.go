package gophermart

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/DimKa163/gophermart/internal/shared/mediatr"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/application/balance"
	"github.com/DimKa163/gophermart/internal/user/application/login"
	"github.com/DimKa163/gophermart/internal/user/application/order"
	"github.com/DimKa163/gophermart/internal/user/application/register"
	"github.com/DimKa163/gophermart/internal/user/application/withdrawal"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/DimKa163/gophermart/internal/user/infrastructure/persistence"
	"github.com/DimKa163/gophermart/internal/user/interfaces/middleware"
	"github.com/DimKa163/gophermart/internal/user/interfaces/rest"
	"github.com/DimKa163/gophermart/internal/user/interfaces/worker"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ServiceContainer struct {
	userAPI     rest.UserAPI
	authService auth.AuthService
	unitOfWork  uow.UnitOfWork
	pgPool      *pgxpool.Pool
	worker      *worker.Worker
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
	unitOfWork := persistence.NewUnitOfWork(pg)
	jwtAuth := auth.NewJWT(auth.JWTConfig{
		TokenExpiration: time.Minute * 30,
		SecretKey:       []byte(conf.Secret),
	})
	authService := auth.NewAuthService(conf.Argon, jwtAuth)

	wrk := worker.NewWorker(cron.New(cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger))),
		conf.Schedule)
	router := gin.New()
	return &Server{
		Config: conf,
		Engine: router,
		Server: &http.Server{
			Addr:    conf.Addr,
			Handler: router.Handler(),
		},
		ServiceContainer: &ServiceContainer{
			userAPI:     rest.NewUserAPI(),
			authService: authService,
			pgPool:      pg,
			worker:      wrk,
			unitOfWork:  unitOfWork,
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
		userAPI := s.userAPI
		userGroup.POST("/register", userAPI.Register)
		userGroup.POST("/login", userAPI.Login)
		userGroup.Use(middleware.Auth(s.authService))
		{
			userGroup.GET("/orders", userAPI.GetOrders)
			userGroup.GET("/withdrawals", userAPI.GetWithdrawals)
			userGroup.POST("/orders", userAPI.Upload)
			balanceGroup := userGroup.Group("/balance")
			{
				balanceGroup.GET("", userAPI.GetBalance)
				balanceGroup.POST("/withdraw", userAPI.Withdraw)
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
	if err := s.worker.Start(ctx); err != nil {
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

func (s *Server) AddMediatr() error {
	var err error
	err = mediatr.Bind[*register.RegisterCommand,
		*types.AppResult[string]](register.New(s.unitOfWork, s.authService))
	if err != nil {
		return err
	}
	err = mediatr.Bind[*login.LoginCommand,
		*types.AppResult[string]](login.New(s.unitOfWork, s.authService))
	if err != nil {
		return err
	}
	err = mediatr.Bind[*order.UploadOrderCommand,
		*types.AppResult[any]](order.NewUploadOrderHandler(s.unitOfWork))
	if err != nil {
		return err
	}
	err = mediatr.Bind[*order.OrderQuery,
		*types.AppResult[[]*model.Order]](order.NewOrderQueryHandler(s.unitOfWork))
	if err != nil {
		return err
	}
	err = mediatr.Bind[*balance.BalanceQuery,
		*types.AppResult[*model.BonusBalance]](balance.NewBalanceQueryHandler(s.unitOfWork))
	if err != nil {
		return err
	}
	err = mediatr.Bind[*balance.WithdrawCommand,
		*types.AppResult[any]](balance.NewWithdrawHandler(s.unitOfWork))
	if err != nil {
		return err
	}
	err = mediatr.Bind[*withdrawal.WithdrawalQuery,
		*types.AppResult[[]*model.BonusMovement]](withdrawal.NewWithdrawalQueryHandler(s.unitOfWork))
	if err != nil {
		return err
	}
	err = mediatr.Bind[*order.TrackOrderCommand,
		*types.AppResult[any]](order.NewTrackOrderHandler(s.unitOfWork))
	if err != nil {
		return err
	}
	return nil
}
