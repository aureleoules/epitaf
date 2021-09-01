package models

import (
	"errors"
	"os"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/lib/cri"
	"github.com/mattn/go-nulltype"
	"go.uber.org/zap"
)

const (
	userSchema = `
		CREATE TABLE users (
			login VARCHAR(256) NOT NULL UNIQUE,

			name VARCHAR(256) NOT NULL,
			email VARCHAR(256) NOT NULL UNIQUE,
			promotion VARCHAR(256),
			class VARCHAR(256),
			region VARCHAR(256),
			semester VARCHAR(256),
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

	updateUserQuery = `
		UPDATE users
			SET 
				promotion = :promotion,
				class = :class,
				region = :region,
				semester = :semester
		WHERE login = :login;
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

	getUsersQuery = `
		SELECT * FROM users;
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

	Login     string              `json:"login" db:"login"`
	Name      string              `json:"name" db:"name"`
	Promotion nulltype.NullInt64  `json:"promotion" db:"promotion"`
	Class     nulltype.NullString `json:"class" db:"class"`
	Region    nulltype.NullString `json:"region" db:"region"`
	Semester  nulltype.NullString `json:"semester" db:"semester"`
	Email     string              `json:"email" db:"email"`
	Teacher   bool                `json:"teacher" db:"teacher"`
}

type UpdateUserReq struct {
	Promotion int    `db:"promotion"`
	Class     string `db:"class"`
	Region    string `db:"region"`
	Semester  string `db:"semester"`
	Login     string `db:"login"`
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

// GetUser retrives user by login
func GetUsers() ([]*User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var users []*User
	err = tx.Select(&users, getUsersQuery)
	return users, err
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

// Update user in DB
func UpdateUser(update *UpdateUserReq) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(updateUserQuery, update)
	if err != nil {
		return err
	}

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

	var group *cri.Group
	for i := len(r.GroupsHistory) - 1; i >= 0; i-- {
		g := r.GroupsHistory[i]
		if g.IsCurrent {
			group, err = client.GetGroup(g.Group.Slug)
			if err != nil {
				return user, jwt.ErrFailedAuthentication
			}

			if group.Kind != "class" {
				continue
			}

			user.Promotion.Set(int64(g.GraduationYear))
			break
		}
	}

	if group == nil {
		return user, nil
	}

	g := strings.Split(group.Name, " ")
	user.Semester.Set(g[0])
	user.Region.Set(g[1])
	if len(g) > 2 {
		user.Class.Set(g[2])
	}

	zap.S().Info("Prepared user data.")
	return user, nil
}

// CanViewTask checks if an user can view a task
func (u User) CanViewTask(task Task) bool {
	// Author
	if task.CreatedByLogin == u.Login {
		return true
	}
	// Student is included in visibility
	if task.Visibility == StudentsVisibility && task.Members.Includes(u.Login) {
		return true
	}

	// Student is in promotion
	if task.Visibility == PromotionVisibility && u.Promotion == task.Promotion && u.Semester == task.Semester {
		return true
	}

	// Student is in class
	if task.Visibility == ClassVisibility && u.Promotion == task.Promotion && u.Semester == task.Semester && task.Class == u.Class && task.Region == u.Region {
		return true
	}

	// If teacher : can view class & promo
	if u.Teacher && (task.Visibility == ClassVisibility || task.Visibility == PromotionVisibility) {
		return true
	}

	return false
}

// CanEditTask checks if an user can edit a task
func (u User) CanEditTask(task Task) bool {
	// Permission is the same for now
	return u.CanViewTask(task)
}

// CanDeleteTask checks if an user can delete a task
func (u User) CanDeleteTask(task Task) bool {
	if task.CreatedByLogin == u.Login {
		return true
	}

	if u.Teacher && (task.Visibility == ClassVisibility || task.Visibility == PromotionVisibility) {
		return true
	}

	return false
}
