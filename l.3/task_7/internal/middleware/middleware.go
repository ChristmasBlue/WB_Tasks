package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	admin   string = "admin"
	manager string = "manager"
	viewer  string = "viewer"
)

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok || userRole == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role missing in token"})
			return
		}

		if userRole != admin && userRole != manager && userRole != viewer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user role mate"})
			return
		}

		requestMethod := c.Request.Method
		switch {
		case requestMethod == http.MethodPost || requestMethod == http.MethodDelete:
			if userRole != admin {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not allowed"})
				return
			}
		case requestMethod == http.MethodPut:
			if userRole != admin && userRole != manager {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not allowed"})
				return
			}
		}

		c.Next()
	}
}
