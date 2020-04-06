package middlewares

import (
	"net/http"
	"os"

	"invink/account-service/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthenticateJWT is a middlware for the endpoints that requires authentication
func AuthenticateJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Only Bearer token is available"})
		c.Abort()
		return
	}

	tokenString := authHeader[7:]
	claims := &models.Claims{}

	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCOUNT_JWT_KEY")), nil
		},
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": err.Error()})
		c.Abort()
		return
	}

	c.Set("username", claims.Username)

	c.Next()
}
