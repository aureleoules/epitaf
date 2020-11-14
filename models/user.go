package models

import (
	"errors"
	"os"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/lib/cri"
	"go.uber.org/zap"
)

const (
	userSchema = `
		CREATE TABLE users (
			login VARCHAR(256) NOT NULL UNIQUE,

			name VARCHAR(256) NOT NULL,
			email VARCHAR(256) NOT NULL UNIQUE,
			promotion VARCHAR(256) NOT NULL DEFAULT 0,
			class VARCHAR(256) NOT NULL DEFAULT "",
			region VARCHAR(256) NOT NULL DEFAULT "",
			semester VARCHAR(256) NOT NULL DEFAULT "",
			teacher BOOLEAN DEFAULT 0,
			
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (login)
		);
	`

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
			promotion,
			class,
			region,
			semester,
			teacher,
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
			promotion,
			class,
			region,
			semester,
			teacher,
			created_at,
			updated_at
		FROM users
		WHERE login = ?;
	`

	searchUserQuery = `
		SELECT 
			login,
			name, 
			email, 
			promotion,
			class,
			region,
			semester,
			teacher,
			created_at,
			updated_at
		FROM users
		WHERE 
			name LIKE ?
			OR login LIKE ?;
	`
)

// User struct
type User struct {
	base

	Login     string `json:"login" db:"login"`
	Name      string `json:"name" db:"name"`
	Promotion int    `json:"promotion" db:"promotion"`
	Class     string `json:"class" db:"class"`
	Region    string `json:"region" db:"region"`
	Semester  string `json:"semester" db:"semester"`
	Email     string `json:"email" db:"email"`
	Teacher   bool   `json:"teacher" db:"teacher"`
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
func (c *User) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(insertUserQuery, c)
	if err != nil {
		return err
	}

	zap.S().Info("User ", c.Name, " just created. ("+c.Email+")")
	return nil
}

// PrepareUser data
func PrepareUser(email string) (User, error) {
	zap.S().Info("Preparing user data...")
	user := User{
		Email: email,
	}
	// CRI req
	client := cri.NewClient(os.Getenv("CRI_USERNAME"), os.Getenv("CRI_PASSWORD"), nil)
	r, err := client.SearchUser(email)
	if err != nil {
		return user, errors.New("not found")
	}

	user.Name = r.FirstName + " " + r.LastName
	user.Login = r.Login

	if r.PrimaryGroup.Slug == "teachers" {
		user.Teacher = true
	}

	var slug string
	for i := len(r.GroupsHistory) - 1; i >= 0; i-- {
		g := r.GroupsHistory[i]
		if g.IsCurrent {
			slug = g.Group.Slug
			user.Promotion = g.GraduationYear
			break
		}
	}

	if slug == "" {
		return user, nil
	}

	group, err := client.GetGroup(slug)
	if err != nil {
		return user, jwt.ErrFailedAuthentication
	}

	g := strings.Split(group.Name, " ")
	user.Semester = g[0]
	user.Region = g[1]
	if len(g) > 2 {
		user.Class = g[2]
	}

	zap.S().Info("Prepared user data.")
	return user, nil
}
