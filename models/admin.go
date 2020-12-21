package models

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/aureleoules/epitaf/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Admin struct
type Admin struct {
	base

	UUID     UUID   `json:"uuid" db:"uuid"`
	RealmID  UUID   `json:"realm_id" db:"realm_id"`
	Login    string `json:"login" db:"login"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

const (
	adminSchema = `
		CREATE TABLE admins (
			uuid BINARY(16) NOT NULL UNIQUE,
			realm_id BINARY(16) NOT NULL,
			login VARCHAR(256) NOT NULL,
			password VARCHAR(128) NOT NULL,
			name VARCHAR(256) NOT NULL,
			email VARCHAR(256) NOT NULL UNIQUE,
			
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (uuid),
			FOREIGN KEY (realm_id) REFERENCES realms (uuid)
		);
	`

	insertAdminQuery = `
		INSERT INTO admins 
			(uuid, realm_id, login, name, email, password) 
		VALUES 
			(:uuid, :realm_id, :login, :name, :email, :password);
	`
)

// Validate admin data
func (a *Admin) Validate() error {
	if !govalidator.IsEmail(a.Email) {
		return errors.New("invalid email")
	}
	if !govalidator.IsLowerCase(a.Login) {
		return errors.New("login should be lowercase")
	}

	if len(a.Password) < 8 {
		return errors.New("password is not long enough")
	}

	return nil
}

// HashPassword hash user's password
func (a *Admin) HashPassword() {
	password, _ := bcrypt.GenerateFromPassword([]byte(a.Password), 12)
	a.Password = string(password)
}

// Insert user in DB
func (a *Admin) Insert() error {
	a.UUID = NewUUID()

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(insertAdminQuery, a)
	if err != nil {
		return err
	}

	zap.S().Info("User ", a.Name, " just created. ("+a.Email+")")
	return nil
}
