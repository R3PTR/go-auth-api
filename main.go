// main.go
package main

import (
	"fmt"

	"github.com/R3PTR/go-auth-api/auth"
	"github.com/R3PTR/go-auth-api/config"
	"github.com/R3PTR/go-auth-api/database"
	"github.com/R3PTR/go-auth-api/emails"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the config package
	config, err := config.NewConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	// Create MongoDB client
	mongoClient, err := database.NewMongoDBClient(config)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	defer mongoClient.Close()
	//
	emailSender := emails.NewEmailSender("ems@te-autoteile.de", "localhost", 1025, "", "")
	// AuthDbService
	authDbService := auth.NewAuthDbService(mongoClient)
	authService := auth.NewAuthService(mongoClient, config, authDbService, emailSender)
	authController := auth.NewAuthController(authService)

	// AuthMiddleware
	authMiddleware := auth.NewMiddleware(mongoClient, authDbService, authService)
	// Create a new router
	router := gin.Default()

	// Cors Config
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"http://localhost:9000", "http://192.168.50.232:9000"}
	cors_config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Password", "Allow-Origin"}
	cors_config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	cors_config.AllowCredentials = true
	router.Use(cors.New(cors_config))

	// Register the routes
	authRouter := router.Group("/auth")
	{
		// GET Routes
		authRouter.GET("/getOwnUser", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authController.GetOwnUser)
		authRouter.GET("/getAllUsers", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER"}, []string{"LoginToken"}), authController.GetAllUsers)
		// POST Routes
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/logout", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authController.Logout)
		authRouter.POST("/createUser", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), authController.CreateUser)
		authRouter.POST("/deleteOtherUser", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), authController.DeleteOtherUser)
		authRouter.POST("/activateUser", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"ActivationToken"}), authController.ActivateUser)
		authRouter.POST("/resetPassword", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"ResetToken"}), authController.ResetPassword)
		authRouter.POST("/forgotPassword", authController.ForgotPassword)
		authRouter.POST("/changePassword", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authMiddleware.PasswordMiddleware(), authController.ChangePassword)
		authRouter.POST("/updateOwnUser", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authController.UpdateOwnUser)
		authRouter.POST("/updateOtherUser", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), authController.UpdateOtherUser)
		authRouter.POST("/getTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authController.GetTOTP)
		authRouter.POST("/activateTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authController.ActivateTOTP)
		authRouter.POST("/deactivateTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authMiddleware.TOTPMiddleware(), authController.DeactivateTOTP)
		authRouter.POST("regenerateBackupCodes", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authMiddleware.TOTPMiddleware(), authController.RegenerateBackupCodes)
	}
	router.Run(":9090")
}
