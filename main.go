// main.go
package main

import (
	"fmt"

	"github.com/R3PTR/go-auth-api/absences"
	"github.com/R3PTR/go-auth-api/auth"
	"github.com/R3PTR/go-auth-api/config"
	"github.com/R3PTR/go-auth-api/database"
	"github.com/R3PTR/go-auth-api/emails"
	"github.com/R3PTR/go-auth-api/sites"
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

	// Create SiteServices
	siteDbService := sites.NewSitesDbService(mongoClient)
	siteService := sites.NewSiteService(siteDbService)
	siteController := sites.NewSiteController(siteService)

	// Create AbsencesServices
	absencesDbService := absences.NewAbsencesDbService(mongoClient)
	absencesService := absences.NewAbsencesService(absencesDbService)
	absencesController := absences.NewAbsencesController(absencesService)

	router := gin.Default()

	// Cors Config
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"*"}
	cors_config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
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
		//authRouter.POST("/getTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authController.GetTOTP)
		//authRouter.POST("/activateTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authController.ActivateTOTP)
		//authRouter.POST("/deactivateTOTP", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authMiddleware.TOTPMiddleware(), authController.DeactivateTOTP)
		//authRouter.POST("regenerateBackupCodes", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), authMiddleware.TOTPMiddleware(), authController.RegenerateBackupCodes)
	}

	// Sites Routes
	siteRouter := router.Group("/sites")
	{
		// GET Routes
		siteRouter.GET("/getSites", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), siteController.GetSites)
		siteRouter.GET("/getWorkspaces", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), siteController.GetWorkspaces)
		// POST Routes
		siteRouter.POST("/createSite", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), siteController.CreateSite)
		siteRouter.POST("/createWorkspace", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), siteController.CreateWorkspace)
		// PUT Routes
		siteRouter.PUT("/updateSite", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), siteController.UpdateSite)
		siteRouter.PUT("/updateWorkspace", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), siteController.UpdateWorkspace)
		// DELETE Routes
		siteRouter.DELETE("/deleteSite/:id", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), siteController.DeleteSite)
		siteRouter.DELETE("/deleteWorkspace/:id", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), siteController.DeleteWorkspace)
	}
	// Absences Routes
	absencesRouter := router.Group("/absences")
	{
		// GET Routes
		absencesRouter.GET("/getAbsences", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), absencesController.GetAllAbsences)
		absencesRouter.GET("/getOwnAbsences", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), absencesController.GetAbsences)

		// POST Routes
		absencesRouter.POST("/createAbsence", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), absencesController.CreateAbsence)

		// PUT Routes
		absencesRouter.PUT("/updateOwnAbsence", authMiddleware.AuthMiddleware([]string{"ADMIN", "USER", "DRIVER"}, []string{"LoginToken"}), absencesController.UpdateOwnAbsence)
		absencesRouter.PUT("/updateAbsenceAsAdmin", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), absencesController.UpdateAbsenceAsAdmin)

		// DELETE Routes
		absencesRouter.DELETE("/deleteAbsence/:id", authMiddleware.AuthMiddleware([]string{"ADMIN"}, []string{"LoginToken"}), absencesController.DeleteAbsence)
	}

	router.Run(":9090")
}
