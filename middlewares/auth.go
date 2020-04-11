package middlewares

import (
	"invink/account-service/errors"
	"invink/account-service/models"
	"invink/account-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthenticateJWT is a middlware for the endpoints that requires authentication
func AuthenticateJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	claims := models.Claims{}

	if authHeader == "" {
		errorCode := errors.EmptyAuthorizationHeaderCode
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorCode, "msg": errors.Messages[errorCode]})
		c.Abort()
		return
	}
	if authHeader[:7] != "Bearer " {
		errorCode := errors.WrongTokenTypeCode
		utils.AbortWithErrorResponse(c, http.StatusUnauthorized, errorCode, "")
		return
	}
	// validating the header Authorization

	if err := claims.ParseIfValidJWT(authHeader[7:]); err != nil {
		utils.AbortWithErrorResponse(c, http.StatusUnauthorized, errors.AuthenticationFailureCode, "")
	}

	c.Set("id", claims.ID)
	c.Next()
}
