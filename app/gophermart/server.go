package gophermart

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/DimKa163/gophermart/internal/shared/tripper"
	"github.com/DimKa163/gophermart/internal/user/application"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/DimKa163/gophermart/internal/user/infrastructure/external/accrual"
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
	worker      *worker.OrderPooler
	crn         *cron.Cron
	accrualCl   accrual.AccrualClient
}
type Server struct {
	Config
	*gin.Engine
	*http.Server
	*ServiceContainer
}

func New(conf Config) *Server {
	router := gin.New()
	return &Server{
		Config: conf,
		Engine: router,
		Server: &http.Server{
			Addr:    conf.Addr,
			Handler: router.Handler(),
		},
		ServiceContainer: &ServiceContainer{},
	}
}

func (s *Server) AddServices() error {
	var err error
	s.pgPool, err = addPgPool(s.Database)
	if err != nil {
		return err
	}
	s.authService = s.addAuthService()
	s.unitOfWork = addUnitOfWork(s.pgPool)
	s.userAPI = rest.NewUserAPI(application.NewUserService(s.unitOfWork, s.authService),
		application.NewOrderService(s.unitOfWork))
	accrualCl := addAccrualClient(s.Accrual)
	s.crn = cron.New(cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	s.worker, err = worker.NewWorker(s.crn,
		s.Schedule, 10, application.NewTrackOrderHandler(s.unitOfWork, application.NewTrackOrderProcessor(accrualCl)))
	if err != nil {
		return err
	}
	return nil
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
	s.crn.Start()
	if err := s.worker.Run(ctx); err != nil {
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

func addPgPool(database string) (*pgxpool.Pool, error) {
	pg, err := pgxpool.New(context.Background(), database)
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func (s *Server) addAuthService() auth.AuthService {
	jwtAuth := auth.NewJWT(auth.JWTConfig{
		TokenExpiration: time.Minute * 30,
		SecretKey:       []byte(s.Secret),
	})
	s.authService = auth.NewAuthService(s.Argon, jwtAuth)
	return s.authService
}

func addUnitOfWork(db db.QueryExecutor) uow.UnitOfWork {
	return persistence.NewUnitOfWork(db)
}

func addAccrualClient(addr string) accrual.AccrualClient {
	tripperFc := []func(transport http.RoundTripper) http.RoundTripper{
		func(transport http.RoundTripper) http.RoundTripper {
			return tripper.NewRetryRoundTripper(transport)
		},
		func(transport http.RoundTripper) http.RoundTripper {
			return tripper.NewLoggingRoundTripper(transport)
		},
	}
	return accrual.New(addr, tripperFc)
}
