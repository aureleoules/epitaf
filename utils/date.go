package utils

import "time"

// TruncateDate truncates date by removing hours, minutes and seconds
func TruncateDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
