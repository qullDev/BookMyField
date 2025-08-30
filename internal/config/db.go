package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	sqlite_driver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabse() {
	dsn := os.Getenv("DATABASE_URL")

	var db *gorm.DB
	var err error

	if dsn == "" {
		// Development fallback - use SQLite
		log.Println("⚠️ DATABASE_URL not found, using SQLite for development")
		db, err = gorm.Open(sqlite_driver.Open("bookmyfield.db"), &gorm.Config{})
	} else {
		// Production - use PostgreSQL
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Error connecting to database: ", err.Error())
	}

	DB = db
	fmt.Println("✅ Database connected")
}
