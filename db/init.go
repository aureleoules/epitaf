package db

import (
	"os"

	"go.uber.org/zap"
)

// Init DB
func Init() {
	if DB == nil {
		connect(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), "")
	}

	// Create database
	zap.S().Info("Creating database...")
	_, err := DB.Exec(`CREATE DATABASE ` + os.Getenv("DB_NAME"))
	if err != nil {
		zap.S().Warn(err)
		// return
	}

	connect(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	// Timezone
	zap.S().Info("Settings Europe/Paris as timezone...")
	_, err = DB.Exec(`SET GLOBAL time_zone = 'Europe/Paris';`)
	if err != nil {
		zap.S().Fatal(err)
	}
}
