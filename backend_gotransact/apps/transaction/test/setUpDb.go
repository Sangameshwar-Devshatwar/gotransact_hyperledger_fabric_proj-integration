package test

import (
	"fmt"
	ac "gotransact/apps/Accounts/models"
	"gotransact/apps/transaction/models"
	"gotransact/pkg/db"
)

func InitTestDB() {
	// dsn := "host=localhost user=trellis password=postgres dbname=mydatabase port=5432 sslmode=disable"
	// var err error
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return err
	// }
	// sqlDB, err := DB.DB()
	// if err != nil {
	// 	return err
	// }
	// _, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	// if err != nil {
	// 	return err
	// }
	// //Run your migrations here if needed
	// err = DB.AutoMigrate(&models.User{}, &models.Company{})
	// if err != nil {
	// 	return err
	// }
	// return nil
	db.InitDB("test")
	if err := db.DB.AutoMigrate(&ac.User{}, &ac.Company{}, &models.PaymentGateway{}, &models.TransactionRequest{}, &models.TransactionHistory{}); err != nil {
		fmt.Printf("Error autoigrating models : %s", err.Error())
	}
}

func CleanUp() {
	sqlDB, _ := db.DB.DB()

	sqlDB.Exec("DELETE FROM companies")
	sqlDB.Exec("DELETE FROM users")
	sqlDB.Exec("DELETE FROM payment_gateways")
	sqlDB.Exec("DELETE FROM transaction_histories")
	sqlDB.Exec("DELETE FROM transaction_requests")

}
func CloseTestDb() {
	fmt.Println("---------------in close db-------------")
	sqlDB, err := db.DB.DB()
	if err != nil {
		fmt.Printf("Error getting sqlDB from gorm DB: %s", err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Printf("Error closing database: %s", err.Error())
	}
}
