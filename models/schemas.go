package models

import (
	"github.com/aureleoules/epitaf/db"
	"go.uber.org/zap"
)

var schemas = []string{userSchema, taskSchema}

// InjectSQLSchemas injects sql schemas in db
func InjectSQLSchemas() error {
	// Schemas
	zap.S().Info("Inserting ", len(schemas), " schemas...")
	for _, s := range schemas {
		zap.S().Info("Inserting schema...", s)
		_, err := db.DB.Exec(s)
		if err != nil {
			return err
		}
		zap.S().Info("Ok")
	}
	return nil
}
