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
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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

func (a *AuthService) CreateUser(createUserRequest CreateUserRequest) error {
	//Check if user exists
	existingUser, _ := a.AuthDbService.GetUserbyUsername(createUserRequest.Username)
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
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println(err)
		return errors.New("something went wrong hashing the password")
	}
	if err != nil {
		log.Fatal(err)
		return err
	}
	timestamp := time.Now()
	user := User{
		Username:            createUserRequest.Username,
		FirstName:           createUserRequest.FirstName,
		LastName:            createUserRequest.LastName,
		Password:            "",
		OneTimePassword:     hashedPassword,
		Role:                createUserRequest.Role,
		Personnelnumber:     createUserRequest.Personnelnumber,
		VacationDaysPerYear: createUserRequest.VacationDaysPerYear,
		TargetHoursPerWeek:  createUserRequest.TargetHoursPerWeek,
		MaximumHoursPerWeek: createUserRequest.MaximumHoursPerWeek,
		State:               NEW,
		InsertedAt:          timestamp,
		UpdatedAt:           timestamp,
	}
	err = a.AuthDbService.CreateUser(user)
	if err != nil {
		return errors.New("something went wrong creating the user")
	}
	a.EmailSender.SendEmail(user.Username, "New User", "Your new password is: "+password)
	return nil
}

func (a *AuthService) DeleteOwnUser(userId string) error {
	err := a.AuthDbService.DeleteUserById(userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) DeleteOtherUser(userId string) error {
	err := a.AuthDbService.DeleteUserById(userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) ActivateUser(user *User, newPassword string) error {
	if user.State != NEW {
		return errors.New("User is already activated")
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
		return errors.New("something went wrong hashing the password")
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

func (a *AuthService) ChangePassword(username, newPassword string) error {
	// Check if user exists
	user, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return errors.New("username or Password incorrect")
	}
	if user.State != ACTIVE {
		return errors.New("User is not active")
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
	a.AuthDbService.DeleteTokensByUserId(user.Id)
	return nil
}

func (a *AuthService) ResetPassword(user *User, newPassword string) error {
	if user.State != ACTIVE {
		return errors.New("User is not active")
	}
	// Check if resetValidUntil is still valid
	if user.ResetValidUntil.Before(time.Now()) {
		return errors.New("reset password is not valid anymore")
	}
	// Generate new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Update user
	filter := bson.M{"username": user.Username}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "oneTimePassword": nil, "resetValidUntil": nil, "updatedAt": time.Now()}}
	_, err = a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	a.AuthDbService.DeleteTokensByUserId(user.Username)
	return nil
}

func (a *AuthService) ForgotPassword(username string) error {
	// Check if user exists
	_, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return errors.New("username or Password incorrect")
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
	err = a.EmailSender.SendEmail(username, "Forgot Password", "Your new password is: "+password)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) Login(username, password, totp string) (*tokenModel, bool, error) {
	// Check if user exists
	user, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return nil, false, errors.New("username or Password incorrect")
	}
	token_type := "LoginToken"
	expires := time.Now().Add(time.Hour * 720)
	if user.State != ACTIVE {
		token_type = "ActivationToken"
		expires = time.Now().Add(time.Minute * 15)
		err := bcrypt.CompareHashAndPassword([]byte(user.OneTimePassword), []byte(password))
		if err != nil {
			return nil, false, errors.New("username or Password incorrect")
		}
	} else {
		// Check if password is correct
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			err = bcrypt.CompareHashAndPassword([]byte(user.OneTimePassword), []byte(password))
			if err == nil {
				token_type = "ResetToken"
				expires = time.Now().Add(time.Minute * 15)
			} else {
				return nil, false, errors.New("username or Password incorrect")
			}
		}
		//Check if TOTP is active
		if user.TotpActive {
			if totp == "" {
				return nil, true, errors.New("no TOTP provided")
			}
			// Check if TOTP is correct
			valid, err := a.VerifyTOTP(user, totp)
			if err != nil {
				return nil, true, err
			}
			if !valid {
				return nil, true, errors.New("TOTP is not valid")
			}
		}
	}
	// Generate JWT
	token_string, err := a.generateJWTToken(username, user.Role, expires.Unix())
	if err != nil {
		return nil, false, err
	}
	// Write token to database
	token, err := a.AuthDbService.WriteTokenToDatabase(user.Id, token_string, token_type, expires)
	if err != nil {
		return nil, false, err
	}
	return token, false, nil
}

