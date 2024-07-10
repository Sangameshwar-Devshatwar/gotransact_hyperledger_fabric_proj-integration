package handlers

import (
	logger "gotransact/apps/Accounts"
	utils "gotransact/apps/Accounts/utils"
	astructutils "gotransact/apps/Astructutils"

	"gotransact/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)
// @BasePath /api
// SignupHandler handles user signup.
// @Summary Signup a new user
// @Description Register a new user with email and company details
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body astructutils.SignupUser true "User Data"
// @Success 200 {object} responses.UserResponse "User Created Successfully"
// @Failure 400 {object} responses.UserResponse "Invalid Input"
// @Failure 500 {object} responses.UserResponse "Internal Server Error"
// @Router /signup [post]
func SignupHandler(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("attempted signup")
	var req astructutils.SignupUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	statuscode, message, data := utils.Signup(req)

	c.JSON(statuscode, responses.UserResponse{Status: statuscode,
		Message: message,
		Data:    data})
}
