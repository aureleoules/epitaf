package models

import (
	"time"

	"github.com/aureleoules/epitaf/lib/chronos"
)

// Calendar struct
type Calendar struct {
	Days []Day `json:"days"`
}

// Day struct
type Day struct {
	Date    time.Time `json:"date"`
	Courses []Course  `json:"courses"`
}

// Course struct
type Course struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Duration  int       `json:"duration"`
	Staff     []string  `json:"staff"`
	Rooms     []string  `json:"rooms"`
}

// FormatCalendar formats Chronos' data into a more readable structure
// and merges courses together when they are the same
func FormatCalendar(cal chronos.Calendar) Calendar {
	var formatted Calendar

	for _, d := range cal.DayList {
		day := Day{
			Date: d.DateTime,
		}

		var currentCourse string
		var currentCourseIndex int
		for i := 0; i < len(d.CourseList); i++ {
			c := d.CourseList[i]

			if currentCourse == c.Name {
				current := &day.Courses[currentCourseIndex]

				current.Duration += c.Duration
				current.EndDate = current.StartDate.Add(time.Duration(current.Duration) * time.Minute)

				continue
			}

			// Format course
			course := Course{
				Name:      c.Name,
				StartDate: c.BeginDate,
			}

			// Format staff
			for _, t := range c.StaffList {
				course.Staff = append(course.Staff, t.Name)
			}

			// Format rooms
			for _, r := range c.RoomList {
				course.Rooms = append(course.Rooms, r.Name)
			}

			course.Duration = c.Duration
			course.EndDate = course.StartDate.Add(time.Duration(course.Duration) * time.Minute)

			// Append course
			day.Courses = append(day.Courses, course)

			currentCourse = c.Name
			currentCourseIndex = len(day.Courses) - 1
		}

		formatted.Days = append(formatted.Days, day)
	}

	return formatted
}
