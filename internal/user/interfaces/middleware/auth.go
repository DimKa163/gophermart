package middleware

import (
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth(jwt *auth.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		tk, cl, err := jwt.ParseJWT(tokenValue)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !tk.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("userId", cl.UserID)
		c.Next()
	}
}
