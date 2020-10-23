package models

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

func GetUserByEmail(email string) (User, error) {
	return User{}, nil
}

// Insert User in DB
func (u *User) Insert() error {
	return nil
}
