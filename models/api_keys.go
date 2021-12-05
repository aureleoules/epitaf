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

	getApiKey = `
		SELECT 
			api_key
		FROM api_keys
		WHERE api_key = ?;
	`
)

// IsApiKeyCorrect checks if the api key is registered
func IsApiKeyCorrect(api_key string) bool {
	zap.S().Info("Retrieving api key...")
	tx, err := db.DB.Beginx()
	if err != nil {
		return false
	}

	defer checkErr(tx, err)

	var exists *[]uint8

	err = tx.Get(&exists, getApiKey, api_key)
	if err != nil {
		zap.S().Error(err)
		zap.S().Info("Api key '", api_key, "' does not exists")
		return false
	}
	zap.S().Info("Retrieved api key.")
	return true
}
