package database

import (
	"finalproject/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var (
	host    = "localhost"
	port    = 5432
	user    = "postgres"
	pasword = "root"
	dbname  = "test"
	db      *gorm.DB
	err     error
)

func ConnectToDatabase() {
	config := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pasword, dbname)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to databases :", err)
	}
	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	return db
}
