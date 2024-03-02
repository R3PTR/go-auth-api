// config/config.go
package config

import (
	"encoding/json"
	"os"
)

// Config represents the MongoDB configuration.
type Config struct {
	MongoDBURI          string `json:"mongodb_uri"`
	UserDatabase        string `json:"user_database"`
	UserCollection      string `json:"user_collection"`
	SiteDatabase        string `json:"site_database"`
	SiteCollection      string `json:"site_collection"`
	WorkspaceCollection string `json:"workspace_collection"`
	AbsenceCollection   string `json:"absence_collection"`
	AbsencesDatabase    string `json:"absences_database"`
	JWTSecret           string `json:"jwt_secret"`
	TokenCollection     string `json:"token_collection"`
	TOTPIssuer          string `json:"totp_issuer"`
}

// Environment represents the environment (development, production, etc.).
type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

var CurrentEnvironment Environment

func NewConfig() (*Config, error) {
	// Set the environment during package initialization
	CurrentEnvironment = Environment(os.Getenv("APP_ENV"))
	if CurrentEnvironment == "" {
		CurrentEnvironment = Development
	}
	return loadConfig()
}

func loadConfig() (*Config, error) {
	configFileName := "config_dev.json"
	if CurrentEnvironment == Production {
		configFileName = "config_prod.json"
	}

	file, err := os.Open("./configs/" + configFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
