package main

import (
	"invink/account-service/controllers"

	"github.com/gin-gonic/gin"

	_ "invink/account-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRoutes(r *gin.Engine) {
	c := controllers.NewController()
	r.POST("/register", c.RegisterUser)
	r.POST("/auth", c.AuthUser)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
