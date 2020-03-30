package main

import (
	"invink/account-service/docs"
	"invink/account-service/models"

	_ "invink/account-service/docs"

	"github.com/gin-gonic/gin"
)

func main() {
	docs.SwaggerInfo.Title = "Invink Account / Profile Service API Documentation"
	docs.SwaggerInfo.Description = "The Account/Profile Service API"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = "not done yet, possibly invink.org"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"https"}
	setupServer().Run()
}

func setupServer() *gin.Engine {
	r := gin.Default()
	db := models.Setup()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	setupRoutes(r)
	return r
}
