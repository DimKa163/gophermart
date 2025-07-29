package middleware

import (
	"fmt"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		logging.Log.Info(
			"got incoming HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		startTime := time.Now()
		c.Next()
		elapsed := time.Since(startTime)
		userID, ok := c.Value("userId").(int64)
		if !ok {
			userID = -1
		}
		logging.Log.Info("Processed HTTP request", zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.Duration("elapsed", elapsed),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("userId", fmt.Sprintf("%d", userID)))
	}
}
