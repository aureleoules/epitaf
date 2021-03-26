package models

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/aureleoules/epitaf/db"
	"github.com/mattn/go-nulltype"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
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

	ID       UUID                `json:"id" db:"id"`
	RealmID  UUID                `json:"realm_id" db:"realm_id"`
	Login    string              `json:"login" db:"login"`
	Name     string              `json:"name" db:"name"`
	Email    string              `json:"email" db:"email"`
	Password nulltype.NullString `json:"password" db:"password"`
}

type UserFilters struct {
	Filters

	ExcludeGroup string `query:"exclude_group"`
}

// Validate user data
func (u *User) Validate() error {
	if !govalidator.IsEmail(u.Email) {
		return errors.New("invalid email")
	}
	if !govalidator.IsLowerCase(u.Login) {
		return errors.New("login should be lowercase")
	}

	if u.Password.Valid() && len(u.Password.String()) < 8 {
		return errors.New("password is not long enough")
	}

	return nil
}

// HashPassword hash user's password
func (u *User) HashPassword() {
	password, _ := bcrypt.GenerateFromPassword([]byte(u.Password.String()), 12)
	u.Password = nulltype.NullStringOf(string(password))
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
	u.ID = NewUUID()

	q, args, err := psql.
		Insert("users").
		Columns("id", "realm_id", "login", "name", "email", "password").
		Values(u.ID, u.RealmID, u.Login, u.Name, u.Email, u.Password).
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
	return err
}

// GetRealmUsers return all users of a realm
func GetRealmUsers(realmID UUID, filters UserFilters) ([]User, error) {
	filters.Defaults()

	// select u.name from users as u where not exists (select 1 from group_users as gu where gu.group_id = '\xab02c45e2018475770cd29852571de90'::bytea and gu.user_id = u.id);
	query := psql.
		Select("u.*").
		From("users AS u").
		Where("u.realm_id = ?", realmID)

	filters.ApplyBase(&query, "u")

	if filters.Query != "" {
		query = query.Where("LOWER(u.name) LIKE '%' || ? || '%'", filters.Query)
	}
	if filters.ExcludeGroup != "" {
		groupID, err := FromUUID(filters.ExcludeGroup)
		if err != nil {
			return nil, err
		}
		query = query.Where("NOT EXISTS (SELECT 1 FROM group_users AS gu WHERE gu.group_id = ? AND gu.user_id = u.id)", groupID)
	}

	q, args, err := query.ToSql()
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
