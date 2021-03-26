package models

import (
	"github.com/aureleoules/epitaf/db"
)

// Realm struct
// A realm contains informatons about a school
type Realm struct {
	base

	ID   UUID   `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	URL  string `json:"url" db:"-"`

	WebsiteURL string `json:"website_url" db:"website_url"`
}

// Insert realm into db
func (r *Realm) Insert() error {
	r.ID = NewUUID()

	q, args, err := psql.Insert("realms").
		Columns("id", "name", "slug", "website_url").
		Values(r.ID, r.Name, r.Slug, r.WebsiteURL).
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
	if err != nil {
		return err
	}

	return nil
}

// GetRealmBySlug retrieve realm by slug
func GetRealmBySlug(slug string) (*Realm, error) {
	q, args, err := psql.Select("r.*").
		From("realms AS r").
		Where("slug = ?", slug).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var realm Realm
	err = tx.Get(&realm, q, args...)

	return &realm, err
}

// GetRealmOfAdmin retrieves realm informations of admin
func GetRealmOfAdmin(id UUID) (*Realm, error) {
	q, args, err := psql.Select("r.*").
		From("admins AS u").
		InnerJoin("realms AS r ON r.id = u.realm_id").
		Where("u.id = ?", id).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var realm Realm
	err = tx.Get(&realm, q, args...)

	return &realm, err
}
