package handlers

import (
	"gotransact/responses"

	logger "gotransact/apps/Accounts"
	utils "gotransact/apps/Accounts/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @BasePath /api
// LogoutHandler handles user logout.
// @Summary Logout a user
// @Description Logs out a user by invalidating the provided token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header with bearer token"
// @Success 200 {object} responses.UserResponse "Logged out successfully"
// @Failure 401 {object} responses.UserResponse "Unauthorized"
// @Failure 500 {object} responses.UserResponse "Internal Server Error"
// @Router /Logout [post]
func LogoutHandler(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("attempted Logout")
	authHeader := c.GetHeader("Authorization")

	status, message, data := utils.Logout(authHeader)

	c.JSON(status, responses.UserResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// if authHeader == "" {
// 	c.JSON(http.StatusUnauthorized, responses.UserResponse{
// 		Status:  http.StatusUnauthorized,
// 		Message: "authorization header missing",
// 	})
// 	return
// }

// //tokenStr := authHeader[len("Bearer "):]
// status, message, data := utils.Logout(authHeader)
// _, err := utils.VerifyPasetoToken(authHeader)
// if err != nil {
// 	c.JSON(http.StatusUnauthorized, responses.UserResponse{
// 		Status:  http.StatusUnauthorized,
// 		Message: "invalid token",
// 		Data: map[string]interface{}{
// 			"data": err.Error(),
// 		},
// 	})
// 	return
// }

// // Blacklist the token by storing it in Redis with an expiration time
// err = rdb.Set(authHeader, "Blacklisted", 24*time.Hour).Err() // adjust expiration time as needed
// if err != nil {
// 	c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 		Status:  http.StatusInternalServerError,
// 		Message: "failed to blacklist token",
// 		Data: map[string]interface{}{
// 			"data": err.Error(),
// 		},
// 	})
// 	return
// }
