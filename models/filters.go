package models

import (
	"time"

	"github.com/aureleoules/epitaf/utils"
)

// Filters struct
type Filters struct {
	StartDate time.Time `json:"start_date" form:"start_date"`
	EndDate   time.Time `json:"end_date" form:"end_date"`

	Completed  *bool       `json:"completed" form:"completed"`
	Visibility *Visibility `json:"visibility" form:"visibility"`
	Subject    *string     `json:"subject" form:"subject"`
}

// Validate filters
func (f *Filters) Validate() error {
	// Default values
	if f.StartDate.IsZero() {
		f.StartDate = utils.TruncateDate(time.Now())
	}
	if f.EndDate.IsZero() {
		f.EndDate = utils.TruncateDate(time.Now().Add(time.Hour * 24 * 365))
	}

	return nil
}
