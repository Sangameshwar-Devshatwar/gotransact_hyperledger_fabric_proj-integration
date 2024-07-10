package test

import (
	accountmodels "gotransact/apps/Accounts/models"
	"gotransact/apps/Base/models"

	//	"gotransact/apps/transaction/functions"
	logger "gotransact/apps/transaction"
	transactionmodels "gotransact/apps/transaction/models"
	"gotransact/apps/transaction/utils"
	"gotransact/pkg/db"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConfirmPayment_Success(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	user := accountmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}
	db.DB.Create(&user)

	gateway := transactionmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	transactionRequest := transactionmodels.TransactionRequest{
		UserID:                 user.ID,
		Status:                 transactionmodels.StatusProcessing,
		Description:            "Test Transaction",
		Amount:                 "100.0",
		PaymentGatewayMethodID: gateway.ID,
		Base: models.Base{
			InternalID: uuid.New(),
		},
	}
	db.DB.Create(&transactionRequest)

	status, message, _ := utils.ConfirmationPay(transactionRequest.InternalID.String(), "true")

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Transaction successfull", message)
	// assert.Equal(t, map[string]interface{}{}, data)

	var updatedTransaction transactionmodels.TransactionRequest
	db.DB.Where("internal_id = ?", transactionRequest.InternalID).First(&updatedTransaction)
	assert.Equal(t, transactionmodels.StatusSuccess, updatedTransaction.Status)

	var transactionHistory transactionmodels.TransactionHistory
	db.DB.Where("transaction_id = ?", updatedTransaction.ID).First(&transactionHistory)
	assert.Equal(t, transactionmodels.StatusSuccess, transactionHistory.Status)

	//CleanUp()
	CloseTestDb()
}

func TestConfirmPayment_Failed(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	user := accountmodels.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password",
	}
	db.DB.Create(&user)

	gateway := transactionmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	transactionRequest := transactionmodels.TransactionRequest{
		UserID:                 user.ID,
		Status:                 transactionmodels.StatusProcessing,
		Description:            "Test Transaction",
		Amount:                 "100.0",
		PaymentGatewayMethodID: gateway.ID,
		Base: models.Base{
			InternalID: uuid.New(),
		},
	}
	db.DB.Create(&transactionRequest)

	status, message, _ := utils.ConfirmationPay(transactionRequest.InternalID.String(), "false")

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Transaction Cancelled", message)
	//assert.Equal(t, map[string]interface{}{}, data)

	var updatedTransaction transactionmodels.TransactionRequest
	db.DB.Where("internal_id = ?", transactionRequest.InternalID).First(&updatedTransaction)
	assert.Equal(t, transactionmodels.StatusFailed, updatedTransaction.Status)

	var transactionHistory transactionmodels.TransactionHistory
	db.DB.Where("transaction_id = ?", updatedTransaction.ID).First(&transactionHistory)
	assert.Equal(t, transactionmodels.StatusFailed, transactionHistory.Status)
	// CleanUp()
	CloseTestDb()
}

func TestConfirmPayment_InvalidTransactionID(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	invalidTransactionID := "invalid-uuid"
	gateway := transactionmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	status, message, data := utils.ConfirmationPay(invalidTransactionID, "true")

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Invalid transaction ID", message)
	assert.Equal(t, map[string]interface{}{}, data)
	// CleanUp()
	CloseTestDb()
}

func TestConfirmPayment_TransactionNotFound(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validUUID := uuid.New().String()

	gateway := transactionmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	status, message, data := utils.ConfirmationPay(validUUID, "true")

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "transaction request not found", message)
	assert.Equal(t, map[string]interface{}{}, data)
	// CleanUp()
	CloseTestDb()
}
