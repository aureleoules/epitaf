package models

import (
	"github.com/aureleoules/epitaf/db"
)

const (
	userSchema = `
		CREATE TABLE users (
			uuid BINARY(16) NOT NULL,
			
			name VARCHAR(256) NOT NULL,
			email VARCHAR(256) NOT NULL UNIQUE,
			promotion VARCHAR(256) NOT NULL,
			class VARCHAR(256) NOT NULL,

			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (uuid),
			UNIQUE INDEX (email)
		);
	`

	insertUserQuery = `
		INSERT INTO users 
			(uuid, name, email, promotion, class) 
		VALUES 
			(:uuid, :name, :email, :promotion, :class);
	`

	getUserByEmailQuery = `
		SELECT 
			uuid, 
			name, 
			email, 
			promotion,
			class,
			created_at,
			updated_at
		FROM users
		WHERE email = ?;
	`

	getUserQuery = `
		SELECT 
			uuid, 
			name, 
			email, 
			promotion,
			class,
			created_at,
			updated_at
		FROM users
		WHERE uuid = ?;
	`
)

// User struct
type User struct {
	base

	Name      string `json:"name" db:"name"`
	Promotion int    `json:"promotion" db:"promotion"`
	Class     string `json:"class" db:"class"`
	Email     string `json:"email" db:"email"`
}

// MicrosoftProfile struct
type MicrosoftProfile struct {
	DisplayName string `json:"displayName"`
	Mail        string `json:"mail"`
}

// GetUserByEmail retrives user by email
func GetUserByEmail(email string) (*User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var user User
	err = tx.Get(&user, getUserByEmailQuery, email)
	return &user, err
}

// GetUser retrives user by uuid
func GetUser(uuid UUID) (*User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var user User
	err = tx.Get(&user, getUserQuery, uuid)
	return &user, err
}

// Insert user in DB
func (c *User) Insert() error {
	c.UUID = NewUUID()

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.NamedExec(insertUserQuery, c)
	return err
}
