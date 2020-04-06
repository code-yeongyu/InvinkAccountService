package main

import (
	"invink/account-service/controllers"
	"invink/account-service/middlewares"

	"github.com/gin-gonic/gin"

	_ "invink/account-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRoutes(r *gin.Engine) {
	r.Use(middlewares.SetupDB)
	c := controllers.NewController()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/register", c.RegisterUser)
	r.POST("/auth", c.AuthUser)

	profile := r.Group("/profile")
	profile.Use(middlewares.AuthenticateJWT)
	{
		profile.GET("/", c.GetMyUsername)
		profile.GET("/:username/", c.GetProfileByUsername)
		profile.PATCH("/:username", nil)
		profile.DELETE("/:username", nil)
		profile.DELETE("/:username/:field_name", nil)
	}

}
