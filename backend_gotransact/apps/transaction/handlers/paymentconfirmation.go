package handlers

import (
	logger "gotransact/apps/transaction"
	"gotransact/apps/transaction/utils"
	"gotransact/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)
// @BasePath /api
// ConfirmationPayment handles payment confirmation.
// @Summary Confirm a payment transaction
// @Description Confirm a payment transaction based on transaction ID and status
// @Tags Payment
// @Accept json
// @Produce json
// @Param transactionid query string true "Transaction ID"
// @Param status query string true "Status of the payment confirmation (true/false)"
// @Failure 400 {object} responses.UserResponse "Invalid input"
// @Failure 500 {object} responses.UserResponse "Internal Server Error"
// @Router /confirm_payment [get]
func ConfirmationPayment(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("attempted confirmpayment")
	id := c.Query("transactionid")
	status := c.Query("status")
	statuses, message, data := utils.ConfirmationPay(id, status)

	if statuses == http.StatusOK && status == "true" {
		c.HTML(http.StatusOK, "success.html", gin.H{
			"TransactionID": data["TransactionID"],
			"Amount":        data["Amount"],
			"DateTime":      data["DateTime"],
		})
	} else if statuses == http.StatusOK && status == "false" {
		c.HTML(http.StatusOK, "failure.html", gin.H{
			"TransactionID": data["TransactionID"],
			"Amount":        data["Amount"],
			"DateTime":      data["DateTime"],
		})
	} else {
		c.JSON(statuses, responses.UserResponse{
			Status:  statuses,
			Message: message,
			Data:    data,
		})
	}
}

//code for reference if anything goes wrong in function

// transid, err := uuid.Parse(id)
// if err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "error while parsing internal_id"})
// 	return
// }

// if status == "true" {
// 	var transreq models.TransactionRequest
// 	if err := db.DB.Where("internal_id=?", transid).First(&transreq).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "can't find in database",
// 			Data: map[string]interface{}{
// 				"error": err.Error(),
// 			},
// 		})
// 	}

// 	// if err:=db.DB.Update("status",models.StatusSuccess)
// 	if err := db.DB.Model(&transreq).Where("internal_id = ?", transid).Update("status", models.StatusSuccess).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "failed to update status",
// 			Data: map[string]interface{}{
// 				"data": err.Error(),
// 			},
// 		})
// 		return
// 	}
// 	transhistory := models.TransactionHistory{
// 		TransactionID: transreq.ID,
// 		Status:        models.StatusSuccess,
// 		Description:   transreq.Description,
// 		Amount:        transreq.Amount,
// 	}
// 	if err := db.DB.Create(&transhistory).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "error while creating record in database",
// 			Data: map[string]interface{}{
// 				"data": err.Error(),
// 			},
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, responses.UserResponse{
// 		Status:  http.StatusOK,
// 		Message: "payment done",
// 		Data: map[string]interface{}{
// 			"data": "success",
// 		},
// 	})
// } else if status == "false" {
// 	var transreq models.TransactionRequest
// 	if err := db.DB.Where("internal_id=?", transid).First(&transreq).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "can't find in database",
// 			Data: map[string]interface{}{
// 				"error": err.Error(),
// 			},
// 		})
// 	}

// 	// if err:=db.DB.Update("status",models.StatusSuccess)
// 	if err := db.DB.Model(&transreq).Where("internal_id = ?", transid).Update("status", models.StatusFailed).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "failed to update status",
// 			Data: map[string]interface{}{
// 				"data": err.Error(),
// 			},
// 		})
// 		return
// 	}
// 	transhistory := models.TransactionHistory{
// 		TransactionID: transreq.ID,
// 		Status:        models.StatusFailed,
// 		Description:   transreq.Description,
// 		Amount:        transreq.Amount,
// 	}
// 	if err := db.DB.Create(&transhistory).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "error while creating record in database",
// 			Data: map[string]interface{}{
// 				"data": err.Error(),
// 			},
// 		})
// 		return
// 	}

// }
