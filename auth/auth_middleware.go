// middleware.go
package auth

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/R3PTR/go-auth-api/database"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	mongoClient   *database.MongoDBClient
	AuthDbService *AuthDbService
	AuthService   *AuthService
}

func NewMiddleware(mongoClient *database.MongoDBClient, authDbService *AuthDbService, authService *AuthService) *AuthMiddleware {
	return &AuthMiddleware{mongoClient: mongoClient, AuthDbService: authDbService, AuthService: authService}
}

func (AuthMiddleware *AuthMiddleware) AuthMiddleware(roleRequired []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		jwt_token, err := ExtractToken(header)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(jwt_token)
		token_model, err := AuthMiddleware.AuthDbService.GetTokenByToken(jwt_token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		if token_model.Expires.Before(time.Now()) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token expired"})
			return
		}
		user, err := AuthMiddleware.AuthDbService.GetUserbyId(token_model.UserId)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		if !slices.Contains(roleRequired, user.Role) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("user", user)
		c.Set("token", token_model)
		c.Next()
	}
}

// Middleware for Methods that require Password Check
func (AuthMiddleware *AuthMiddleware) PasswordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_unasserted, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "User not found"})
			return
		}
		user := user_unasserted.(*User)
		password := c.GetHeader("Password")
		if password == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Password not provided"})
			return
		}
		err := AuthMiddleware.AuthService.VerifyPassword(user, password)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Next()
	}
}

func (AuthMiddleware *AuthMiddleware) TOTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_unasserted, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "User not found"})
			return
		}
		user := user_unasserted.(*User)
		if !user.TotpActive {
			c.AbortWithStatusJSON(401, gin.H{"error": "TOTP not activated"})
			return
		}
		totp := c.GetHeader("TOTP")
		if totp == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "TOTP not provided"})
			return
		}
		correct, error := AuthMiddleware.AuthService.VerifyTOTP(user, totp)
		if error != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": error.Error()})
			return
		}
		if !correct {
			c.AbortWithStatusJSON(401, gin.H{"error": "TOTP not valid"})
			return
		}
		c.Next()
	}
}

func ExtractToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	token := strings.Split(header, " ")
	if len(token) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return token[1], nil
}
