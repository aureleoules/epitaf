package db

import (
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
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
	DB, err = sqlx.Connect("mysql", user+":"+pass+"@("+host+")/"+database+"?charset=utf8mb4,utf8&parseTime=true")
	if err != nil {
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
	DB.Close()
	zap.S().Info("Closed DB")
}

func Delete() {
	if DB == nil {
		return
	}

	zap.S().Info("Deleting epitaf...")
	_, err := DB.Exec("DROP DATABASE epitaf;")
	if err != nil {
		zap.S().Fatal(err)
	}
}
