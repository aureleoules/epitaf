package models

import (
	"database/sql"
	"time"
)

const (
	groupSchema = `
		CREATE TABLE groups (
			id VARCHAR(16) NOT NULL UNIQUE,
			realm_id VARCHAR(16) NOT NULL,
			usable BOOLEAN NOT NULL,
			slug VARCHAR(256) NOT NULL UNIQUE,
			name VARCHAR(256) NOT NULL,

			parent_id VARCHAR(16),
			active_at TIMESTAMP,
			archived BOOLEAN NOT NULL DEFAULT 0,
			archived_at TIMESTAMP,
			
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (short_id)
		);
	`
)

// Group struct
type Group struct {
	base

	RealmID string `json:"realm_id" db:"realm_id"`
	ID      string `json:"id" db:"id"`

	// Can add tasks to this group
	Usable sql.NullBool `json:"usable" db:"usable"`

	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`

	ParentID  string  `json:"parent_id" db:"parent_id"`
	Subgroups []Group `json:"subgroups" db:"-"`

	Archived   bool      `json:"archived" db:"archived"`
	ActiveAt   time.Time `json:"active_at" db:"active_at"`
	ArchivedAt time.Time `json:"archived_at" db:"archived_at"`
}
