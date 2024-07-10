package utils

import (
	"context"
	"fmt"
	logger "gotransact/apps/Accounts"
	"gotransact/apps/Accounts/models"
	validators "gotransact/apps/Accounts/validator"
	astructutils "gotransact/apps/Astructutils"
	"gotransact/fabric"
	"gotransact/pkg/db"
	"log"
	"net/http"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
)

func SendMail(mail string) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted Sendmail() email to", mail)
	abc := gomail.NewMessage()

	abc.SetHeader("From", "sangatrellis123@gmail.com")
	abc.SetHeader("To", mail)
	abc.SetHeader("Subject", "confirmation")
	abc.SetBody("text/plain", "user has been created successfully ,this is a confirmation mail")

	a := gomail.NewDialer("smtp.gmail.com", 587, "sangatrellis123@gmail.com", "mhch pnah ljze lsyw")
	if err := a.DialAndSend(abc); err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error sending mail")
		log.Fatal(err.Error())
	}
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("mail sent successfully to ", mail)
}

// implementing creation of token
var secretKey = paseto.NewV4AsymmetricSecretKey() // don't share this!!!
var publicKey = secretKey.Public()

func GeneratePasetoToken(user models.User) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetExpiration(exp)

	token.Set("User", user)
	signed := token.V4Sign(secretKey, nil)
	return signed, nil
}

func VerifyPasetoToken(signed string) (any, error) {
	val, err := rdb.Get(ctx, signed).Result()
	if err == nil && val == "Blacklisted" {
		return nil, fmt.Errorf("token has been revoked")
	}

	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	token, err := parser.ParseV4Public(publicKey, signed, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token %w", err)
	}
	var user models.User
	err = token.Get("User", &user)
	if err != nil {
		return nil, fmt.Errorf("subject claim not found in token")
	}
	return user, nil
}
func Signup(requestInputs astructutils.SignupUser) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted signup method with email", requestInputs.Email, "and company", requestInputs.CompanyName)
	if err := validators.GetValidator().Struct(requestInputs); err != nil {
		return http.StatusBadRequest, "Password should contain atleast one upper case character,one lower case character,one number and one special character", map[string]interface{}{}
	}
	var count int64
	if err := db.DB.Model(&models.User{}).Where("email = ?", requestInputs.Email).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "email already exists", map[string]interface{}{}
	}
	if err := db.DB.Model(&models.Company{}).Where("name = ?", requestInputs.CompanyName).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "company already exists", map[string]interface{}{}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestInputs.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, "failed to hash the password", map[string]interface{}{}
	}

	//values stored into original model to add user to database
	user := models.User{
		FirstName: requestInputs.FirstName,
		LastName:  requestInputs.LastName,
		Email:     requestInputs.Email,
		Password:  string(hashedPassword),
		Company: models.Company{
			Name: requestInputs.CompanyName,
		},
	}
	//user created by using user passesd values binded into struct
	if err := db.DB.Create(&user).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error while creating database")
		return http.StatusInternalServerError, "failed to create user database", map[string]interface{}{}
	}
	Smartcontract := fabric.Initfabric()
	result, err := Smartcontract.SubmitTransaction("Userregister", user.FirstName, user.LastName, user.Email, user.Company.Name)
	if err != nil {
		log.Printf("Failed to submit transaction: %v", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return http.StatusInternalServerError, "failed to submit transaction", map[string]interface{}{}
	}
	var userss models.User
	if err := db.DB.Model(&userss).Where("email = ?", user.Email).Update("transaction_hash", result).Error; err != nil {
		return http.StatusInternalServerError, "failed to update status", map[string]interface{}{}
	}
	//sent mail using goroutines
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("User created successfully with email", requestInputs.Email, "and company", requestInputs.CompanyName)
	go SendMail(user.Email)
	return http.StatusOK, "User Created Successfully", map[string]interface{}{}
}

// login function
func Login(inputuser astructutils.LoginInput) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted login method using mail", inputuser.Email)

	if err := validators.GetValidator().Struct(inputuser); err != nil {
		return http.StatusBadRequest, "error while validating fields", map[string]interface{}{}
	}
	var user models.User
	if err := db.DB.Where("email=?", inputuser.Email).First(&user).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error finding mail in database")
		return http.StatusInternalServerError, "no user found in database", map[string]interface{}{}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputuser.Password)); err != nil {
		return http.StatusInternalServerError, "invalid password", map[string]interface{}{}
	}

	token, err := GeneratePasetoToken(user)
	if err != nil {
		return http.StatusInternalServerError, "error generating token", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("Logged in successfully using mail", inputuser.Email)
	return http.StatusOK, "logged in successfully", map[string]interface{}{"data": token}
}

// logout function to implement logout handler
func Logout(authHeader string) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"message": "in logout func",
	}).Info("attempted logout method")

	if authHeader == "" {
		return http.StatusUnauthorized, "authorization header missing", map[string]interface{}{}
	}

	//tokenStr := authHeader[len("Bearer "):]

	_, err := VerifyPasetoToken(authHeader)
	if err != nil {
		return http.StatusUnauthorized, "invalid token", map[string]interface{}{}
	}

	// Blacklist the token by storing it in Redis with an expiration time
	err = rdb.Set(ctx, authHeader, "Blacklisted", 24*time.Hour).Err() // adjust expiration time as needed
	if err != nil {
		return http.StatusInternalServerError, "failed to blacklist token", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{"logout": "success"}).Info("Logged out successfully")
	return http.StatusOK, "logged out successfully", map[string]interface{}{}
}
