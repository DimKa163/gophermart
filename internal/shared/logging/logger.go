package logging

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Log = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

func Logger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if !ok {
		logger = Log
	}
	return logger
}

func SetLogger(ctx context.Context, l *zap.Logger) context.Context {
	//чёртов gin.Context
	gCtx, ok := ctx.(*gin.Context)
	if !ok {
		ctx = context.WithValue(ctx, "logger", l)
		return ctx
	}
	gCtx.Set("logger", l)
	return gCtx
}
