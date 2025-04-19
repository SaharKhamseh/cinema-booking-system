package database

import (
	"log"
	"os"

	"github.com/SaharKhamseh/cinema-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error .env file")
	}
	dsn := os.Getenv("DSN")
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	} else {
		log.Println("Connect successfully")
	}
	DB = database

	DB.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Theater{},
		&models.Screen{},
		&models.Seat{},
		&models.ShowTime{},
		&models.Booking{},
	)
}
