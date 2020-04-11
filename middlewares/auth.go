package middlewares

import (
	"invink/account-service/errors"
	"invink/account-service/models"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AuthenticateJWT is a middlware for the endpoints that requires authentication
func AuthenticateJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		errorCode := errors.EmptyAuthorizationHeaderCode
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorCode, "msg": errors.Messages[errorCode]})
		c.Abort()
		return
	}

	if authHeader[:7] != "Bearer " {
		errorCode := errors.WrongTokenTypeCode
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorCode, "msg": errors.Messages[errorCode]})
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
		errorCode := errors.AuthenticationFailureCode
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorCode, "msg": errors.Messages[errorCode], "detail": err.Error()})
		c.Abort()
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Model(models.User{}).Where("ID = ?", claims.ID).First(&models.User{}).Error; err != nil {
		errorCode := errors.AuthenticationFailureCode
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorCode, "msg": errors.Messages[errorCode], "detail": "Username has changed, therefore you have to refresh your token"})
		c.Abort()
		return
	}

	c.Set("id", claims.ID)

	c.Next()
}
