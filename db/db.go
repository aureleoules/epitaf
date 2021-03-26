package db

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	// Import postgres driver
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	// DB holds DB socket
	DB    *sqlx.DB
	tries = 0
)

// Open DB
func Connect() {
	connect(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
}

func connect(host, port, user, pass, database string) {
	if tries > 3 {
		zap.S().Fatal("could not connect to database.")
	}
	zap.S().Info("Connecting to database...")

	connInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, database)
	zap.S().Info(connInfo)

	var err error
	DB, err = sqlx.Connect("postgres", connInfo)
	if err != nil {
		zap.S().Error(err)

		time.Sleep(5 * time.Second)
		connect(host, port, user, pass, database)
	}

	zap.S().Info("Connected to database.")
}

// Close DB
func Close() {
	if DB == nil {
		return
	}
	zap.S().Info("Closing DB...")
	err := DB.Close()
	if err != nil {
		zap.S().Error()
		return
	}
	zap.S().Info("Closed DB")
}

// Delete DB
func Delete() {
	if DB == nil {
		return
	}

	zap.S().Info("Deleting...")
	_, err := DB.Exec("DROP DATABASE IF EXISTS " + os.Getenv("DB_NAME"))
	if err != nil {
		zap.S().Fatal(err)
	}
}
