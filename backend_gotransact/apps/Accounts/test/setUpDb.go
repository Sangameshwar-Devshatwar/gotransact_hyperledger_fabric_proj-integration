package test

import (
	"fmt"
	"gotransact/apps/Accounts/models"
	"gotransact/pkg/db"
)

func InitTestDB() {
	db.InitDB("test")
	if err := db.DB.AutoMigrate(&models.User{}, &models.Company{}); err != nil {
		fmt.Printf("Error autoigrating models : %s", err.Error())
	}
}

func CleanUp() {
	sqlDB, _ := db.DB.DB()
	sqlDB.Exec("TRUNCATE TABLE users CASCADE")
	sqlDB.Exec("TRUNCATE TABLE companies CASCADE")
}
func CloseTestDb() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		fmt.Printf("Error getting sqlDB from gorm DB: %s", err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Printf("Error closing database: %s", err.Error())
	}
}
