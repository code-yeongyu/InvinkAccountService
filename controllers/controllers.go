package controllers

import (
	"invink/account-service/errors"
	"invink/account-service/models"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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

func verifyUsername(s string) bool {
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
		case unicode.IsLower(c):
		case unicode.IsUpper(c):
		case c == '.' || c == '-' || c == '_':
		default:
			return false
		}
	}
	return true
}

func validateUsername(db *gorm.DB, username string) (errorCode int) {
	if err := db.Where("username = ?", username).First(&models.User{}).Error; err == nil {
		errorCode = errors.UsernameExistsCode
		return
	} // validating username duplicates
	if !verifyUsername(username) {
		errorCode = errors.UsernameFormatErrorCode
		return
	} // validating if username is in proper format
	return -1
}

func validatePassword(password string) (errorCode int) {
	if len(password) < 8 {
		errorCode = errors.PasswordTooShortCode
		return
	} // validating password length

	if !verifyPassword(password) {
		errorCode = errors.PasswordVulnerableErrorCode
		return
	}
	return -1
}

func isPasswordCorrect(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
