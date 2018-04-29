package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	models "github.com/peacecwz/go-social-app/models"
)

func InitDB() {
	db := DB()

	db.AutoMigrate(&(models.User{}))
	db.AutoMigrate(&(models.Post{}))
	db.AutoMigrate(&(models.ProfileView{}))
	db.AutoMigrate(&(models.Like{}))
	db.AutoMigrate(&(models.Follow{}))
}

// DB function
func DB() *gorm.DB {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	_db := os.Getenv("DB")
	conStr := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", host, port, user, _db, password)

	db, err := gorm.Open(os.Getenv("DB_TYPE"), conStr)
	if err != nil {
		panic(err)
	}
	if os.Getenv("MODE") == "DEV" {
		return db.Debug()
	}

	return db
}
