package main

import (
	"invink/account-service/docs"

	_ "invink/account-service/docs"

	"github.com/gin-gonic/gin"
)

func main() {
	docs.SwaggerInfo.Title = "Invink Account / Profile Service API Documentation"
	docs.SwaggerInfo.Description = "The Account/Profile Service API"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = "possibly invink.org"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"https"}
	setupServer().Run(":3000")
}

func setupServer() *gin.Engine {
	r := gin.Default()
	setupRoutes(r)
	return r
}
