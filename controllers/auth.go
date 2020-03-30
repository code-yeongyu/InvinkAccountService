package controllers

import (
	"invink/account-service/errors"
	"invink/account-service/forms"
	"invink/account-service/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// AuthUser godoc
func (ctrler *Controller) AuthUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	var inputForm forms.Authentication

	if err := c.ShouldBindJSON(&inputForm); err != nil {
		errorCode := errors.FormErrorCode
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode], "detail": err.Error()})
		return
	}

	if err := db.Where("email = ? OR username = ?", inputForm.ID, inputForm.ID).First(&user).Error; err != nil {
		abortWith400ErrorResponse(c, errors.EmailExistsCode)
		return
	} // checking ID

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputForm.Password)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to authenticate."})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("jwtkey"))

	if err != nil {
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
