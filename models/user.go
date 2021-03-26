package models

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/aureleoules/epitaf/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	insertUserQuery = `
		INSERT INTO users 
			(login, name, email, promotion, class, region, semester, teacher) 
		VALUES 
			(:login, :name, :email, :promotion, :class, :region, :semester, :teacher);
	`

	getUserByEmailQuery = `
		SELECT 
			login,
			name, 
			email, 
			created_at,
			updated_at
		FROM users
		WHERE email = ?;
	`

	getUserQuery = `
		SELECT 
			login,
			name, 
			email, 
			created_at,
			updated_at
		FROM users
		WHERE login = ? AND realm_id=?;
	`

	searchUserQuery = `
		SELECT 
			login,
			name, 
			email, 

			created_at,
			updated_at
		FROM users
		WHERE 
			(name LIKE ?
			OR login LIKE ?) AND realm_id=?;
	`
)

// User struct
type User struct {
	base

	ID       UUID   `json:"id" db:"id"`
	RealmID  UUID   `json:"realm_id" db:"realm_id"`
	Login    string `json:"login" db:"login"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

// Validate user data
func (u *User) Validate() error {
	if !govalidator.IsEmail(u.Email) {
		return errors.New("invalid email")
	}
	if !govalidator.IsLowerCase(u.Login) {
		return errors.New("login should be lowercase")
	}

	if len(u.Password) < 8 {
		return errors.New("password is not long enough")
	}

	return nil
}

// HashPassword hash user's password
func (u *User) HashPassword() {
	password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	u.Password = string(password)
}

// GetUserByEmail retrives user by email
func GetUserByEmail(email string) (*User, error) {
	zap.S().Info("Retrieving user by email...")
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var user User
	err = tx.Get(&user, getUserByEmailQuery, email)
	if err != nil {
		return nil, errors.New("not found")
	}
	zap.S().Info("Retrieved user by email.")
	return &user, err
}

func GetGroupUsers(realmID, groupID UUID) ([]User, error) {
	q, args, err := psql.
		Select("u.*").
		From("users AS u").
		InnerJoin("group_users AS gu ON gu.user_id = u.id").
		Where("u.realm_id = ? AND gu.group_id = ?", realmID, groupID).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var users []User
	err = tx.Select(&users, q, args...)

	return users, err
}

// GetUser retrives user by login
func GetUser(login string) (*User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var user User
	err = tx.Get(&user, getUserQuery, login)
	return &user, err
}

// SearchUser returns slice of users
func SearchUser(query string) ([]User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var users []User

	i := "%" + query + "%"
	err = tx.Select(&users, searchUserQuery, i, i)
	if err != nil {
		zap.S().Error(err)
	}
	return users, err
}

// Insert user in DB
func (u *User) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(insertUserQuery, u)
	if err != nil {
		return err
	}

	zap.S().Info("User ", u.Name, " just created. ("+u.Email+")")
	return nil
}

// GetRealmUsers return all users of a realm
func GetRealmUsers(realmID UUID) ([]User, error) {
	q, args, err := psql.
		Select("u.*").
		From("users AS u").
		Where("u.realm_id = ?", realmID).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var users []User
	err = tx.Select(&users, q, args...)

	return users, err
}
