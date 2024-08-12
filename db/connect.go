package db

import (
	"fmt"
	"os"

	"github.com/BIQ-Cat/easyserver/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var ModelsList []interface{}

func Connect() error {

	err := godotenv.Load()
	if err != nil {
		if config.Config.Debug {
			fmt.Println(fmt.Errorf("WARNING: %w", err))
		}
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Создать строку подключения
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
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
