package models

import (
	"time"

	"github.com/aureleoules/epitaf/db"
	"github.com/mattn/go-nulltype"
)

// Group struct
type Group struct {
	base

	ID      UUID `json:"id" db:"id"`
	RealmID UUID `json:"realm_id" db:"realm_id"`

	// Can add tasks to this group
	Usable nulltype.NullBool `json:"usable" db:"usable"`

	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`

	// Dynamic fields
	Users     []User    `json:"users,omitempty" db:"-"`
	Subjects  []Subject `json:"subjects,omitempty" db:"-"`
	Subgroups []Group   `json:"subgroups,omitempty" db:"-"`

	ParentID *UUID `json:"parent_id" db:"parent_id"`

	Archived   bool       `json:"archived" db:"archived"`
	ArchivedAt *time.Time `json:"archived_at" db:"archived_at"`

	Active   bool      `json:"active" db:"active"`
	ActiveAt time.Time `json:"active_at" db:"active_at"`
}

// Insert group in DB
func (g *Group) Insert() error {
	g.ID = NewUUID()

	q, args, err := psql.Insert("groups").
		Columns("id", "realm_id", "name", "slug", "parent_id", "active", "active_at", "archived", "archived_at").
		Values(g.ID, g.RealmID, g.Name, g.Slug, g.ParentID, g.Active, g.ActiveAt, g.Archived, g.ArchivedAt).
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

func getGroup(realmID UUID, id UUID) (*Group, error) {
	q, args, err := psql.Select("g.*").
		From("groups AS g").
		Where("realm_id = ? AND id = ?", realmID, id).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var group Group
	err = tx.Get(&group, q, args...)

	return &group, err
}

// GetGroup returns group by uuid
func GetGroup(realmID, id UUID) (*Group, error) {
	group, err := getGroup(realmID, id)
	if err != nil {
		return nil, err
	}
	group.Users, err = GetGroupUsers(realmID, id)
	if err != nil {
		return nil, err
	}
	group.Subjects, err = GetGroupSubjects(id)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// GetGroupTree returns group tree of a realm
func GetGroupTree(realmID UUID) (*Group, error) {
	group, err := getRootGroup(realmID)
	if err != nil {
		return nil, err
	}

	group.Subgroups, err = getSubGroupsRec(realmID, group.ID)
	if err != nil {
		return nil, err
	}

	return group, err
}

// GetSubGroupsRec finds groups of subgroups in a realm
func getSubGroupsRec(realmID UUID, parentID UUID) ([]Group, error) {
	groups, err := GetSubGroups(realmID, parentID)
	if err != nil {
		return nil, err
	}

	for i, g := range groups {
		groups[i].Subgroups, err = getSubGroupsRec(realmID, g.ID)
		if err != nil {
			return nil, err
		}
	}

	return groups, nil
}

func getRootGroup(realmID UUID) (*Group, error) {
	q, args, err := psql.Select("g.*").
		From("groups AS g").
		Where("realm_id = ? AND parent_id IS NULL", realmID).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var group Group
	err = tx.Get(&group, q, args...)

	return &group, err
}

// GetSubGroups retrieves subgroups of groups
func GetSubGroups(realmID UUID, parentID UUID) ([]Group, error) {
	q, args, err := psql.Select("g.*").
		From("groups AS g").
		Where("realm_id = ? AND parent_id = ?", realmID, parentID).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var groups []Group
	err = tx.Select(&groups, q, args...)

	return groups, err
}

// DeleteGroup removes a group from db
func DeleteGroup(realmID UUID, id UUID) error {
	return deleteSubGroupsRec(realmID, id)
}

// GetSubGroupsRec finds groups of subgroups in a realm
func deleteSubGroupsRec(realmID UUID, parentID UUID) error {
	groups, err := GetSubGroups(realmID, parentID)
	if err != nil {
		return err
	}

	for _, g := range groups {
		err = deleteSubGroupsRec(realmID, g.ID)
		if err != nil {
			return err
		}
	}

	return deleteGroup(realmID, parentID)
}

func deleteGroup(realmID, id UUID) error {
	q, args, err := psql.Delete("groups").
		Where("realm_id = ? AND id = ?", realmID, id).
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
	return err
}

// AddGroupUsers add users to a specific group
func AddGroupUsers(id UUID, userIDs []UUID) error {
	q := psql.Insert("group_users").
		Columns("group_id", "user_id")

	for _, userID := range userIDs {
		q = q.Values(id, userID)
	}

	query, args, err := q.ToSql()

	if err != nil {
		return err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.Exec(query, args...)
	return err
}

// IsUserInGroup checks if a user is already in a group
func IsUserInGroup(groupID, userID UUID) (bool, error) {
	q, args, err := psql.Select("COUNT(*)").
		From("group_users AS gu").
		Where("gu.group_id = ? AND gu.user_id = ?", groupID, userID).
		Limit(1).
		ToSql()

	if err != nil {
		return false, err
	}

	tx, err := db.DB.Beginx()
	if err != nil {
		return false, err
	}

	defer checkErr(tx, err)

	var c int
	err = tx.Get(&c, q, args...)
	return c > 0, err
}
