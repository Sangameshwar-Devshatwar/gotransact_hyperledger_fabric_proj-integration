package test

import (
	"gotransact/apps/Accounts/models"
	"gotransact/apps/Accounts/utils"
	validators "gotransact/apps/Accounts/validator"
	astructutils "gotransact/apps/Astructutils"
	"gotransact/pkg/db"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// var DB *gorm.DB

// func ClearDatabase() {
// 	fmt.Println("in clear db --------------1")
// 	db.DB.Exec("DELETE FROM companies")
// 	fmt.Println("in clear db --------------2")
// 	db.DB.Exec("DELETE FROM users")
// 	fmt.Println("in clear db --------------3")

// }

func TestSignup_success(t *testing.T) {
	InitTestDB()
	CleanUp()
	validators.Init()
	// logger.Init()

	input := astructutils.SignupUser{
		FirstName:   "testfirstname",
		LastName:    "testlastname",
		Email:       "test@gmail.com",
		CompanyName: "trellissoft",
		Password:    "Password@123",
	}

	status, message, data := utils.Signup(input)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "User Created Successfully", message)
	assert.Equal(t, map[string]interface{}{}, data)

	var user models.User
	err := db.DB.Where("email = ?", input.Email).First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, input.FirstName, user.FirstName)
	assert.Equal(t, input.LastName, user.LastName)
	assert.Equal(t, input.Email, user.Email)

	var company models.Company
	err = db.DB.Where("name = ?", input.CompanyName).First(&company).Error
	assert.NoError(t, err)
	assert.Equal(t, input.CompanyName, company.Name)
	CloseTestDb()
}

func TestSignup_EmailAreadyExist(t *testing.T) {
	InitTestDB()
	CleanUp()
	// logger.Init()
	validators.Init()

	existingUser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Company: models.Company{
			Name: "trellissoft",
		},
	}

	db.DB.Create(&existingUser)

	input := astructutils.SignupUser{
		FirstName:   "otherfirstname",
		LastName:    "otherlastname",
		Email:       "test@gmail.com",
		CompanyName: "Google",
		Password:    "Password@123",
	}

	status, message, data := utils.Signup(input)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "email already exists", message)
	assert.Equal(t, map[string]interface{}{}, data)
	CloseTestDb()
}

func TestSignup_CompanyAlreadyExist(t *testing.T) {
	InitTestDB()
	CleanUp()
	validators.Init()
	// logger.Init()

	existingUser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test1@gmail.com",
		Company: models.Company{
			Name: "trellissoft",
		},
	}

	db.DB.Create(&existingUser)

	input := astructutils.SignupUser{
		FirstName:   "otherfirstname",
		LastName:    "otherlastname",
		Email:       "test@gmail.com",
		CompanyName: "trellissoft",
		Password:    "Password@123",
	}

	status, message, data := utils.Signup(input)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "company already exists", message)
	assert.Equal(t, map[string]interface{}{}, data)
	CloseTestDb()
}

func TestSignup_InvaldPassword(t *testing.T) {
	InitTestDB()
	CleanUp()
	validators.Init()
	// logger.Init()

	input := astructutils.SignupUser{
		FirstName:   "otherfirstname",
		LastName:    "otherlastname",
		Email:       "testgmail.com",
		CompanyName: "trellissoft",
		Password:    "password@123",
	}

	status, message, data := utils.Signup(input)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Password should contain atleast one upper case character,one lower case character,one number and one special character", message)
	assert.Equal(t, map[string]interface{}{}, data)
	CloseTestDb()
}
