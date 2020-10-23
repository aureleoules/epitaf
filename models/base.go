package models

import (
	"time"
)

type base struct {
	UUID      UUID      `db:"uuid" json:"uuid"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
