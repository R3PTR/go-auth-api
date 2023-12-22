// database/mongodb.go
package database

import (
	"context"
	"fmt"
	"log"

	"github.com/R3PTR/go-auth-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient is a struct that holds the MongoDB client.
type MongoDBClient struct {
	client *mongo.Client
	Config *config.Config
}

// NewMongoDBClient creates a new MongoDB client.
func NewMongoDBClient(config *config.Config) (*MongoDBClient, error) {
	clientOptions := options.Client().ApplyURI(config.MongoDBURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return &MongoDBClient{client: client, Config: config}, nil
}

// Close closes the MongoDB client connection.
func (mc *MongoDBClient) Close() {
	err := mc.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection closed.")
}

// GetCollection returns a MongoDB collection.
func (mc *MongoDBClient) GetCollection(database, collection string) *mongo.Collection {
	return mc.client.Database(database).Collection(collection)
}
