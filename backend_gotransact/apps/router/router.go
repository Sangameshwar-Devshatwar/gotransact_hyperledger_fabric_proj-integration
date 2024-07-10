package router

import (
	"gotransact/apps/Accounts/handlers"
	handler2 "gotransact/apps/transaction/handlers"
	"gotransact/docs"
	"gotransact/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routing() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"
	r.LoadHTMLGlob("/home/trellis/Desktop/folder-1/golang-training/backend_gotransact/apps/transaction/handlers/templates/*")

	r.POST("/api/signup", handlers.SignupHandler)
	r.POST("/api/signin", handlers.LoginHandler)
	r.GET("/api/confirm_payment", handler2.ConfirmationPayment)
	auth := r.Group("/api/protected")
	{
		auth.Use(middleware.AuthMiddleware())
		auth.POST("/postpayment", handler2.PostPayment)
		auth.POST("/Logout", handlers.LogoutHandler)

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8000")
}
