package middlewares

import (
	"invink/account-service/models"

	"github.com/gin-gonic/gin"
)

// SetupDB is a middlware for the database connection
func SetupDB(c *gin.Context) {
	c.Set("db", models.Setup())
	c.Next()
}
