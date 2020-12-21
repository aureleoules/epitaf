package models

import (
	"github.com/aureleoules/epitaf/db"
)

// Realm struct
// A realm contains informatons about a school
type Realm struct {
	base

	UUID UUID   `json:"uuid" db:"uuid"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	URL  string `json:"url" db:"-"`

	WebsiteURL string `json:"website_url" db:"website_url"`
}

const (
	realmSchema = `
		CREATE TABLE realms (
			uuid BINARY(16) NOT NULL UNIQUE,
			name VARCHAR(256) NOT NULL,
			slug VARCHAR(256) NOT NULL UNIQUE,

			website_url VARCHAR(1024),
			
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (uuid)
		);
	`
	insertRealmQuery = `
		INSERT INTO realms 
			(uuid, name, slug, website_url) 
		VALUES 
			(:uuid, :name, :slug, :website_url);
	`

	getRealmBySlugQuery = `
		SELECT
			uuid,
			name,
			slug,
			website_url,
			created_at,
			updated_at
		FROM realms
		WHERE slug=?;
	`
)

// Insert realm into db
func (r *Realm) Insert() error {
	r.UUID = NewUUID()

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(insertRealmQuery, r)
	if err != nil {
		return err
	}

	return nil
}

// GetRealmBySlug retrieve realm by slug
func GetRealmBySlug(slug string) (*Realm, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var realm Realm
	err = tx.Get(&realm, getRealmBySlugQuery, slug)
	return &realm, err
}
