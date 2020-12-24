package models

import (
	"time"

	"github.com/aureleoules/epitaf/db"
	"github.com/mattn/go-nulltype"
	"go.uber.org/zap"
)

const (
	groupSchema = `
		CREATE TABLE groups (
			uuid BINARY(16) NOT NULL UNIQUE,
			realm_id BINARY(16) NOT NULL,
			usable BOOLEAN NOT NULL DEFAULT 1,
			slug VARCHAR(256) NOT NULL,
			name VARCHAR(256) NOT NULL,

			parent_id BINARY(16),
			active_at TIMESTAMP,
			archived BOOLEAN NOT NULL DEFAULT 0,
			archived_at TIMESTAMP,
			
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
			PRIMARY KEY (uuid),
			FOREIGN KEY (realm_id) REFERENCES realms (uuid),
			FOREIGN KEY (parent_id) REFERENCES groups (uuid)
		);
	`

	getRealmGroupsQuery = `
		SELECT 
			uuid,
			realm_id,
			usable,
			slug, 
			name, 
			parent_id,
			active_at,
			archived,
			archived_at
			created_at,
			updated_at
		FROM groups
		WHERE 
			realm_id=? AND parent_id=?;
	`

	getRootGroupQuery = `
		SELECT 
			uuid,
			realm_id,
			usable,
			slug, 
			name, 
			parent_id,
			active_at,
			archived,
			archived_at
			created_at,
			updated_at
		FROM groups
		WHERE 
			realm_id=? AND parent_id IS NULL;
	`

	getGroupQuery = `
		SELECT 
			uuid,
			realm_id,
			usable,
			slug, 
			name, 
			parent_id,
			active_at,
			archived,
			archived_at
			created_at,
			updated_at
		FROM groups
		WHERE 
			realm_id=? AND uuid=?;
	`

	insertGroupQuery = `
		INSERT INTO groups 
			(uuid, realm_id, name, slug, parent_id, active_at, archived, archived_at) 
		VALUES 
			(:uuid, :realm_id, :name, :slug, :parent_id, :active_at, archived, archived_at);
	`
)

// Group struct
type Group struct {
	base

	RealmID UUID `json:"realm_id" db:"realm_id"`
	UUID    UUID `json:"uuid" db:"uuid"`

	// Can add tasks to this group
	Usable nulltype.NullBool `json:"usable" db:"usable"`

	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`

	Users []*User `json:"users,omitempty" db:"-"`

	ParentID  *UUID    `json:"parent_id" db:"parent_id"`
	Subgroups []*Group `json:"subgroups,omitempty" db:"-"`

	Archived   bool       `json:"archived" db:"archived"`
	ActiveAt   time.Time  `json:"active_at" db:"active_at"`
	ArchivedAt *time.Time `json:"archived_at" db:"archived_at"`
}

// Insert group in DB
func (g *Group) Insert() error {
	g.UUID = NewUUID()

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(insertGroupQuery, g)
	if err != nil {
		return err
	}

	return nil
}

// GetGroup returns group by uuid
func GetGroup(realmID UUID, groupID UUID) (*Group, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var group Group

	err = tx.Get(&group, getGroupQuery, realmID, groupID)
	if err != nil {
		zap.S().Error(err)
	}

	return &group, err
}

// GetGroupTree returns group tree of a realm
func GetGroupTree(realmID UUID) (*Group, error) {
	group, err := getRootGroup(realmID)
	if err != nil {
		return nil, nil
	}

	group.Subgroups, err = getSubGroupsRec(realmID, group.UUID)
	if err != nil {
		return nil, nil
	}

	return group, err
}

// GetSubGroupsRec finds groups of subgroups in a realm
func getSubGroupsRec(realmID UUID, parentID UUID) ([]*Group, error) {
	groups, err := GetSubGroups(realmID, parentID)
	if err != nil {
		return nil, nil
	}

	for _, g := range groups {
		g.Subgroups, err = getSubGroupsRec(realmID, g.UUID)
		if err != nil {
			return nil, nil
		}
	}

	return groups, nil
}

func getRootGroup(realmID UUID) (*Group, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var group Group

	err = tx.Get(&group, getRootGroupQuery, realmID)
	if err != nil {
		zap.S().Error(err)
	}

	return &group, err
}

// GetSubGroups retrieves subgroups of groups
func GetSubGroups(realmID UUID, parentID UUID) ([]*Group, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var groups []*Group

	err = tx.Select(&groups, getRealmGroupsQuery, realmID, parentID)
	if err != nil {
		zap.S().Error(err)
	}

	return groups, err
}
