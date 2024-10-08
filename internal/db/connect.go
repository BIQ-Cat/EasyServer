package db

import (
	"fmt"

	config "github.com/BIQ-Cat/easyserver/config/base"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Enable postgres driver
)

var db *gorm.DB
var ModelsList []interface{}

func Connect() error {
	username := config.EnvConfig.DBUser
	password := config.EnvConfig.DBPass
	dbName := config.EnvConfig.DBName
	dbHost := config.EnvConfig.DBHost
	dbPort := config.EnvConfig.DBPort

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=disable password=%s", dbHost, username, dbName, dbPort, password) //Создать строку подключения

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		return err
	}

	db = conn
	db.Debug().AutoMigrate(ModelsList...)
	return nil
}

func GetDB() *gorm.DB {
	return db
}
