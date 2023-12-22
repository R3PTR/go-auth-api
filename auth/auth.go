// auth.go
package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/R3PTR/go-auth-api/config"
	"github.com/R3PTR/go-auth-api/database"
	"github.com/R3PTR/go-auth-api/emails"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	mongoClient   *database.MongoDBClient
	config        *config.Config
	AuthDbService *AuthDbService
	EmailSender   *emails.EmailSender
}

const (
	NEW    = "NEW"
	ACTIVE = "ACTIVE"
)

const (
	ADMIN  = "ADMIN"
	USER   = "USER"
	DRIVER = "DRIVER"
)

// NewAuthService creates a new AuthService with the provided MongoDB client.
func NewAuthService(mongoClient *database.MongoDBClient, config *config.Config, authDbService *AuthDbService, emailSender *emails.EmailSender) *AuthService {
	return &AuthService{mongoClient: mongoClient, config: config, AuthDbService: authDbService, EmailSender: emailSender}
}

func (a *AuthService) CreateUser(username, role string) error {
	//Check if user exists
	existingUser, err := a.AuthDbService.GetUserbyUsername(username)
	if existingUser != nil {
		return errors.New("User already exists")
	}
	// Generate random password
	password, err := generateRandomPassword(8)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(password)
	hashedPassword, err := HashPassword(password)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println(err)
		return errors.New("Something went wrong hashing the password")
	}
	if err != nil {
		log.Fatal(err)
		return err
	}
	timestamp := time.Now()
	user := User{
		Username:        username,
		Password:        "",
		OneTimePassword: hashedPassword,
		Role:            role,
		State:           NEW,
		InsertedAt:      timestamp,
		UpdatedAt:       timestamp,
	}
	err = a.AuthDbService.CreateUser(user)
	if err != nil {
		return errors.New("Something went wrong creating the user")
	}
	// Send email
	// TODO implement email sending
	return nil
}

func (a *AuthService) DeleteOwnUser(username string) error {
	err := a.AuthDbService.DeleteUserByUsername(username)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) DeleteOtherUser(username string) error {
	err := a.AuthDbService.DeleteUserByUsername(username)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) ActivateUser(user User, OneTimePassword, newPassword string) error {
	if user.State != NEW {
		return errors.New("User is already activated")
	}
	// Check if password is correct
	err := bcrypt.CompareHashAndPassword([]byte(user.OneTimePassword), []byte(OneTimePassword))
	if err != nil {
		fmt.Println(err)
		return errors.New("Username or Password incorrect")
	}
	// Hash newPassword
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Double check if passwordhash is correct
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(newPassword))
	if err != nil {
		fmt.Println(err)
		return errors.New("Something went wrong hashing the password")
	}
	// Update user
	filter := bson.M{"username": user.Username}
	update := bson.M{"$set": bson.M{"state": ACTIVE, "password": hashedPassword, "oneTimePassword": nil, "resetValidUntil": nil, "updatedAt": time.Now()}}
	_, err = a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) ChangePassword(username, oldPassword, newPassword string) error {
	// Check if user exists
	user, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return errors.New("Username or Password incorrect")
	}
	if user.State != ACTIVE {
		return errors.New("User is not active")
	}
	// Check if password is correct
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		fmt.Println(err)
		return errors.New("Username and Password incorrect")
	}
	// Generate new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Update user
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "updatedAt": time.Now()}}
	_, err = a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	a.AuthDbService.DeleteTokenByUsername(username)
	return nil
}

func (a *AuthService) ResetPassword(username, oneTimePassword, newPassword string) error {
	// Check if user exists
	user, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return errors.New("Username or Password incorrect")
	}
	if user.State != ACTIVE {
		return errors.New("User is not active")
	}
	// Check if password is correct
	err := bcrypt.CompareHashAndPassword([]byte(user.OneTimePassword), []byte(oneTimePassword))
	if err != nil {
		fmt.Println(err)
		return errors.New("Username or Password incorrect")
	}
	// Check if resetValidUntil is still valid
	if user.ResetValidUntil.Before(time.Now()) {
		return errors.New("Reset password link is not valid anymore")
	}
	// Generate new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Update user
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "oneTimePassword": nil, "resetValidUntil": nil, "updatedAt": time.Now()}}
	_, err = a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	a.AuthDbService.DeleteTokenByUsername(username)
	return nil
}

func (a *AuthService) ForgotPassword(username string) error {
	// Check if user exists
	user, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return errors.New("Username or Password incorrect")
	}
	if user.State != ACTIVE {
		return errors.New("User is not active")
	}
	// Generate random password
	password, err := generateRandomPassword(8)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(password)
	// Set resetValidUntil
	resetValidUntil := time.Now().Add(time.Minute * 15)
	// Update user
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"oneTimePassword": hashedPassword, "resetValidUntil": resetValidUntil}}
	_, err = a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	a.EmailSender.SendEmail(username, "Forgot Password", "Your new password is: "+password)
	return nil
}

func (a *AuthService) Login(username, password string) (string, error) {
	// Check if user exists
	user, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return "", errors.New("Username or Password incorrect")
	}
	if user.State != ACTIVE {
		// User is not active
		return "", errors.New("User is not active")
	}
	// Check if password is correct
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Username or Password incorrect")
	}
	// Generate JWT
	expires := time.Now().Add(time.Hour * 720)
	token, err := a.generateJWTToken(username, expires.Unix())
	if err != nil {
		return "", err
	}
	// Write token to database
	err = a.AuthDbService.WriteTokenToDatabase(username, token, expires)
	if err != nil {
		return "", err
	}
	fmt.Printf("Found user: %+v\n", user)
	fmt.Printf("Token: %s\n", token)
	return token, nil
}

func (a *AuthService) Logout(username string) error {
	// Check if user exists
	user, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return errors.New("Username or Password incorrect")
	}
	if user.State != ACTIVE {
		// User is not active
		return errors.New("User is not active")
	}
	// Delete token from database
	err := a.AuthDbService.DeleteTokenByUsername(username)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) generateJWTToken(username string, expires int64) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.New(jwt.SigningMethodHS512)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = expires

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(a.config.JWTSecret))
}

func HashPassword(password string) (string, error) {
	// Generate a hashed representation of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password strings.Builder

	maxIndex := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return "", err
		}

		password.WriteByte(charset[randomIndex.Int64()])
	}

	return password.String(), nil
}
