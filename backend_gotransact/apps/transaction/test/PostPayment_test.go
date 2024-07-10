package test

import (
	"gotransact/apps/Accounts/models"
	logger "gotransact/apps/transaction"
	transmodels "gotransact/apps/transaction/models"
	"gotransact/apps/transaction/structutils"
	"gotransact/apps/transaction/utils"
	"gotransact/apps/transaction/validators"
	"gotransact/pkg/db"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// func ClearDatabase() {
// 	fmt.Println("in clear db --------------1")
// 	db.DB.Exec("DELETE FROM companies")
// 	fmt.Println("in clear db --------------2")
// 	db.DB.Exec("DELETE FROM users")
// 	fmt.Println("in clear db --------------3")

// }
func TestPostPayment_success(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validators.InitValidation()

	// Create a test user
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)
	user := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "tests@gmail.com",
		Password:  string(passwordHash),
		Company: models.Company{
			Name: "goggle",
		},
	}
	db.DB.Create(&user)

	// Create a payment gateway
	paymentGateway := transmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&paymentGateway)

	input := structutils.PaymentRequest{
		CardNumber:  "1234567812345678",
		ExpiryDate:  "12/25",
		CVV:         "123",
		Amount:      "100",
		Description: "Test payment",
	}

	status, message, data := utils.PostPayment(input, user)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "mail sent for confirmation of transaction to user", message)
	assert.Equal(t, map[string]interface{}{}, data)

	var transaction transmodels.TransactionRequest
	err := db.DB.Where("user_id = ?", user.ID).First(&transaction).Error
	assert.NoError(t, err)
	assert.Equal(t, input.Description, transaction.Description)
	assert.Equal(t, input.Amount, transaction.Amount)
	// CleanUp()
	CloseTestDb()
}
func TestPostPayment_validationCard(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validators.InitValidation()

	input := structutils.PaymentRequest{
		CardNumber:  "65423152",
		ExpiryDate:  "12/25",
		CVV:         "123",
		Amount:      "100",
		Description: "Test payment",
	}

	user := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
	}

	status, message, data := utils.PostPayment(input, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Validation error", message)
	assert.Empty(t, data)

	CloseTestDb()
}
func TestPostPayment_validationExpiryDate(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validators.InitValidation()

	input := structutils.PaymentRequest{
		CardNumber:  "6542315212345678",
		ExpiryDate:  "12/15",
		CVV:         "123",
		Amount:      "100",
		Description: "Test payment",
	}

	user := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
	}

	status, message, data := utils.PostPayment(input, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Validation error", message)
	assert.Empty(t, data)

	CloseTestDb()
}
func TestPostPayment_validationCVV(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validators.InitValidation()

	input := structutils.PaymentRequest{
		CardNumber:  "6542315212345678",
		ExpiryDate:  "12/15",
		CVV:         "123343",
		Amount:      "100",
		Description: "Test payment",
	}

	user := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
	}

	status, message, data := utils.PostPayment(input, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Validation error", message)
	assert.Empty(t, data)

	CloseTestDb()
}
func TestPostPayment_validationAmount(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validators.InitValidation()

	input := structutils.PaymentRequest{
		CardNumber:  "6542315212345678",
		ExpiryDate:  "12/15",
		CVV:         "123",
		Amount:      "hefiwhefeqb",
		Description: "Test payment",
	}

	user := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
	}

	status, message, data := utils.PostPayment(input, user)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Validation error", message)
	assert.Empty(t, data)

	CloseTestDb()
}
