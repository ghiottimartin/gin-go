package helpers

import (
	"net/http"
	"os"
	"strings"

	"github.com/cookit/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldCheckToken(c.Request.RequestURI) {
			return
		}
		tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set the token and its claims in the context for later use
		c.Set("token", token)
		c.Set("claims", token.Claims.(*models.AppClaims))
		c.Next()
	}
}
