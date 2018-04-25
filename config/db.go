package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB function
func DB() *gorm.DB {

	user := "postgres"  //os.Getenv("DB_USER")
	password := "1234"  //os.Getenv("DB_PASSWORD")
	host := "localhost" //os.Getenv("DB_HOST")
	port := "5433"
	_db := "socialapp" //os.Getenv("DB")
	conStr := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", host, port, user, _db, password)

	db, err := gorm.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}
	return db
}
