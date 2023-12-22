// controllers/auth_controller.go
package controllers

import (
	"net/http"

	"github.com/R3PTR/go-auth-api/auth"
	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related requests
type AuthController struct {
	authService *auth.AuthService
}

func NewAuthController(authService *auth.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login handles user login and issues a JWT
func (ac *AuthController) Login(c *gin.Context) {
	var loginRequest auth.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	login, error := ac.authService.Login(loginRequest.Username, loginRequest.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": login})
}

// CreateUser handles user creation
func (ac *AuthController) CreateUser(c *gin.Context) {
	var createUserRequest auth.CreateUserRequest
	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.authService.CreateUser(createUserRequest.Username, createUserRequest.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Activates User
func (ac *AuthController) ActivateUser(c *gin.Context) {
	var activeUserRequest auth.ActivateUserRequest
	if err := c.ShouldBindJSON(&activeUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, err := user_unasserted.(auth.User)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	error := ac.authService.ActivateUser(user, activeUserRequest.OneTimePassword, activeUserRequest.NewPassword)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{"message": "User activated successfully"})
}

// Reset Password
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var resetPasswordRequest auth.ResetPasswordRequest
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, err := user_unasserted.(auth.User)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	error := ac.authService.ResetPassword(user.Username, resetPasswordRequest.OneTimePassword, resetPasswordRequest.NewPassword)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{"message": "Password reset successfully"})
}

// Change Password
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var changePasswordRequest auth.ChangePasswordRequest
	if err := c.ShouldBindJSON(&changePasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, err := user_unasserted.(auth.User)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	error := ac.authService.ChangePassword(user.Username, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{"message": "Password changed successfully"})
}

// Send Reset Password
func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var forgotPasswordRequest auth.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&forgotPasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	error := ac.authService.ForgotPassword(forgotPasswordRequest.Username)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{"message": "Reset password email sent successfully"})
}

// Logout handles user logout
func (ac *AuthController) Logout(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, err := user_unasserted.(auth.User)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	error := ac.authService.Logout(user.Username)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// DeleteOwnUser handles user deletion of own user
func (ac *AuthController) DeleteOwnUser(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, err := user_unasserted.(auth.User)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	error := ac.authService.DeleteOwnUser(user.Username)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// DeleteOtherUser handles user deletion of other user
func (ac *AuthController) DeleteOtherUser(c *gin.Context) {
	var deleteUserRequest auth.DeleteOtherUserRequest
	if err := c.ShouldBindJSON(&deleteUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.authService.DeleteOtherUser(deleteUserRequest.UsernameToDelete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
