package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // to manage the postgres database
)

// SetupModels initializes and migrates models
func SetupModels() *gorm.DB {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("ACCOUNT_DB_HOST"), os.Getenv("ACCOUNT_DB_PORT"),
		os.Getenv("ACCOUNT_DB_USER"), os.Getenv("ACCOUNT_DB_DBNAME"),
		os.Getenv("ACCOUNT_DB_PASSWORD"), os.Getenv("ACCOUNT_DB_SSLMODE"))
	db, err := gorm.Open("postgres", dbInfo)

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&UserModel{})
	return db
}
