package models

import (
	"github.com/aureleoules/epitaf/db"
	"go.uber.org/zap"
)

const (
	apiKeySchema = `
		CREATE TABLE api_keys (
			token VARCHAR(64) NOT NULL UNIQUE,
			label VARCHAR(64) NOT NULL UNIQUE,

			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),

			PRIMARY KEY (token)
		);
	`

	getAPIKey = `
		SELECT 
			token
		FROM api_keys
		WHERE token = ?;
	`

	insertAPIKey = `
		INSERT INTO api_keys
			(token, label)
		VALUES
			(:token, :label);
	`
)

// APIKey struct
type APIKey struct {
	base

	Token string
	Label string
}

// Insert inserts API key in DB
func (a *APIKey) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(insertAPIKey, a)
	return err
}

// IsAPIKeyCorrect checks if the API key is registered
func IsAPIKeyCorrect(token string) bool {
	zap.S().Info("Retrieving api key...")
	tx, err := db.DB.Beginx()
	if err != nil {
		return false
	}

	defer checkErr(tx, err)

	var exists *[]uint8

	err = tx.Get(&exists, getAPIKey, token)
	if err != nil {
		zap.S().Error(err)
		zap.S().Info("Api key '", token, "' does not exists")
		return false
	}
	zap.S().Info("Retrieved api key.")
	return true
}
