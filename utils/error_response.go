package utils

import (
	"invink/account-service/errors"

	"github.com/gin-gonic/gin"
)

// AbortWithErrorResponse aborts the request with the given error
func AbortWithErrorResponse(c *gin.Context, statusCode int, errorCode int, detail string) {
	c.AbortWithStatusJSON(statusCode, gin.H{"error": errorCode, "msg": errors.Messages[errorCode], "detail": detail})
}
