package models

import (
	"github.com/aureleoules/epitaf/db"
	"github.com/mattn/go-nulltype"
)

type Subject struct {
	base

	ID       UUID                `json:"id" db:"id"`
	GroupID  UUID                `json:"group_id" db:"group_id"`
	Name     string              `json:"name" db:"name"`
	Color    nulltype.NullString `json:"color" db:"color"`
	IconURL  nulltype.NullString `json:"icon_url" db:"icon_url"`
	Archived bool                `json:"archived" db:"archived"`

	Inherited bool `json:"inherited,omitempty" db:"-"`
}

// Insert subject into db
func (s *Subject) Insert() error {
	s.ID = NewUUID()

	q, args, err := psql.Insert("group_subjects").
		Columns("id", "group_id", "name", "color", "icon_url").
		Values(s.ID, s.GroupID, s.Name, s.Color, s.IconURL).
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

// ArchiveSubject archives a group subject
func ArchiveSubject(id UUID) error {
	q, args, err := psql.
		Update("group_subjects").
		Set("archived", true).
		Where("id = ?", id).
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

func getGroupSubjects(groupID UUID) ([]Subject, error) {
	q, args, err := psql.
		Select("gs.*").
		From("group_subjects AS gs").
		Where("gs.group_id = ? AND gs.archived = false", groupID).
		ToSql()

	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var subjects []Subject
	err = tx.Select(&subjects, q, args...)

	return subjects, err
}

// GetGroupSubjects returns a group's subjects
func GetGroupSubjects(groupID UUID, inheritance bool) ([]Subject, error) {
	if !inheritance {
		return getGroupSubjects(groupID)
	}

	var subjects []Subject

	var id *UUID = &groupID

	for id != nil {
		s, err := getGroupSubjects(*id)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, s...)

		g, err := getGroup(*id)
		if err != nil {
			return nil, err
		}
		if g.ParentID == nil {
			break
		}

		id = g.ParentID
	}

	return subjects, nil
}
