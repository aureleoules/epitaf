package models

import (
	"github.com/aureleoules/epitaf/db"
	"github.com/mattn/go-nulltype"
)

// Class struct
type Class struct {
	Promotion nulltype.NullInt64  `json:"promotion" db:"promotion"`
	Semester  nulltype.NullString `json:"semester" db:"semester"`
	Region    nulltype.NullString `json:"region" db:"region"`
	Class     nulltype.NullString `json:"class" db:"class"`
}

const (
	getClassesQuery = `
		SELECT DISTINCT
			promotion, 
			semester, 
			region, 
			class 
		FROM users
		WHERE teacher=0;
	`
)

// GetClasses fetches distinct classes from users in DB
func GetClasses() ([]Class, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var classes []Class
	err = tx.Select(&classes, getClassesQuery)
	return classes, err
}
