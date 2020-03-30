package controllers

import (
	"invink/account-service/errors"
	"invink/account-service/forms"
	"invink/account-service/models"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func abortWith404ErrorResponse(c *gin.Context, errorCode int) {
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

// RegisterUser is a handler for registering a user
func RegisterUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	var inputForm forms.Registration

	if err := c.ShouldBindJSON(&inputForm); err != nil {
		errorCode := errors.FormErrorCode
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode], "detail": err.Error()})
		return
	}

	if err := db.Where("email = ?", inputForm.Email).First(&user).Error; err == nil {
		abortWith404ErrorResponse(c, errors.EmailExistsCode)
		return
	} // validating email duplicates

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`) // email format
	if !re.MatchString(inputForm.Email) {
		abortWith404ErrorResponse(c, errors.EmailFormatErrorCode)
		return
	} // validating if email is in proper format

	/* email validation */

	if err := db.Where("username = ?", inputForm.Username).First(&user).Error; err == nil {
		abortWith404ErrorResponse(c, errors.UsernameExistsCode)
		return
	} // validating username duplicates

	re = regexp.MustCompile(`^[0-9a-zA-Z._]+$`)
	if !re.MatchString(inputForm.Username) {
		abortWith404ErrorResponse(c, errors.UsernameFormatErrorCode)
		return
	} // validating if username is in proper format

	/* username validation */

	if len(inputForm.Password) < 8 {
		abortWith404ErrorResponse(c, errors.PasswordTooShortCode)
		return
	} // validating password length

	if !verifyPassword(inputForm.Password) {
		abortWith404ErrorResponse(c, errors.PasswordVulnerableErrorCode)
		return
	} // validating password format

	/* password validation */

	if !strings.Contains(inputForm.PublicKey, "PUBLIC KEY") {
		abortWith404ErrorResponse(c, errors.PublicKeyErrorCode)
		return
	}

	/* publickey validation */

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(inputForm.Password), 15)

	user = models.User{
		Username:  inputForm.Username,
		Email:     inputForm.Email,
		Password:  string(passwordHash),
		Nickname:  inputForm.Nickname,
		Bio:       inputForm.Bio,
		PublicKey: inputForm.PublicKey,
	}

	db.Create(&user)

	c.Data(http.StatusCreated, gin.MIMEHTML, nil)
}
