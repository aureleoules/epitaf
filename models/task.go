package models

import "github.com/teris-io/shortid"

const (
	taskSchema = `
		CREATE TABLE tasks (
			uuid BINARY(16) NOT NULL,
			short_id VARCHAR(16) NOT NULL,
			
			promotion VARCHAR(256) NOT NULL,
			class VARCHAR(256) NOT NULL,
			title VARCHAR(256) NOT NULL,
			subject VARCHAR(256) NOT NULL,
			content VARCHAR(65536) NOT NULL DEFAULT "",
			
			created_by BINARY(16) NOT NULL,
			updated_by BINARY(16) NOT NULL,

			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),

			PRIMARY KEY (uuid),
			SECONDARY KEY (short_id),
			FOREIGN KEY (created_by) REFERENCES users (uuid),
			FOREIGN KEY (updated_by) REFERENCES users (uuid)
		);
	`
)

// Task truct
type Task struct {
	base

	ShortID shortid.Shortid `json:"short_id" db:"short_id"`

	Promotion int    `json:"promotion" db:"promotion"`
	Class     string `json:"class" db:"class"`
	Title     string `json:"title" db:"title"`
	Subject   string `json:"subject" db:"subject"`
	Content   string `json:"content" db:"content"`

	CreatedByID UUID   `json:"-" db:"created_by"`
	CreatedBy   string `json:"created_by" db:"-"`

	UpdatedByID UUID   `json:"-" db:"updated_by"`
	UpdatedBy   string `json:"updated_by" db:"-"`
}
