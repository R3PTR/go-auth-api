package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related requests
type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login handles user login and issues a JWT
func (ac *AuthController) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	header := c.GetHeader("TOTP")
	login, requires2FA, error := ac.authService.Login(loginRequest.Username, loginRequest.Password, header)
	if requires2FA {
		c.JSON(http.StatusOK, gin.H{"message": "Credentials correct", "requires_2fa": true})
		return
	}
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": login.Token, "token_type": login.TokenType})
}

// CreateUser handles user creation
func (ac *AuthController) CreateUser(c *gin.Context) {
	var createUserRequest CreateUserRequest
	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.authService.CreateUser(createUserRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Activates User
func (ac *AuthController) ActivateUser(c *gin.Context) {
	var activeUserRequest ActivateUserRequest
	if err := c.ShouldBindJSON(&activeUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, ok := user_unasserted.(*User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	err := ac.authService.ActivateUser(user, activeUserRequest.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{"message": "User activated successfully"})
}

// Reset Password
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var resetPasswordRequest ResetPasswordRequest
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, err := user_unasserted.(User)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	error := ac.authService.ResetPassword(&user, resetPasswordRequest.NewPassword)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{"message": "Password reset successfully"})
}

// Change Password
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var changePasswordRequest ChangePasswordRequest
	if err := c.ShouldBindJSON(&changePasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, ok := user_unasserted.(*User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	error := ac.authService.ChangePassword(user.Username, changePasswordRequest.NewPassword)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{"message": "Password changed successfully"})
}

// Send Reset Password
func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var forgotPasswordRequest ForgotPasswordRequest
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
	user, err := user_unasserted.(User)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	token_unasserted, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
		return
	}
	token, err := token_unasserted.(*tokenModel)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
		return
	}
	error := ac.authService.Logout(user.Username, token.Token)
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
	user, err := user_unasserted.(User)
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
	var deleteUserRequest DeleteOtherUserRequest
	if err := c.ShouldBindJSON(&deleteUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.authService.DeleteOtherUser(deleteUserRequest.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetTOTP
func (ac *AuthController) GetTOTP(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, ok := user_unasserted.(*User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	if user.TotpActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TOTP already activated"})
		return
	}
	otp, backupCodes, err := ac.authService.GetTOTP(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"otp_secret": otp.Secret(), "otp_url": otp.URL(), "backup_codes": backupCodes})
}

// ActivateTOTP
func (ac *AuthController) ActivateTOTP(c *gin.Context) {
	var activateTOTPRequest ActivateTOTPRequest
	if err := c.ShouldBindJSON(&activateTOTPRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, ok := user_unasserted.(*User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	err := ac.authService.ActivateTOTP(user, activateTOTPRequest.TOTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "TOTP activated successfully"})
}

// DeactivateTOTP
func (ac *AuthController) DeactivateTOTP(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, ok := user_unasserted.(*User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	err := ac.authService.DeactivateTOTP(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "TOTP deactivated successfully"})
}

// RegenerateBackupCodes
func (ac *AuthController) RegenerateBackupCodes(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, ok := user_unasserted.(*User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	backupCodes, err := ac.authService.GenerateBackupCodes(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"backup_codes": backupCodes})
}

// Get OwnUser
func (ac *AuthController) GetOwnUser(c *gin.Context) {
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, ok := user_unasserted.(*User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userOutput, err := ac.authService.GetOwnUser(user.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userOutput})
}

// Get All Users
func (ac *AuthController) GetAllUsers(c *gin.Context) {
	users, err := ac.authService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Update Own User
func (ac *AuthController) UpdateOwnUser(c *gin.Context) {
	var updateOwnUserRequest UpdateOwnUserRequest
	if err := c.ShouldBindJSON(&updateOwnUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_unasserted, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	user, ok := user_unasserted.(*User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	err := ac.authService.UpdateUser(user.Id, updateOwnUserRequest.Username, updateOwnUserRequest.FirstName, updateOwnUserRequest.LastName, "", "", 0, 0, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Update Other User
func (ac *AuthController) UpdateOtherUser(c *gin.Context) {
	var updateOtherUserRequest UpdateOtherUserRequest
	if err := c.ShouldBindJSON(&updateOtherUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.authService.UpdateUser(updateOtherUserRequest.Id, updateOtherUserRequest.Username, updateOtherUserRequest.FirstName, updateOtherUserRequest.LastName, updateOtherUserRequest.Role, updateOtherUserRequest.Personnelnumber, updateOtherUserRequest.VacationDaysPerYear, updateOtherUserRequest.TargetHoursPerWeek, updateOtherUserRequest.MaximumHoursPerWeek)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
