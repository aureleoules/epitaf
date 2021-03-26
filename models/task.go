package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/utils"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"
)

const (
	insertTaskQuery = `
		INSERT INTO tasks 
			(short_id, promotion, visibility, members, class, region, semester, title, subject, content, due_date, created_by_login, updated_by_login) 
		VALUES 
			(:short_id, :promotion, :visibility, :members, :class, :region, :semester, :title, :subject, :content, :due_date, :created_by_login, :updated_by_login) 
		`
	getTaskQuery = `
		SELECT 
			tasks.short_id,
			tasks.promotion,
			tasks.visibility,
			tasks.members,
			tasks.class,
			tasks.region,
			tasks.semester,
			tasks.title,
			tasks.subject,
			tasks.content,
			tasks.due_date,
			ct.task_id IS NOT NULL as completed,
			ct.completed_at as completed_at,
			tasks.created_by_login,
			tasks.updated_by_login,
			users.name as created_by,
			updated_user.name as updated_by,
			tasks.created_at,
			tasks.updated_at
		FROM tasks
		LEFT JOIN users
		ON 
			users.login = tasks.created_by_login
		LEFT JOIN users as updated_user
		ON
			updated_user.login = tasks.updated_by_login
		LEFT JOIN completed_tasks as ct
		ON 
			ct.task_id = tasks.short_id AND ct.login = ?
		WHERE short_id = ? AND deleted=0;
	`

	getTasksRangeQuery = `
		SELECT
			tasks.short_id,
			tasks.promotion,
			tasks.visibility,
			tasks.members,
			tasks.class,
			tasks.title,
			tasks.subject,
			tasks.semester,
			tasks.content,
			tasks.region,
			tasks.due_date,
			ct.task_id IS NOT NULL as completed,
			ct.completed_at as completed_at,
			users.name as created_by,
			updated_user.name as updated_by,
			tasks.created_by_login,
			tasks.updated_by_login,
			tasks.created_at,
			tasks.updated_at
		FROM tasks
		LEFT JOIN users
		ON 
			users.login = tasks.created_by_login
		LEFT JOIN users as updated_user
		ON
			updated_user.login = tasks.updated_by_login
		LEFT JOIN completed_tasks as ct
		ON
			ct.task_id = tasks.short_id AND ct.login = ?
		WHERE 
			tasks.visibility=COALESCE(?, tasks.visibility)
			AND tasks.subject=COALESCE(?, tasks.subject)
			AND (
				COALESCE(?, ct.task_id IS NOT NULL) = (ct.task_id IS NOT NULL)
			)
			AND 
			(
				(
					tasks.visibility = 'self'
					AND tasks.created_by_login = ?
				) 
				OR (
					tasks.visibility = 'class'
					AND tasks.promotion = ?
					AND tasks.class = ?
					AND tasks.region = ?
					AND tasks.semester = ?
				) 
				OR (
					tasks.visibility = 'promotion'
					AND tasks.promotion = ? 
					AND tasks.semester = ? 
				)
				OR (
					tasks.visibility = 'students'
					AND (tasks.members LIKE ? OR tasks.created_by_login = ?)
				)
			)
			AND deleted=0
			AND due_date >= ? 
			AND due_date <= ?;
	`

	updateTaskQuery = `
		UPDATE tasks
		SET
			title=COALESCE(NULLIF(:title,''), tasks.title),
			subject=COALESCE(NULLIF(:subject,''), tasks.subject),
			visibility=COALESCE(NULLIF(:visibility,''), tasks.visibility),
			members=IF(visibility = 'students', COALESCE(NULLIF(:members,''), tasks.members), NULL),
			content=COALESCE(NULLIF(:content,''), tasks.content),
			due_date=COALESCE(NULLIF(:due_date,''), tasks.due_date),
			semester=:semester,
			region=:region,
			class=:class,
			promotion=:promotion,
			updated_by_login=:updated_by_login
		WHERE short_id=:short_id
	`

	deleteTaskQuery = `
		UPDATE tasks set deleted=1 WHERE short_id=?;
	`

	markTaskQuery = `
		INSERT INTO completed_tasks 
			(task_id, login)
		VALUES
			(?, ?);
	`

	unMarkTaskQuery = `
		DELETE FROM completed_tasks WHERE task_id=? AND login=?;
	`
)