func (a *AuthService) Logout(username, token string) error {
	// Check if user exists
	user, error := a.AuthDbService.GetUserbyUsername(username)
	if error != nil {
		return errors.New("username or Password incorrect")
	}
	if user.State != ACTIVE {
		// User is not active
		return errors.New("User is not active")
	}
	// Delete token from database
	err := a.AuthDbService.DeleteTokenByUserId(user.Id, token)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) generateJWTToken(username, role string, expires int64) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.New(jwt.SigningMethodHS512)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = expires
	claims["role"] = role

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
	const charset = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"
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

func (a *AuthService) GetTOTP(user *User) (*otp.Key, []string, error) {
	issuer := a.config.TOTPIssuer
	key, err := generateTOTPKey(user.Username, issuer)
	if err != nil {
		return nil, nil, err
	}
	backupCodes, err := a.GenerateBackupCodes(user)
	if err != nil {
		return nil, nil, err
	}
	user.TotpSecret = key.Secret()
	user.TotpActive = false
	err = a.AuthDbService.UpdateUser(user)
	if err != nil {
		return nil, nil, err
	}
	return key, backupCodes, nil
}

func (a *AuthService) GenerateBackupCodes(user *User) ([]string, error) {
	var codes []string
	var hashed_codes []string
	for i := 0; i < 8; i++ {
		code, err := generateRandomPassword(8)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
		hashed_code, err := HashPassword(code)
		if err != nil {
			return nil, err
		}
		hashed_codes = append(hashed_codes, hashed_code)
	}
	user.BackupCodes = hashed_codes
	a.AuthDbService.UpdateUser(user)
	return codes, nil
}

func (a *AuthService) checkBackupCodes(user *User, code string) bool {
	for _, c := range user.BackupCodes {
		err := bcrypt.CompareHashAndPassword([]byte(c), []byte(code))
		if err == nil {
			// delete code from backup codes
			var newCodes []string
			for _, c := range user.BackupCodes {
				hashed_code, err := HashPassword(code)
				if err != nil {
					return false
				}
				if c != hashed_code {
					newCodes = append(newCodes, c)
				}
			}
			user.BackupCodes = newCodes
			authDbService := a.AuthDbService
			authDbService.UpdateUser(user)
			return true
		}
	}
	return false
}

func (a *AuthService) ActivateTOTP(user *User, otp string) error {
	valid := totp.Validate(otp, user.TotpSecret)
	if valid {
		user.TotpActive = true
		err := a.AuthDbService.UpdateUser(user)
		if err != nil {
			return err
		}
		a.AuthDbService.DeleteTokensByUserId(user.Username)
		return nil
	}
	return errors.New("OTP is not valid")
}

func (a *AuthService) DeactivateTOTP(user *User) error {
	user.TotpActive = false
	user.BackupCodes = nil
	user.TotpSecret = ""
	err := a.AuthDbService.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func generateTOTPKey(accountName, issuer string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		log.Fatal(err)
	}
	return key, nil
}

func (a *AuthService) VerifyTOTP(user *User, otp string) (bool, error) {
	if user.State != ACTIVE {
		// User is not active
		return false, errors.New("User is not active")
	}
	valid := totp.Validate(otp, user.TotpSecret)
	if valid {
		return true, nil
	}
	valid = a.checkBackupCodes(user, otp)
	return valid, nil
}

// Get All Users
func (a *AuthService) GetAllUsers() ([]UserOutputAll, error) {
	users, err := a.AuthDbService.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Get User By Username
func (a *AuthService) GetUserByUsername(username string) (*User, error) {
	user, err := a.AuthDbService.GetUserbyUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update User
func (a *AuthService) UpdateUser(userId, username, firstName, lastName, role, personnelnumber string, vacationDaysPerYear int, targetHoursPerWeek, maximumHoursPerWeek float32) error {
	user, err := a.AuthDbService.GetUserbyId(userId)
	if err != nil {
		return err
	}
	if username != "" {
		user.Username = username
	}
	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}
	if role != "" {
		user.Role = role
	}
	if personnelnumber != "" {
		user.Personnelnumber = personnelnumber
	}
	if vacationDaysPerYear != 0 {
		user.VacationDaysPerYear = vacationDaysPerYear
	}
	if targetHoursPerWeek != 0 {
		user.TargetHoursPerWeek = targetHoursPerWeek
	}
	if maximumHoursPerWeek != 0 {
		user.MaximumHoursPerWeek = maximumHoursPerWeek
	}
	err = a.AuthDbService.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

// VerifyPassword
func (a *AuthService) VerifyPassword(user *User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("username or Password incorrect")
	}
	return nil
}

// GetOwnUser
func (a *AuthService) GetOwnUser(userId string) (*UserOutput, error) {
	user, err := a.AuthDbService.GetOwnUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
