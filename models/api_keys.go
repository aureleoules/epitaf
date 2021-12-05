package models

import (
	"github.com/aureleoules/epitaf/db"
	"go.uber.org/zap"
)

const (
	apiKeySchema = `
		CREATE TABLE api_keys (
			api_key VARCHAR(64) NOT NULL UNIQUE,

			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),

			PRIMARY KEY (api_key)
		);
	`

	getAPIKey = `
		SELECT 
			api_key
		FROM api_keys
		WHERE api_key = ?;
	`

	insertAPIKey = `
		INSERT INTO api_keys
			(api_key)
		VALUES
			(?);
	`
)

// InsertAPIKey inserts API key in DB
func InsertAPIKey(apiKey string) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.Exec(insertAPIKey, apiKey)
	return err
}

// IsAPIKeyCorrect checks if the API key is registered
func IsAPIKeyCorrect(apiKey string) bool {
	zap.S().Info("Retrieving api key...")
	tx, err := db.DB.Beginx()
	if err != nil {
		return false
	}

	defer checkErr(tx, err)

	var exists *[]uint8

	err = tx.Get(&exists, getAPIKey, apiKey)
	if err != nil {
		zap.S().Error(err)
		zap.S().Info("Api key '", apiKey, "' does not exists")
		return false
	}
	zap.S().Info("Retrieved api key.")
	return true
}
