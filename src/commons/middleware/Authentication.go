package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenValidator interface {
	ValidateToken(token string) (string, error)
}

func Authentication(tokenValidator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "authorization header tidak ditemukan",
				},
			)
			return
		}

		parts := strings.SplitN(authorization, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "format authorization tidak valid",
				},
			)
			return
		}

		token := parts[1]

		userID, err := tokenValidator.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "token tidak valid",
				},
			)
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
