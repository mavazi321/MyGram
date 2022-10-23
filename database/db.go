package database

import (
	"finalproject/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var (
	host    = os.Getenv("PGHOST")
	port    = os.Getenv("PGPORT")
	user    = os.Getenv("PGUSER")
	pasword = os.Getenv("PGPASSWORD")
	dbname  = os.Getenv("PGDATABASE")
	db      *gorm.DB
	err     error
)

func ConnectToDatabase() {
	config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pasword, dbname)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to databases :", err)
	}
	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	return db
}
