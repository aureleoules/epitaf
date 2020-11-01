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
			uuid BINARY(16) NOT NULL,
			
			name VARCHAR(256) NOT NULL,
			email VARCHAR(256) NOT NULL UNIQUE,
			promotion VARCHAR(256) NOT NULL DEFAULT 0,
			class VARCHAR(256) NOT NULL DEFAULT "",
			region VARCHAR(256) NOT NULL DEFAULT "",
			semester VARCHAR(256) NOT NULL DEFAULT "",
			teacher BOOLEAN DEFAULT 0,
			
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (uuid),
			UNIQUE INDEX (email)
		);
	`

	insertUserQuery = `
		INSERT INTO users 
			(uuid, name, email, promotion, class, region, semester, teacher) 
		VALUES 
			(:uuid, :name, :email, :promotion, :class, :region, :semester, :teacher);
	`

	getUserByEmailQuery = `
		SELECT 
			uuid, 
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
			uuid, 
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
		WHERE uuid = ?;
	`
)

// User struct
type User struct {
	base

	Name      string `json:"name" db:"name"`
	Promotion int    `json:"promotion" db:"promotion"`
	Class     string `json:"class" db:"class"`
	Region    string `json:"region" db:"region"`
	Semester  string `json:"semester" db:"semester"`
	Email     string `json:"email" db:"email"`
	Teacher   bool   `json:"teacher" db:"teacher"`
}

// MicrosoftProfile struct
type MicrosoftProfile struct {
	DisplayName string `json:"displayName"`
	Mail        string `json:"mail"`
}

// GetUserByEmail retrives user by email
func GetUserByEmail(email string) (*User, error) {
	zap.S().Info("Retrieving user by email...")
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

	zap.S().Info("Retrieved user by email.")
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

	if r.PrimaryGroup.Slug == "teachers" {
		user.Teacher = true
	}

	var slug string
	for _, g := range r.GroupsHistory {
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
