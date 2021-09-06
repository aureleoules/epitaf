package zeus

import "time"

type Calendar struct {
	IDReservation int         `json:"idReservation"`
	IDCourse      interface{} `json:"idCourse"`
	Name          string      `json:"name"`
	IDType        int         `json:"idType"`
	StartDate     time.Time   `json:"startDate"`
	EndDate       time.Time   `json:"endDate"`
	IsOnline      bool        `json:"isOnline"`
	Rooms         []struct {
		ID         int    `json:"id"`
		Capacity   int    `json:"capacity"`
		Name       string `json:"name"`
		IDRoomType int    `json:"idRoomType"`
		IDLocation int    `json:"idLocation"`
	} `json:"rooms"`
	Groups []struct {
		ID         int         `json:"id"`
		IDParent   interface{} `json:"idParent"`
		Name       string      `json:"name"`
		Path       interface{} `json:"path"`
		Count      interface{} `json:"count"`
		IsReadOnly interface{} `json:"isReadOnly"`
		IDSchool   int         `json:"idSchool"`
		Color      string      `json:"color"`
	} `json:"groups"`
	Teachers   interface{} `json:"teachers"`
	IDSchool   int         `json:"idSchool"`
	SchoolName string      `json:"schoolName"`
}
