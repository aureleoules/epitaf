package models

// UserGroup group of user
type UserGroup struct {
	base

	UserID  UUID `json:"user_id" db:"user_id"`
	GroupID UUID `json:"group_id" db:"group_id"`
	Active  bool `json:"active" db:"active"`
}

const (
	userGroupSchema = `
		CREATE TABLE user_groups (
			user_id BINARY(16) NOT NULL,
			group_id BINARY(16) NOT NULL,
			active BOOLEAN NOT NULL DEFAULT 0,

			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (uuid),
			FOREIGN KEY (user_id) REFERENCES users (uuid),
			FOREIGN KEY (group_id) REFERENCES groups (uuid)
		);
	`
)
