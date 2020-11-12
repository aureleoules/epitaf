package models

import "github.com/aureleoules/epitaf/db"

// Class struct
type Class struct {
	Promotion int    `json:"promotion" db:"promotion"`
	Semester  string `json:"semester" db:"semester"`
	Region    string `json:"region" db:"region"`
	Class     string `json:"class" db:"class"`
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
