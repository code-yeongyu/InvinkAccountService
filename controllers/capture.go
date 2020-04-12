package controllers

import (
	"invink/account-service/models"

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
// @Success 200 {object} EmptyResponse "Capture At"
// @Failure 400 {object} TypicalErrorResponse "Normal Form error, like username duplicate"
// @Router /register/ [post]
func (ctrler *Controller) IncreaseCaptureAttemptCount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ID := c.MustGet("id").(uint64)
	var user models.User
	user.SetUserByID(db, ID)
	user.CaptureCnt++
	db.Save(&user)
}
