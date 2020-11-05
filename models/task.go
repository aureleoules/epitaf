package models

import (
	"errors"
	"time"

	"github.com/aureleoules/epitaf/db"
	"github.com/teris-io/shortid"
)

const (
	taskSchema = `
		CREATE TABLE tasks (
			uuid BINARY(16) NOT NULL,
			short_id VARCHAR(16) NOT NULL,
			
			promotion VARCHAR(256) NOT NULL,
			global BOOLEAN NOT NULL DEFAULT 0,
			class VARCHAR(256) NOT NULL,
			region VARCHAR(256) NOT NULL,
			semester VARCHAR(256) NOT NULL,
			title VARCHAR(256) NOT NULL,
			subject VARCHAR(256) NOT NULL,
			content TEXT NOT NULL DEFAULT "",
			due_date TIMESTAMP NOT NULL,
			created_by_id BINARY(16) NOT NULL,
			updated_by_id BINARY(16) NOT NULL,

			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),

			PRIMARY KEY (uuid),
			FOREIGN KEY (created_by_id) REFERENCES users (uuid),
			FOREIGN KEY (updated_by_id) REFERENCES users (uuid)
		);
	`

	insertTaskQuery = `
		INSERT INTO tasks 
			(uuid, short_id, promotion, global, class, region, semester, title, subject, content, due_date, created_by_id, updated_by_id) 
		VALUES 
			(:uuid, :short_id, :promotion, :global, :class, :region, :semester, :title, :subject, :content, :due_date, :created_by_id, :updated_by_id) 
		`
	getTaskQuery = `
		SELECT 
			tasks.uuid,
			tasks.short_id,
			tasks.promotion,
			tasks.global,
			tasks.class,
			tasks.region,
			tasks.semester,
			tasks.title,
			tasks.subject,
			tasks.content,
			tasks.due_date,
			users.name as created_by,
			updated_user.name as updated_by,
			tasks.created_at,
			tasks.updated_at
		FROM tasks
		LEFT JOIN users
		ON 
			users.uuid = tasks.created_by_id
		LEFT JOIN users as updated_user
		ON
			updated_user.uuid = tasks.updated_by_id
		WHERE short_id = ?;
	`

	getTasksRangeQuery = `
		SELECT
			tasks.uuid,
			tasks.short_id,
			tasks.promotion,
			tasks.global,
			tasks.class,
			tasks.title,
			tasks.subject,
			tasks.semester,
			tasks.content,
			tasks.region,
			tasks.due_date,
			users.name as created_by,
			updated_user.name as updated_by,
			tasks.created_by_id,
			tasks.updated_by_id,
			tasks.created_at,
			tasks.updated_at
		FROM tasks
		LEFT JOIN users
		ON 
			users.uuid = tasks.created_by_id
		LEFT JOIN users as updated_user
		ON
			updated_user.uuid = tasks.updated_by_id
		WHERE 
			(
				(
					tasks.promotion = ?
					AND tasks.class = ?
					AND tasks.region = ?
					AND tasks.semester = ?
				) OR (tasks.promotion = ? AND tasks.global = 1)
			)
			AND due_date > ? 
			AND due_date < ?;
	`

	getAllTasksRangeQuery = `
		SELECT
			tasks.uuid,
			tasks.short_id,
			tasks.promotion,
			tasks.global,
			tasks.class,
			tasks.title,
			tasks.subject,
			tasks.semester,
			tasks.content,
			tasks.region,
			tasks.due_date,
			users.name as created_by,
			updated_user.name as updated_by,
			tasks.created_by_id,
			tasks.updated_by_id,
			tasks.created_at,
			tasks.updated_at
		FROM tasks
		LEFT JOIN users
		ON 
			users.uuid = tasks.created_by_id
		LEFT JOIN users as updated_user
		ON
			updated_user.uuid = tasks.updated_by_id
		WHERE 
			due_date > ? 
			AND due_date < ?;
	`

	updateTaskQuery = `
		UPDATE tasks
		SET
			title=:title,
			subject=:subject,
			content=:content,
			updated_by_id=:updated_by_id,
			due_date=:due_date,
			semester=:semester,
			region=:region,
			class=:class,
			promotion=:promotion
		WHERE short_id=:short_id
	`

	deleteTaskQuery = `
		DELETE FROM tasks WHERE short_id=?
	`
)

// Task truct
type Task struct {
	base

	ShortID string `json:"short_id" db:"short_id"`

	Promotion int       `json:"promotion" db:"promotion"`
	Global    bool      `json:"global" db:"global"`
	Class     string    `json:"class" db:"class"`
	Region    string    `json:"region" db:"region"`
	Semester  string    `json:"semester" db:"semester"`
	Title     string    `json:"title" db:"title"`
	Subject   string    `json:"subject" db:"subject"`
	Content   string    `json:"content" db:"content"`
	DueDate   time.Time `json:"due_date" db:"due_date"`

	CreatedByID UUID   `json:"-" db:"created_by_id"`
	CreatedBy   string `json:"created_by" db:"created_by"`

	UpdatedByID UUID   `json:"-" db:"updated_by_id"`
	UpdatedBy   string `json:"updated_by" db:"updated_by"`
}

// Validate task data
func (t Task) Validate() error {
	if t.Title == "" {
		return errors.New("title empty")
	}
	if t.Subject == "" {
		return errors.New("subject empty")
	}
	if t.Content == "" {
		return errors.New("content empty")
	}
	if t.Promotion == 0 {
		return errors.New("no promotion")
	}

	if !t.Global {
		if t.Semester == "" {
			return errors.New("no semester")
		}
		if t.Region == "" {
			return errors.New("no region")
		}
		if t.Class == "" {
			return errors.New("no class")
		}
	}

	return nil
}

// Insert task in DB
func (t *Task) Insert() error {
	t.UUID = NewUUID()
	t.ShortID = shortid.MustGenerate()

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.NamedExec(insertTaskQuery, t)
	return err
}

// DeleteTask from db
func DeleteTask(id string) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec(deleteTaskQuery, id)
	return err
}

// UpdateTask in DB
func UpdateTask(task Task) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.NamedExec(updateTaskQuery, task)
	return err
}

// GetTask by shortID
func GetTask(id string) (*Task, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var task Task
	err = tx.Get(&task, getTaskQuery, id)
	return &task, err
}

// GetTasksRange returns list of tasks in a time for a specific class promotion
func GetTasksRange(promotion int, semester string, class string, region string, start, end time.Time) ([]Task, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var tasks []Task
	err = tx.Select(&tasks, getTasksRangeQuery, promotion, class, region, semester, promotion, start, end)
	return tasks, err
}

// GetAllTasksRange returns list of tasks in a time range (for teachers)
func GetAllTasksRange(start, end time.Time) ([]Task, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var tasks []Task
	err = tx.Select(&tasks, getAllTasksRangeQuery, start, end)
	return tasks, err
}
