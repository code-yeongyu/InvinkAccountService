package controllers

import (
	"invink/account-service/models"
	"strings"
	"unicode"

	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getProfile(db *gorm.DB, username string) (user models.User, err error) {
	err = db.Where("username = ?", username).First(&user).Error
	return user, err
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
// @Failure 404
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
