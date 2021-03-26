package models

import (
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
)

// Filters struct stores informations about a search query
type Filters struct {
	Query string `query:"query"`

	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`

	Limit  int `query:"limit"`
	Offset int `query:"offset"`

	Order string `query:"sort"`
}

// ApplyBase applies basic filters such as `created_at`, `limit` and `offset` to the SQL query
func (f *Filters) ApplyBase(q *squirrel.SelectBuilder, table string) {
	*q = q.
		OrderByClause(table + ".created_at " + f.OrderBy()).
		Limit(uint64(f.Limit)).
		Offset(uint64(f.Offset))

	if f.StartDate != "" {
		*q = q.Where(table+".created_at >= ?", f.StartDate)
	}
	if f.EndDate != "" {
		*q = q.Where(table+".created_at <= ?", f.EndDate)
	}
}

// Defaults sets default values of query
func (f *Filters) Defaults() {
	if f.Limit == 0 {
		f.Limit = 50
	}
	if f.StartDate == "" {
		f.StartDate = time.Now().Add(-time.Hour * 24 * 30).Format(time.RFC3339)
	}

	if f.EndDate == "" {
		f.EndDate = time.Now().Add(time.Hour * 24).Format(time.RFC3339)
	}
}

// OrderBy returns order of query
func (f Filters) OrderBy() string {
	if strings.ToLower(f.Order) == "asc" {
		return "ASC"
	}
	return "DESC"
}
