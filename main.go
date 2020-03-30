package main

import (
	"invink/account-service/models"

	"github.com/gin-gonic/gin"
)

func main() {
	setupServer().Run()
}

func setupServer() *gin.Engine {
	r := gin.Default()
	db := models.SetupModels()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	setupRoutes(r)

	return r
}
