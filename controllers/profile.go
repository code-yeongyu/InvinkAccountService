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

// GetMyProfile godoc
// @Summary Get my profile
// @Description Get a requested user's profile
// @Produce json
// @Success 200 {object} MyProfileResponse "When request to other's profile"
// @Failure 404
// @Router /profile/ [get]
func (ctrler *Controller) GetMyProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ID := c.MustGet("id").(uint64)
	var user models.User
	user.SetUserByID(db, ID)
	c.JSON(http.StatusOK, user.ToMyProfileMap())
}

// GetProfileByUsername godoc
// @Summary Get a profile by username
// @Description Get a profile by username with given information
// @Produce json
// @Success 200 {object} PublicProfileResponse "When request to other's profile"
// @Failure 404 {object} EmptyResponse "No such user"
// @Router /profile/:username [get]
func (ctrler *Controller) GetProfileByUsername(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ID := c.MustGet("id").(uint64)
	var loginUser, targetUser models.User

	loginUser.SetUserByID(db, ID)
	if err := targetUser.SetUserByUsername(db, c.Param("username")); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if c.Param("username") == loginUser.Username {
		ctrler.GetMyProfile(c)
		return
	}
	c.JSON(http.StatusOK, targetUser.ToPublicProfileMap())
}

// UpdateMyProfile godoc
// @Summary Update my profile
// @Description Update my profile with given information
// @Param username path string false "Username"
// @Param password path string false "Password"
// @Param nickname path string false "Nickname"
// @Param picture_url path string false "PictureURL"
// @Param bio path string false "Bio"
// @Param my_keys path string false "MyKeys"
// @Param current_password path string false "CurrentPassword"
// @Produce json
// @Success 200 {object} EmptyResponse "No errors occurred, profile was successfully updated"
// @Failure 400 {object} TypicalErrorResponse "Wrong format or invalid information"
// @Router /profile/ [patch]
func (ctrler *Controller) UpdateMyProfile(c *gin.Context) {
	var loginUser models.User
	var inputForm forms.Profile
	db := c.MustGet("db").(*gorm.DB)
	ID := c.MustGet("id").(uint64)

	c.ShouldBindJSON(&inputForm)

	loginUser.SetUserByID(db, ID)

	if (inputForm.Username != "" || inputForm.Password != "") &&
		!loginUser.IsPasswordCorrect(inputForm.CurrentPassword) {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errors.AuthenticationFailureCode, "If you want to change your username or password, please enter your current password correctly.")
		return
	}

	if inputForm.Username != "" {
		if errorCode := models.NewUser().ValidateUsername(db, inputForm.Username); errorCode != -1 {
			utils.AbortWithErrorResponse(c, http.StatusBadRequest, errorCode, "")
			return
		} // if username is not valid, then return with the reason error code
	}
	if inputForm.Password != "" {
		if errorCode := models.NewUser().ValidatePassword(inputForm.Password); errorCode != -1 {
			utils.AbortWithErrorResponse(c, http.StatusBadRequest, errorCode, "")
			return
		} // if username is not valid, then return with the reason error code
		inputForm.Password = string(models.NewUser().GenerateHashedPassword(inputForm.Password))
	}
	/*
		if inputForm.PictureURL != "" {
			// logics about checking if the image url is valid
		}
		if inputForm.MyKeys != "" {
			// logics about checking whether the user actually has the following keys from the following user
			// should be placed here
		}
	*/
	db.Model(models.User{}).Where(&loginUser).Updates(inputForm)
}

// ClearFieldFromMyProfile godoc
// @Summary Update my profile
// @Description Update my profile with given information. picture_url or bio will be only accepted for :field_name
// @Produce json
// @Success 200 {object} EmptyResponse "No errors occurred, profile was successfully removed"
// @Failure 400 {object} TypicalErrorResponse "Wrong password"
// @Router /profile/:field_name [delete]
func (ctrler *Controller) ClearFieldFromMyProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ID := c.MustGet("id").(uint64)
	fieldName := c.Param("field_name")
	if fieldName != "picture_url" && fieldName != "bio" && fieldName != "nickname" {
		errorCode := errors.ParameterErrorCode
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errorCode, "")
	}
	db.Model(models.User{}).Where("ID = ?", ID).Update(fieldName, gorm.Expr("NULL"))
}

// RemoveMyProfile godoc
// @Summary Update my profile
// @Description Update my profile with given information
// @Produce json
// @Success 200 {object} EmptyResponse "No errors occurred, profile was successfully removed"
// @Failure 400 {object} TypicalErrorResponse "Wrong password"
// @Router /profile/ [delete]
func (ctrler *Controller) RemoveMyProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ID := c.MustGet("id").(uint64)
	var loginUser models.User
	var inputForm forms.Profile
	loginUser.SetUserByID(db, ID)
	c.ShouldBindJSON(&inputForm)

	if !loginUser.IsPasswordCorrect(inputForm.CurrentPassword) {
		utils.AbortWithErrorResponse(c, http.StatusBadRequest, errors.AuthenticationFailureCode, "")
		return
	}

	db.Model(loginUser).Delete(loginUser)
}
