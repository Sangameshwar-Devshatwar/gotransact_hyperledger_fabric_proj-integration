package utils

import (
	"bytes"
	"fmt"
	"gotransact/apps/Accounts/models"
	logger "gotransact/apps/transaction"
	mod "gotransact/apps/transaction/models"
	"gotransact/apps/transaction/structutils"
	"gotransact/apps/transaction/validators"
	"gotransact/fabric"
	"gotransact/pkg/db"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

type EmailData struct {
	Username         string
	TransactionID    uuid.UUID
	Amount           string
	DateTime         time.Time
	ConfirmationLink string
	FailPaymentLink  string
	AdminEmail       string
	CompanyName      string
	CompanyContact   string
	CompanyWebsite   string
}

func SendMail(user models.User, transdetails mod.TransactionRequest) {
	fmt.Println("====================startof the mail=============")
	tmpl, err := template.ParseFiles("/home/trellis/Desktop/folder-1/golang-training/backend_gotransact/apps/transaction/utils/mail_template.html")
	if err != nil {
		log.Fatal("Error loading email template:", err)
	}
	//iid, err := uuid.Parse(transdetails.InternalID)

	data := EmailData{
		Username:         user.FirstName,
		TransactionID:    transdetails.InternalID,
		Amount:           transdetails.Amount,
		DateTime:         transdetails.CreatedAt,
		ConfirmationLink: "http:" + "//localhost:8000/api/confirm_payment?transactionid=" + transdetails.InternalID.String() + "&status=true", // Adjust the confirmation link as needed
		FailPaymentLink:  "http://" + "localhost:8000/api/confirm_payment?transactionid=" + transdetails.InternalID.String() + "&status=false",
		AdminEmail:       "Admin123@gmail.com", // Adjust the admin email
		CompanyName:      "trellisMagic",
		CompanyContact:   "Tower 4,infocity ,Gandhinagar,Gujrat",
		CompanyWebsite:   "http://trellisMagic.com",
	}

	// Execute the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Fatal("Error executing email template:", err)
	}

	// Create the email
	m := gomail.NewMessage()
	m.SetHeader("From", "sangatrellis123@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Payment Confirmation")
	m.SetBody("text/html", body.String())

	// Configure the SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "sangatrellis123@gmail.com", "mhch pnah ljze lsyw")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Error sending email:", err)
	}
	fmt.Println("====================end of the mail=============")

}
func FetchTransactionsLast24Hours() []mod.TransactionRequest {
	var transactions []mod.TransactionRequest
	last24Hours := time.Now().Add(-24 * time.Hour)
	db.DB.Where("created_at >= ?", last24Hours).Find(&transactions)
	return transactions
}
func GenerateExcel(transactions []mod.TransactionRequest) (string, error) {
	f := excelize.NewFile()
	sheetName := "Transactions"
	index := f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", "ID")
	f.SetCellValue(sheetName, "B1", "InternalID")
	f.SetCellValue(sheetName, "C1", "UserID")
	f.SetCellValue(sheetName, "D1", "Status")
	f.SetCellValue(sheetName, "E1", "PaymentGatewayID")
	f.SetCellValue(sheetName, "F1", "Description")
	f.SetCellValue(sheetName, "G1", "Amount")
	f.SetCellValue(sheetName, "H1", "CreatedAt")
	f.SetCellValue(sheetName, "I1", "UpdatedAt")

	for i, tr := range transactions {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), tr.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), tr.InternalID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), tr.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), tr.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), tr.PaymentGatewayMethodID)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), tr.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), tr.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), tr.CreatedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), tr.UpdatedAt)
	}

	f.SetActiveSheet(index)
	filePath := "transactions.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func SendMailWithAttachment(email, filePath string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "sangatrellis123@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Daily Transactions Report")
	m.SetBody("text/plain", "Please find attached the daily transactions report.")
	m.Attach(filePath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "sangatrellis123@gmail.com", "mhch pnah ljze lsyw")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("could not send email: %v", err)
	}
	fmt.Println("Email sent successfully")
}
func PostPayment(payreq structutils.PaymentRequest, user models.User) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted postpayment method with email", user.Email, " and company ", user.Company)
	if err := validators.GetValidator().Struct(payreq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, fieldErr := range validationErrors {
				fieldName := fieldErr.Field()
				// Using the field name to get the custom error message
				if customMessage, found := validators.CustomErrorMessages[fieldName]; found {
					errors[fieldName] = customMessage
				} else {
					errors[fieldName] = err.Error() // Default error message
				}
			}
			return http.StatusBadRequest, "Validation error", map[string]interface{}{}
		}

		return http.StatusBadRequest, "error while validating", map[string]interface{}{}
	}

	var paymentgateway mod.PaymentGateway
	if err := db.DB.Where("slug=?", "card").First(&paymentgateway).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error searching for slug")
		return http.StatusInternalServerError, "no user found in database", map[string]interface{}{}
	}
	//var TransactionReq transmodel.TransactionRequest
	transactiondetails := mod.TransactionRequest{
		UserID:                 user.ID,
		Status:                 mod.StatusProcessing,
		Description:            payreq.Description,
		Amount:                 payreq.Amount,
		PaymentGatewayMethodID: paymentgateway.ID,
	}
	if err := db.DB.Create(&transactiondetails).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error creating transaction details in database")
		return http.StatusInternalServerError, "failed to create transactiondetails in database", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{
		"data": "transaction details",
	}).Info("transaction details added to database successfully")
	transactionHistory := mod.TransactionHistory{
		TransactionID: transactiondetails.ID,
		Status:        transactiondetails.Status,
		Description:   transactiondetails.Description,
		Amount:        transactiondetails.Amount,
	}
	if err := db.DB.Create(&transactionHistory).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error adding transaction to transaction history")
		return http.StatusInternalServerError, "failed to create transactiondetails in database", map[string]interface{}{}
	}
	Smartcontract := fabric.Initfabric()
	resulttx, err := Smartcontract.SubmitTransaction("PaymentInitiate", transactionHistory.Status, transactionHistory.Amount, user.Email, string(transactionHistory.TransactionID))
	if err != nil {
		log.Printf("Failed to submit transaction: %v", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return http.StatusInternalServerError, "failed to submit transaction", map[string]interface{}{}
	}
	fmt.Println(string(resulttx))
	var transreqs mod.TransactionHistory
	if err := db.DB.Model(&transreqs).Where("transaction_id = ?", transactiondetails.ID).Update("transaction_hash", string(resulttx)).Error; err != nil {
		return http.StatusInternalServerError, "failed to update status", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("transaction details added to transaction history successfully")

	go SendMail(user, transactiondetails)

	return http.StatusOK, "mail sent for confirmation of transaction to user", map[string]interface{}{}
}
func ConfirmationPay(id string, status string) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted confirm payement method", "transaction id:", id)
	transid, err := uuid.Parse(id)
	if err != nil {
		return http.StatusBadRequest, "Invalid transaction ID", map[string]interface{}{}
	}

	if status == "true" {
		var transreq mod.TransactionRequest
		if err := db.DB.Where("internal_id=?", transid).First(&transreq).Error; err != nil {
			return http.StatusBadRequest, "transaction request not found", map[string]interface{}{}
		}

		// if err:=db.DB.Update("status",models.StatusSuccess)
		if err := db.DB.Model(&transreq).Where("internal_id = ?", transid).Update("status", mod.StatusSuccess).Error; err != nil {
			logger.ErrorLogger.WithFields(logrus.Fields{"error": err.Error()}).Error("error whike updating status as success")
			return http.StatusInternalServerError, "failed to update status", map[string]interface{}{}
		}
		transhistory := mod.TransactionHistory{
			TransactionID: transreq.ID,
			Status:        mod.StatusSuccess,
			Description:   transreq.Description,
			Amount:        transreq.Amount,
		}
		if err := db.DB.Create(&transhistory).Error; err != nil {
			logger.ErrorLogger.WithFields(logrus.Fields{"error": err.Error()}).Error("error while creating transaction history in database")
			return http.StatusInternalServerError, "error while creating record in database", map[string]interface{}{}
		}
		Smartcontract := fabric.Initfabric()
		result, err := Smartcontract.SubmitTransaction("ConfirmPayment", string(transhistory.TransactionID), transhistory.Status, transhistory.Amount)
		if err != nil {
			log.Printf("Failed to submit transaction: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return http.StatusInternalServerError, "failed to submit transaction", map[string]interface{}{}
		}
		fmt.Println(string(result))
		var transreqs mod.TransactionHistory
		if err := db.DB.Model(&transreqs).Where("id = ?", transhistory.ID).Update("transaction_hash", result).Error; err != nil {
			return http.StatusInternalServerError, "failed to update status", map[string]interface{}{}
		}
		logger.InfoLogger.WithFields(logrus.Fields{}).Info("transaction successfull")
		return http.StatusOK, "Transaction successfull", map[string]interface{}{
			"TransactionID": transreq.InternalID,
			"Amount":        transreq.Amount,
			"DateTime":      transreq.CreatedAt,
		}

	} else if status == "false" {
		var transreq mod.TransactionRequest
		if err := db.DB.Where("internal_id=?", transid).First(&transreq).Error; err != nil {
			return http.StatusInternalServerError, "can't find in database", map[string]interface{}{}
		}

		// if err:=db.DB.Update("status",models.StatusSuccess)
		if err := db.DB.Model(&transreq).Where("internal_id = ?", transid).Update("status", mod.StatusFailed).Error; err != nil {
			return http.StatusInternalServerError, "failed to update status", map[string]interface{}{}
		}
		transhistory := mod.TransactionHistory{
			TransactionID: transreq.ID,
			Status:        mod.StatusFailed,
			Description:   transreq.Description,
			Amount:        transreq.Amount,
		}
		if err := db.DB.Create(&transhistory).Error; err != nil {
			logger.ErrorLogger.WithFields(logrus.Fields{"error": err.Error()}).Error("error while creating transaction history in database")
			return http.StatusInternalServerError, "error while creating record in database", map[string]interface{}{}
		}
		Smartcontract := fabric.Initfabric()
		result, err := Smartcontract.SubmitTransaction("ConfirmPayment", string(transhistory.TransactionID), transhistory.Status, transhistory.Amount)
		if err != nil {
			log.Printf("Failed to submit transaction: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return http.StatusInternalServerError, "failed to submit transaction", map[string]interface{}{}
		}
		var transreqs mod.TransactionRequest
		if err := db.DB.Model(&transreqs).Where("transaction_id = ?", transhistory.TransactionID).Update("transaction_hash", result).Error; err != nil {
			return http.StatusInternalServerError, "failed to update status", map[string]interface{}{}
		}
		logger.InfoLogger.WithFields(logrus.Fields{}).Info("transaction Cancelled")
		return http.StatusOK, "Transaction Cancelled", map[string]interface{}{
			"TransactionID": transreq.InternalID,
			"Amount":        transreq.Amount,
			"DateTime":      transreq.CreatedAt,
		}
	}

	logger.InfoLogger.WithFields(logrus.Fields{}).Info("transaction successfull")

	return http.StatusOK, "Transaction successfull", map[string]interface{}{}
}
