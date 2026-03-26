package routes

import (
	"go-auth-app/middleware"
	"go-auth-app/repositories"
	"go-auth-app/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	repo := &repositories.UserRepoDB{}
	authController := controllers.AuthController{
		UserRepo: repo,
	}
	api := router.Group("/api")

	// Public
	api.POST("/register", authController.Register)
	api.POST("/login", authController.Login)
	api.POST("/refresh", middleware.RefreshToken)

	// Protected
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/profile", authController.Profile)
	protected.POST("/logout", middleware.AuthMiddleware(), authController.Logout)

	// Admin only
	admin := protected.Group("/admin")
	admin.Use(middleware.AdminOnly())

		admin.GET("/dashboard", authController.Dashboard)
}