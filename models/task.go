package models

import (
	"time"

	"github.com/aureleoules/epitaf/db"
	"github.com/mattn/go-nulltype"
)

// Task truct
type Task struct {
	base

	ID      UUID   `json:"id" db:"id"`
	ShortID string `json:"short_id" db:"short_id"`
	GroupID UUID   `json:"group_id" db:"group_id"`

	SubjectID *UUID    `json:"subject_id,omitempty" db:"subject_id"`
	Subject   *Subject `json:"subject,omitempty" db:"-"`

	Title        string              `json:"title" db:"title"`
	Content      nulltype.NullString `json:"content" db:"content"`
	DueDateStart time.Time           `json:"due_date_start" db:"due_date_start"`
	DueDateEnd   time.Time           `json:"due_date_end" db:"due_date_end"`

	Completed   bool       `json:"completed" db:"completed"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`

	Archived bool `json:"-" db:"archived"`
}

// Validate task data
func (t *Task) Validate() error {
	return nil
}

// Insert task in DB
func (t *Task) Insert() error {
	t.ID = NewUUID()

	q, args, err := psql.Insert("tasks").
		Columns("id", "group_id", "short_id", "title", "content", "due_date_start", "due_date_end").
		Values(t.ID, t.GroupID, t.ShortID, t.Title, t.Content, t.DueDateStart, t.DueDateEnd).
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

// Mark task as done by user
func (t *Task) Mark(login string) error {
	return nil
}

// Unmark task by user
func (t *Task) Unmark(login string) error {
	return nil
}

func ArchiveTask(id string) error {
	return nil
}

// UpdateTask in DB
func UpdateTask(task Task) error {
	return nil
}

func GetUserTask(id, login string) (*Task, error) {
	return nil, nil
}
