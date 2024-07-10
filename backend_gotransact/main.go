package main

import (
	logger "gotransact/apps/Accounts"
	accounts "gotransact/apps/Accounts/models"
	validators "gotransact/apps/Accounts/validator"
	base "gotransact/apps/Base/models"
	"gotransact/apps/router"

	translogger "gotransact/apps/transaction"
	transaction "gotransact/apps/transaction/models"

	validators2 "gotransact/apps/transaction/validators"
	"gotransact/config"
	database "gotransact/pkg/db"

	extutils "gotransact/utils"
)



// @BasePath /api
// @title Go Transaction API
// @version 1.0
// @description This is a sample server for a transaction system.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email sangadevshatwar143@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8000

// @in header
// @name Authorization
func main() {
	config.Loadenv()
	database.InitDB("")
	validators.Init()
	validators2.InitValidation()
	logger.Init()
	translogger.Init()
	database.DB.AutoMigrate(&base.Base{}, &accounts.User{}, &accounts.Company{}, &transaction.PaymentGateway{}, &transaction.TransactionRequest{}, &transaction.TransactionHistory{})
	extutils.Cron()
	router.Routing()
	//Smartcontract = fabric.Initfabric()
}

// r.POST("/api/create", func(c *gin.Context) {
// 	gatway := transaction.PaymentGateway{
// 		Slug:  "card",
// 		Label: "Card",
// 	}
// 	if err := db.DB.Create(&gatway).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "error",
// 			Data:    map[string]interface{}{"data": err.Error()},
// 		})
// 		return
// 	}
// })

//////the above code is used for creating record in database about to add details about payment gateway