// Members string
type Members []string

// Includes checks if s is included in slice
func (m Members) Includes(s string) bool {
	for _, a := range m {
		if a == s {
			return true
		}
	}
	return false
}

// MarshalJSON interface method
func (m Members) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(m))
}

// UnmarshalJSON interface method
func (m *Members) UnmarshalJSON(b []byte) error {
	var a []string
	err := json.Unmarshal(b, &a)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	*m = Members(a)
	return nil
}

// String method
func (m Members) String() string {
	l := ""
	for i, s := range m {
		l += string(s)
		if i != len(m)-1 {
			l += ","
		}
	}

	return l
}

// Value of members
func (m Members) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}

	return driver.Value(m.String()), nil
}

// Scan Members
func (m *Members) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	*m = Members(strings.Split(string(src.([]byte)), ","))
	return nil
}

// Visibility enum
type Visibility string

const (
	// SelfVisibility only the author of task can access it
	SelfVisibility Visibility = "self"
	// GroupVisibility users part of group or child groups can access it
	GroupVisibility Visibility = "group"
	// StudentsVisibility only selected students can access it
	StudentsVisibility Visibility = "students"
)

// Task truct
type Task struct {
	base

	ID uint `json:"id" db:"id"`

	// Meta
	RealmID UUID   `json:"realm_id" db:"realm_id"`
	ShortID string `json:"short_id" db:"short_id"`

	Visibility Visibility `json:"visibility" db:"visibility"`
	GroupID    UUID       `json:"group_id" db:"group_id"`
	Members    Members    `json:"members" db:"members"`

	// Body
	Title   string    `json:"title" db:"title"`
	Subject string    `json:"subject" db:"subject"`
	Content string    `json:"content" db:"content"`
	DueDate time.Time `json:"due_date" db:"due_date"`

	Completed   bool       `json:"completed" db:"completed"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`

	Deleted bool `json:"-" db:"deleted"`

	// Meta
	CreatedByLogin string `json:"created_by_login" db:"created_by_login"`
	CreatedBy      string `json:"created_by" db:"created_by"`
	UpdatedByLogin string `json:"updated_by_login" db:"updated_by_login"`
	UpdatedBy      string `json:"updated_by" db:"updated_by"`
}

// PrepareUpdate prepares new update data
func (t Task) PrepareUpdate(def Task, user User) Task {
	var update Task
	// Due date
	if t.DueDate.IsZero() {
		update.DueDate = def.DueDate
	} else {
		update.DueDate = t.DueDate
	}
	update.DueDate = update.DueDate.Local()

	// Content
	update.Content = setValueDefaultString(t.Content, def.Content)
	update.Title = setValueDefaultString(t.Title, def.Title)
	update.Subject = setValueDefaultString(t.Subject, def.Subject)

	return update
}

// Validate task data
func (t *Task) Validate() error {

	if t.Title == "" {
		return errors.New("title empty")
	}
	if t.Subject == "" {
		return errors.New("subject empty")
	}
	if t.Content == "" {
		return errors.New("content empty")
	}

	// Truncate minutes and seconds of due date
	t.DueDate = utils.TruncateDate(t.DueDate)

	if t.DueDate.Before(utils.TruncateDate(time.Now())) {
		return errors.New("invalid due date")
	}

	if len(t.Members) == 0 {
		t.Members = nil
	}

	return nil
}

// Insert task in DB
func (t *Task) Insert() error {
	t.ShortID = shortid.MustGenerate()

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(insertTaskQuery, t)
	return err
}

// Mark task as done by user
func (t *Task) Mark(login string) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.Exec(markTaskQuery, t.ShortID, login)
	return err
}

// Unmark task by user
func (t *Task) Unmark(login string) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.Exec(unMarkTaskQuery, t.ShortID, login)
	return err
}

// DeleteTask from db
func DeleteTask(id string) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.Exec(deleteTaskQuery, id)
	return err
}

// UpdateTask in DB
func UpdateTask(task Task) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)

	_, err = tx.NamedExec(updateTaskQuery, task)
	return err
}

// GetUserTask by shortID
func GetUserTask(id, login string) (*Task, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer checkErr(tx, err)

	var task Task
	err = tx.Get(&task, getTaskQuery, login, id)
	return &task, err
}
