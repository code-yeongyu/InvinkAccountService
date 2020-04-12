package controllers

import (
	"invink/account-service/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// IncreaseCaptureAttemptCount godoc
// @Summary Register an user
// @Description Register an user with given information
// @Success 200 {object} EmptyResponse "Capture Attempt Count Increased"
// @Router /capture/ [post]
func (ctrler *Controller) IncreaseCaptureAttemptCount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ID := c.MustGet("id").(uint64)
	var user models.User
	user.SetUserByID(db, ID)
	user.CaptureCnt++
	db.Save(&user)
}
