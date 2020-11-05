package models

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func checkErr(tx *sqlx.Tx, err error) {
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			zap.S().Error(err)
		}
	} else {
		err = tx.Commit()
		if err != nil {
			zap.S().Error(err)
		}
	}
}
