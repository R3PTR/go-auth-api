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
}

func NewMiddleware(mongoClient *database.MongoDBClient, authDbService *AuthDbService) *AuthMiddleware {
	return &AuthMiddleware{mongoClient: mongoClient, AuthDbService: authDbService}
}

func (AuthMiddleware *AuthMiddleware) AuthMiddleware(roleRequired []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		jwt_token, err := extractBearerToken(header)
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
		c.Next()
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
