package controllers

import (
	"invink/account-service/errors"
	"invink/account-service/forms"
	"invink/account-service/models"
	"invink/account-service/utils"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AuthUser godoc
// @Summary Authenticate an user
// @Description Authenticate an user with given information, to get a jwt token
// @Accept  json
// @Produce  json
// @Param id path string true "Username or Email"
// @Param password path string true "Password"
// @Success 200 {object} AuthenticatedResponse "Valid information, authenticated"
// @Failure 400 {object} TypicalErrorResponse "Wrong format or invalid information"
// @Router /auth/ [post]
func (ctrler *Controller) AuthUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	var inputForm forms.Authentication

	if err := c.ShouldBindJSON(&inputForm); err != nil {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errors.FormErrorCode, err.Error())
		return
	}

	if err := user.SetUserByEmailOrID(db, inputForm.ID); err != nil {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errors.AuthenticationFailureCode, "")
		return
	} // no such user
	if !user.IsPasswordCorrect(inputForm.Password) {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errors.AuthenticationFailureCode, "")
		return
	}

	claims := models.IssueToken(user.ID, time.Now().Add(15*time.Minute))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("ACCOUNT_JWT_KEY")))
	if err != nil {
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
