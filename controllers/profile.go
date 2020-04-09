package controllers

import (
	"invink/account-service/errors"
	"invink/account-service/forms"
	"invink/account-service/models"
	"strings"
	"unicode"

	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func getProfile(db *gorm.DB, username string) (user models.User, err error) {
	err = db.Where("username = ?", username).First(&user).Error
	return
}

func modelToMyProfileMap(profile models.User) (myProfileMap map[string]interface{}) {
	myProfileMap = structs.Map(profile)

	myProfileMap["my_keys"] = myProfileMap["MyKeys"]
	delete(myProfileMap, "MyKeys")
	myProfileMap["picture_url"] = myProfileMap["PictureURL"]
	delete(myProfileMap, "PictureURL")
	myProfileMap["public_key"] = myProfileMap["PublicKey"]
	delete(myProfileMap, "PublicKey")

	delete(myProfileMap, "ID")
	delete(myProfileMap, "Password")
	delete(myProfileMap, "Follower")
	delete(myProfileMap, "Following")

	for k, v := range myProfileMap {
		if unicode.IsUpper(rune(k[0])) {
			if v != "" {
				myProfileMap[strings.ToLower(k)] = v
			}
			delete(myProfileMap, k)
		} else {
			if v == "" {
				delete(myProfileMap, k)
			}
		}
	}

	followingUsername := make([]string, len(profile.Following))
	followerUsername := make([]string, len(profile.Follower))

	for i, v := range profile.Following {
		followingUsername[i] = v.Username
	}
	for i, v := range profile.Follower {
		followerUsername[i] = v.Username
	}
	myProfileMap["following_username"] = followingUsername
	myProfileMap["follower_username"] = followerUsername

	return
}

func modelToPublicProfileMap(profile models.User) (publicProfileMap map[string]interface{}) {
	publicProfileMap = map[string]interface{}{
		"username":      profile.Username,
		"following_cnt": len(profile.Following),
		"follower_cnt":  len(profile.Follower),
	}
	if profile.Nickname != "" {
		publicProfileMap["nickname"] = profile.Nickname
	}
	if profile.Bio != "" {
		publicProfileMap["bio"] = profile.Bio
	}
	if profile.PictureURL != "" {
		publicProfileMap["picture_url"] = profile.PictureURL
	}
	return
}

// GetMyProfile godoc
// @Summary Get my profile
// @Description Get a requested user's profile
// @Produce json
// @Success 200 {object} MyProfileResponse "When request to other's profile"
// @Failure 404
// @Router /profile/ [get]
func (ctrler *Controller) GetMyProfile(c *gin.Context) {
	username := c.MustGet("username").(string)

	db := c.MustGet("db").(*gorm.DB)
	profile, _ := getProfile(db, username)

	c.JSON(http.StatusOK, modelToMyProfileMap(profile))
}

// GetProfileByUsername godoc
// @Summary Get a profile by username
// @Description Get a profile by username with given information
// @Produce json
// @Success 200 {object} PublicProfileResponse "When request to other's profile"
// @Failure 404 {object} EmptyResponse "No such user"
// @Router /profile/:username [get]
func (ctrler *Controller) GetProfileByUsername(c *gin.Context) {
	username := c.Param("username")
	db := c.MustGet("db").(*gorm.DB)

	if c.Param("username") == c.MustGet("username").(string) {
		ctrler.GetMyProfile(c)
		return
	}

	profile, err := getProfile(db, username)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, modelToPublicProfileMap(profile))
}

// UpdateMyProfile godoc
// @Summary Update my profile
// @Description Update my profile with given information
// @Produce json
// @Success 200 {object} EmptyResponse "No errors occurred, profile was successfully updated"
// @Failure 400 {object} TypicalErrorResponse "Wrong format or invalid information"
// @Router /profile/ [patch]
func (ctrler *Controller) UpdateMyProfile(c *gin.Context) {
	var profile models.User
	var inputForm forms.Profile
	username := c.MustGet("username").(string)
	db := c.MustGet("db").(*gorm.DB)

	if err := c.ShouldBindJSON(&inputForm); err != nil {
		errorCode := errors.FormErrorCode
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode], "detail": err.Error()})
		return
	}

	profile, _ = getProfile(db, username)

	if bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(inputForm.CurrentPassword)) != nil {
		errorCode := errors.AuthenticationFailureCode
		c.JSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode]})
		return
	}
	// checking current_password

	if inputForm.Username != "" {
		if errorCode := validateUsername(db, inputForm.Username); errorCode != -1 {
			abortWith400ErrorResponse(c, errorCode)
			return
		}
		profile.Username = inputForm.Username
	}
	if inputForm.Password != "" {
		if errorCode := validatePassword(inputForm.Password); errorCode != -1 {
			abortWith400ErrorResponse(c, errorCode)
			return
		}
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(inputForm.Password), 15)
		profile.Password = string(passwordHash)
	}
	if inputForm.Nickname != "" {
		profile.Nickname = inputForm.Nickname
	}
	if inputForm.PictureURL != "" {
		profile.PictureURL = inputForm.PictureURL
	}
	if inputForm.Bio != "" {
		profile.Bio = inputForm.Bio
	}
	if inputForm.MyKeys != "" {
		// logics about checking whether the user actually has the following keys from the following user
		// should be placed here
		profile.MyKeys = inputForm.MyKeys
	}

	db.Save(&profile)

	c.Data(http.StatusOK, gin.MIMEHTML, nil)
}

// RemoveMyProfile godoc
// @Summary Update my profile
// @Description Update my profile with given information
// @Produce json
// @Success 200 {object} EmptyResponse "No errors occurred, profile was successfully removed"
// @Failure 400 {object} TypicalErrorResponse "Wrong password"
<<<<<<< HEAD
// @Router /profile/ [delete]
=======
// @Router /profile/ [patch]
>>>>>>> 3f32be481b67952a0423cf29d441cd5f50adb780
func (ctrler *Controller) RemoveMyProfile(c *gin.Context) {
	var inputForm forms.Profile
	username := c.MustGet("username").(string)
	db := c.MustGet("db").(*gorm.DB)
	profile, _ := getProfile(db, username)
	if err := c.ShouldBindJSON(&inputForm); err != nil {
		errorCode := errors.FormErrorCode
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode], "detail": err.Error()})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(inputForm.CurrentPassword)) != nil {
		errorCode := errors.AuthenticationFailureCode
		c.JSON(http.StatusBadRequest, gin.H{"error": errorCode, "msg": errors.Messages[errorCode]})
		return
	}

	db.Where("username = ?", username).Delete(models.User{})

	c.Data(http.StatusOK, gin.MIMEHTML, nil)
}
