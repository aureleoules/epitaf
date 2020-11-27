package models

const (
	realmSchema = `
		CREATE TABLE realms (

			id VARCHAR(16) NOT NULL UNIQUE,
			slug VARCHAR(256) NOT NULL UNIQUE,
			name VARCHAR(256) NOT NULL,
			website_url VARCHAR(512) NOT NULL DEFAULT "",

			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (short_id)
		);
	`
)

// Realm struct
// A realm contains informatons about a school
type Realm struct {
	base

	ID   string `json:"id" db:"id"`
	Slug string `json:"slug" db:"slug"`
	URL  string `json:"url" db:"url"`

	Name       string `json:"name" db:"name"`
	WebsiteURL string `json:"website_url" db:"website_url"`
}
