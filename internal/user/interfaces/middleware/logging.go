package middleware

import (
	"bytes"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := c.GetRawData()
		logging.Log.Info(
			"got incoming HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("body", string(data)),
		)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		logger := logging.Log.With(zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path))
		logging.SetLogger(c, logger)
		startTime := time.Now()
		c.Next()
		elapsed := time.Since(startTime)
		log := logging.Logger(c)
		logger = log.With(zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.Duration("elapsed", elapsed),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path))
		logger.Info("Processed HTTP request")
	}
}
