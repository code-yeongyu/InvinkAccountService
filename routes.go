package main

import (
	"invink/account-service/controllers"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.RegisterUser)
}
