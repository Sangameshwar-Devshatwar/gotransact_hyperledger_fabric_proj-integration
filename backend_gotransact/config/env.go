package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbTimezone string
)

func Loadenv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")
	DbTimezone = os.Getenv("DB_TIMEZONE")
}

//dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", dbHost, dbUser, dbPassword, dbName, dbPort, dbTimezone)
