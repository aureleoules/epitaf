package chronos

import "time"

// Calendar struct
type Calendar struct {
	ID      int `json:"Id"`
	DayList []struct {
		CourseList []struct {
			ID        int       `json:"Id"`
			Name      string    `json:"Name"`
			BeginDate time.Time `json:"BeginDate"`
			EndDate   time.Time `json:"EndDate"`
			Duration  int       `json:"Duration"`
			StaffList []struct {
				ID       int    `json:"Id"`
				Name     string `json:"Name"`
				Type     int    `json:"Type"`
				ParentID int    `json:"ParentId"`
			} `json:"StaffList"`
			RoomList []struct {
				Rooms    interface{} `json:"Rooms"`
				ID       int         `json:"Id"`
				Name     string      `json:"Name"`
				Type     int         `json:"Type"`
				ParentID int         `json:"ParentId"`
			} `json:"RoomList"`
			GroupList []struct {
				Groups   interface{} `json:"Groups"`
				ID       int         `json:"Id"`
				Name     string      `json:"Name"`
				Type     int         `json:"Type"`
				ParentID int         `json:"ParentId"`
			} `json:"GroupList"`
			Code interface{} `json:"Code"`
			Type interface{} `json:"Type"`
			URL  interface{} `json:"Url"`
			Info interface{} `json:"Info"`
		} `json:"CourseList"`
		DateTime time.Time `json:"DateTime"`
	} `json:"DayList"`
}
