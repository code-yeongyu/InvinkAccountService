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
	r.POST("/register/", c.RegisterUser)
	r.POST("/auth/", c.AuthUser)

	r.POST("/capture/", middlewares.AuthenticateJWT, c.IncreaseCaptureAttemptCount)
	profile := r.Group("/profile")
	profile.Use(middlewares.AuthenticateJWT)
	{
		profile.DELETE("/:username/:field_name/", c.ClearFieldFromMyProfile)
		profile.GET("/:username/", c.GetProfileByUsername)
		profile.GET("/", c.GetMyProfile)
		profile.PATCH("/", c.UpdateMyProfile)
		profile.DELETE("/", c.RemoveMyProfile)
	}
}
