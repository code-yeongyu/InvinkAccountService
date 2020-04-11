package controllers

import (
	"invink/account-service/errors"
	"invink/account-service/forms"
	"invink/account-service/models"
	"invink/account-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RegisterUser godoc
// @Summary Register an user
// @Description Register an user with given information
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Param email path string true "Email"
// @Param password path string true "Password"
// @Param publicKey path string true "RSA 2048 PublicKey"
// @Param nickname path string false "Nickname"
// @Param bio path string false "Bio"
// @Success 201 {object} EmptyResponse "User account created"
// @Failure 400 {object} TypicalErrorResponse "Normal Form error, like username duplicate"
// @Router /register/ [post]
func (ctrler *Controller) RegisterUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	var inputForm forms.Registration

	if err := c.ShouldBindJSON(&inputForm); err != nil {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errors.FormErrorCode, err.Error())
		return
	}

	if errorCode := models.NewUser().ValidateUsername(db, inputForm.Username); errorCode != -1 {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errorCode, "")
		return
	}
	if errorCode := models.NewUser().ValidateEmail(db, inputForm.Email); errorCode != -1 {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errorCode, "")
		return
	}
	if errorCode := models.NewUser().ValidatePassword(inputForm.Password); errorCode != -1 {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errorCode, "")
		return
	}
	if !models.NewUser().IsPublicKeyInFormat(inputForm.PublicKey) {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errors.PublicKeyErrorCode, "")
		return
	}

	hashedPassword := models.NewUser().GenerateHashedPassword(inputForm.Password)
	user = models.User{
		Username:  inputForm.Username,
		Email:     inputForm.Email,
		Password:  string(hashedPassword),
		Nickname:  inputForm.Nickname,
		Bio:       inputForm.Bio,
		PublicKey: inputForm.PublicKey,
	}
	db.Create(&user)
	c.Data(http.StatusCreated, gin.MIMEHTML, nil)
}
