package db

import (
	"os"
	"time"

	// Import SQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	// DB holds DB socket
	DB    *sqlx.DB
	tries = 0
)

// Connect to DB
func Connect() {
	connect(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
}

func connect(host, user, pass, database string) {
	if tries > 3 {
		zap.S().Fatal("could not connect to database.")
	}
	zap.S().Info("Connecting to database...")
	var err error
	DB, err = sqlx.Connect("mysql", user+":"+pass+"@("+host+")/"+database+"?charset=utf8mb4,utf8&parseTime=true&time_zone=UTC")
	if err != nil {
		zap.S().Error(err)
		time.Sleep(5 * time.Second)
		connect(host, user, pass, database)
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

	zap.S().Info("Deleting epitaf...")
	_, err := DB.Exec("DROP DATABASE " + os.Getenv("DB_NAME"))
	if err != nil {
		zap.S().Fatal(err)
	}
}
