package controllers

import (
	"invink/account-service/errors"
	"invink/account-service/forms"
	"invink/account-service/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// AuthUser godoc
// @Summary Authenticate an user
// @Description Authenticate an user with given information, to get a jwt token
// @Accept  json
// @Produce  json
// @Param id path string true "Username or Email"
// @Param password path string true "Password"
// @Success 200 {object} AuthenticatedResponse "Valid information, authenticated"
// @Failure 400 {object} EmptyResponse "Wrong format or invalid information"
// @Router /auth/ [post]
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
		errorCode := errors.AuthenticationFailureCode
		c.JSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode]})
		return
	} // checking ID

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputForm.Password)) != nil {
		errorCode := errors.AuthenticationFailureCode
		c.JSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode]})
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &models.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("ACCOUNT_JWT_KEY")))

	if err != nil {
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
