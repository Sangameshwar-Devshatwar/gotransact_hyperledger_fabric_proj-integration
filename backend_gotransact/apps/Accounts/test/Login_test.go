package test

import (
	logger "gotransact/apps/Accounts"
	"gotransact/apps/Accounts/models"
	"gotransact/apps/Accounts/utils"
	validators "gotransact/apps/Accounts/validator"
	astructutils "gotransact/apps/Astructutils"
	"gotransact/pkg/db"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_success(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validators.Init()

	// Create a test user
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)
	user := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Password:  string(passwordHash),
	}
	db.DB.Create(&user)

	input := astructutils.LoginInput{
		Email:    "test@gmail.com",
		Password: "Password@123",
	}

	status, message, data := utils.Login(input)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "logged in successfully", message)
	assert.NotEmpty(t, data["data"])

	CloseTestDb()
}
func TestLogin_validationError(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validators.Init()

	input := astructutils.LoginInput{
		Email:    "invalid-email",
		Password: "short",
	}

	status, message, data := utils.Login(input)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "error while validating fields", message)
	assert.Empty(t, data)

	CloseTestDb()
}
func TestLogin_userNotFound(t *testing.T) {
	InitTestDB()
	CleanUp()
	logger.Init()
	validators.Init()

	input := astructutils.LoginInput{
		Email:    "notfound@gmail.com",
		Password: "Password@123",
	}

	status, message, data := utils.Login(input)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, "no user found in database", message)
	assert.Empty(t, data)

	CloseTestDb()
}
func TestLogin_invalidPassword(t *testing.T) {
	InitTestDB()
	logger.Init()
	CleanUp()
	validators.Init()

	// Create a test user
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)
	user := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Password:  string(passwordHash),
	}
	db.DB.Create(&user)

	input := astructutils.LoginInput{
		Email:    "test@gmail.com",
		Password: "InvalidPassword@123",
	}

	status, message, data := utils.Login(input)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, "invalid password", message)
	assert.Empty(t, data)
	CleanUp()
	CloseTestDb()
}
