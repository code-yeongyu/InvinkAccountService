package controllers

import (
	"invink/account-service/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// IncreaseReportAttemptCount godoc
// @Summary Increase Report Attempt Count
// @Description Increase requested user's Report Attempt Count
// @Success 200 {object} EmptyResponse "Report Attempt Count Increased"
// @Router /attempt/report [post]
func (ctrler *Controller) IncreaseReportAttemptCount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	ID := c.MustGet("id").(uint64)
	var user models.User
	user.SetUserByID(db, ID)
	user.ReportCnt++
	db.Save(&user)
}
