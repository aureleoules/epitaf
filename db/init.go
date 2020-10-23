package db

import (
	"os"

	"go.uber.org/zap"
)

func Init() {
	if DB == nil {
		connect(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), "")
	}

	// Create database
	zap.S().Info("Creating database...")
	_, err := DB.Exec(`CREATE DATABASE ` + os.Getenv("DB_NAME"))
	if err != nil {
		zap.S().Fatal(err)
	}
	zap.S().Info("Created database: " + os.Getenv("DB_NAME"))

	connect(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	// Timezone
	zap.S().Info("Settings UTC as timezone...")
	_, err = DB.Exec(`SET GLOBAL time_zone = 'UTC';`)
	if err != nil {
		zap.S().Fatal(err)
	}
}
