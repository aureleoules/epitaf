package models

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-nulltype"
	"go.uber.org/zap"
)

var psql squirrel.StatementBuilderType

func init() {
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

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

func setValueDefaultNString(value, def nulltype.NullString) nulltype.NullString {
	if !value.Valid() {
		return def
	}
	return value
}
func setValueDefaultString(value, def string) string {
	if value == "" {
		return def
	}
	return value
}

func setValueDefaultNInt64(value, def nulltype.NullInt64) nulltype.NullInt64 {
	if !value.Valid() {
		return def
	}
	return value
}
