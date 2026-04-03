package routes

import (
	"go-auth-app/config"
	"go-auth-app/controllers"
	"go-auth-app/middleware"
	"go-auth-app/repositories"
	"go-auth-app/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	userRepo := &repositories.UserRepoDB{}
	refreshTokenRepo := repositories.NewRefreshTokenRepo()
	emailVerificationRepo := repositories.NewEmailVerificationTokenRepo()
	socialRepo := repositories.NewSocialAccountRepository()
	jwtService := services.NewJWTService(config.JWTSecret)
	emailSvc := services.NewEmailService()

	authService := services.NewAuthService(
		userRepo,
		refreshTokenRepo,
		emailVerificationRepo,
		socialRepo,
		jwtService,
		emailSvc,
	)

	userService := &services.UserService{
		UserRepo: userRepo,
	}

	authController := controllers.AuthController{
		AuthService: authService,
	}

	userController := controllers.UserController{
		UserService: userService,
	}

	auth := router.Group("/auth")

	// ========================
	// 🔓 PUBLIC ROUTES
	// ========================
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
	auth.POST("/refresh", authController.RefreshToken)
	auth.GET("/verify", authController.VerifyEmail)
	auth.GET("/resend-verification", authController.ResendVerification)
	auth.POST("/forgot-password", authController.ForgotPassword)
	auth.POST("/reset-password", authController.ResetPassword)
	auth.POST("/social-login", authController.SocialLogin)

	// ========================
	// 🔐 PROTECTED ROUTES
	// ========================
	protected := auth.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService))
	protected.GET("/profile", authController.Profile)
	protected.POST("/logout", authController.Logout)

	// ========================
	// 👑 ADMIN ROUTES
	// ========================
	admin := protected.Group("/admin")
	admin.Use(middleware.AdminOnly())
	admin.GET("/users", userController.GetAllUsers)
	admin.DELETE("/users/:id", userController.DeleteUser)
}
