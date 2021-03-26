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

	ID       UUID   `json:"id" db:"id"`
	RealmID  UUID   `json:"realm_id" db:"realm_id"`
	Login    string `json:"login" db:"login"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}

const (
	getAdminByEmailQuery = `
		SELECT 
			uuid,
			login,
			name, 
			email,
			password,
			created_at,
			updated_at
		FROM admins
		WHERE email = ?;
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
	a.ID = NewUUID()

	q, args, err := psql.Insert("admins").
		Columns("id", "realm_id", "login", "name", "email", "password").
		Values(a.ID, a.RealmID, a.Login, a.Name, a.Email, a.Password).
		ToSql()

	if err != nil {
		return err
	}

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.Exec(q, args...)
	if err != nil {
		return err
	}

	zap.S().Info("User ", a.Name, " just created. ("+a.Email+")")
	return nil
}

// GetAdminByEmail retrieves admin by email
func GetAdminByEmail(email string) (*Admin, error) {
	q, args, err := psql.Select("a.*").
		From("admins AS a").
		Where("email = ?", email).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var admin Admin
	err = tx.Get(&admin, q, args...)

	return &admin, err
}

// GetAdmin retrieves admin by uuid
func GetAdmin(id UUID) (*Admin, error) {
	q, args, err := psql.Select("a.*").
		From("admins AS a").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var admin Admin
	err = tx.Get(&admin, q, args...)

	return &admin, err
}
