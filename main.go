// main.go
package main

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

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
	router.Use(cors.Default())
	testRouter := router.Group("/test").Use(authMiddleware.AuthMiddleware([]string{"ADMIN"}))
	testRouter.GET("/Test", func(c *gin.Context) {
		err := bcrypt.CompareHashAndPassword([]byte("$2a$10$rEieyUa4kUIl7CxRszYywuj7SoWXUqw9vofzZ66B1mKM8p3qEdYSS"), []byte("123456"))
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Password is correct"})
	})
	// Register the routes
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/createUser", authMiddleware.AuthMiddleware([]string{"ADMIN"}), authController.CreateUser)
		authRouter.POST("/activateUser", authController.ActivateUser)
		authRouter.POST("/resetPassword", authController.ResetPassword)
		authRouter.POST("/forgotPassword", authController.ForgotPassword)
		authRouter.POST("/changePassword", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}), authMiddleware.PasswordMiddleware(), authController.ChangePassword)
		authRouter.POST("/changeOwnUsername", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}), authMiddleware.PasswordMiddleware(), authController.ChangeOwnUsername)
		authRouter.POST("/changeOtherUsername", authMiddleware.AuthMiddleware([]string{"ADMIN"}), authMiddleware.PasswordMiddleware(), authController.ChangeOtherUsername)
		authRouter.GET("/getOwnUser", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}), authController.GetOwnUser)
		authRouter.GET("/getAllUsers", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER"}), authController.GetAllUsers)
		authRouter.POST("/logout", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}), authController.Logout)
		authRouter.POST("/getTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}), authController.GetTOTP)
		authRouter.POST("/activateTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}), authController.ActivateTOTP)
		authRouter.POST("/deactivateTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}), authMiddleware.TOTPMiddleware(), authController.DeactivateTOTP)
		authRouter.POST("regenerateBackupCodes", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}), authMiddleware.TOTPMiddleware(), authController.RegenerateBackupCodes)
	}
	router.Run(":9090")
}
