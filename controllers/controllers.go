package controllers

import (
	"invink/account-service/errors"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
)

// Controller is a wrapper for all the controllers
type Controller struct {
}

// NewController returns a new Controller instance
func NewController() *Controller {
	return &Controller{}
}

func abortWith400ErrorResponse(c *gin.Context, errorCode int) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode]})
}

func verifyPassword(s string) bool {
	var number, lower, upper, special bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}
	return number && lower && upper && special
}
