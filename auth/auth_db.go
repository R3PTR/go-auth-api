package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/R3PTR/go-auth-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthDbService struct {
	mongoClient *database.MongoDBClient
}

func NewAuthDbService(mongoClient *database.MongoDBClient) *AuthDbService {
	return &AuthDbService{mongoClient: mongoClient}
}

// Create User
func (a *AuthDbService) CreateUser(user User) error {
	_, err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).InsertOne(context.Background(), user)
	return err
}

func (a *AuthDbService) GetUserbyUsername(username string) (*User, error) {
	user := &User{}
	err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).FindOne(context.Background(), bson.M{"username": username}).Decode(user)
	if err != nil {
		// Handle errors, e.g., user not found
		fmt.Println("Error:", err)
		return nil, err
	}
	return user, nil
}

// Get User by Id
func (a *AuthDbService) GetUserbyId(id string) (*User, error) {
	user := &User{}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	err = a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).FindOne(context.Background(), bson.M{"_id": objectId}).Decode(user)
	if err != nil {
		// Handle errors, e.g., user not found
		fmt.Println("Error:", err)
		return nil, err
	}
	return user, nil
}

func (a *AuthDbService) GetTokenByToken(token string) (*tokenModel, error) {
	token_model := &tokenModel{}
	err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.TokenCollection).FindOne(context.Background(), bson.M{"token": token}).Decode(token_model)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return token_model, nil
}

func (a *AuthDbService) WriteTokenToDatabase(userId, token string, expires time.Time) error {
	token_struct := tokenModel{
		UserId:     userId,
		Token:      token,
		InsertedAt: time.Now(),
		UpdatedAt:  time.Now(),
		Expires:    expires,
	}
	insertResult, err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.TokenCollection).InsertOne(context.Background(), token_struct)
	if err != nil {
		return err
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}

// Delete all tokens for a user
func (a *AuthDbService) DeleteTokensByUserId(userId string) error {
	_, err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.TokenCollection).DeleteMany(context.Background(), bson.M{"user_id": userId})
	if err != nil {
		return err
	}
	return nil
}

// Delete all tokens for a user
func (a *AuthDbService) DeleteTokenByUserId(userId, token string) error {
	_, err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.TokenCollection).DeleteOne(context.Background(), bson.M{"user_id": userId})
	if err != nil {
		return err
	}
	return nil
}

// Delete User
func (a *AuthDbService) DeleteUserByUsername(username string) error {
	_, err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).DeleteMany(context.Background(), bson.M{"username": username})
	if err != nil {
		return err
	}
	return nil
}

// Update User
func (a *AuthDbService) UpdateUser(user *User) error {
	// Delete Id from stuct, to prevent overwriting
	fmt.Println("Updating user")
	user_id := user.Id
	user.Id = ""
	objectId, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		log.Println("Invalid id")
	}
	filter := bson.M{"_id": objectId}
	result, err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).ReplaceOne(context.Background(), filter, user)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 0 {
		return nil
	}
	result, err = a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).UpdateOne(context.Background(), bson.M{"_id": user_id}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	if result.ModifiedCount != 0 {
		return nil
	}
	return fmt.Errorf("something went wrong, user not updated")
}

// Get All Users
func (a *AuthDbService) GetAllUsers() ([]UserOutputAll, error) {
	var users []UserOutputAll
	cursor, err := a.mongoClient.GetCollection(a.mongoClient.Config.UserDatabase, a.mongoClient.Config.UserCollection).Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user UserOutputAll
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
