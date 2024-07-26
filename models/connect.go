package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var modelsList []interface{}

func init() {

	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
		err = nil
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Создать строку подключения
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		log.Fatal(err)
	}

	db = conn
	db.Debug().AutoMigrate(modelsList...)
}

func GetDB() *gorm.DB {
	return db
}
