package middleware

import (
	"fmt"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Auth(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		cl, err := authService.Verify(tokenValue)
		logger := logging.Logger(c)
		if err != nil {
			logger.Warn("Authorization Error", zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		auth.SetUser(c, cl.UserID)
		logging.SetLogger(c, logger.With(zap.String("userId", fmt.Sprintf("%d", cl.UserID))))
		c.Next()
	}
}
