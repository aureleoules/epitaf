package cri

type ProfileSearchReq struct {
	Detail   string          `json:"detail"`
	Count    int             `json:"count"`
	Next     interface{}     `json:"next"`
	Previous interface{}     `json:"previous"`
	Results  []ProfileResult `json:"results"`
}

type ProfileResult struct {
	URL          string `json:"url"`
	Login        string `json:"login"`
	UID          int    `json:"uid"`
	PrimaryGroup struct {
		URL  string `json:"url"`
		Slug string `json:"slug"`
	} `json:"primary_group"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	GroupsHistory []struct {
		Group struct {
			URL  string `json:"url"`
			Slug string `json:"slug"`
		} `json:"group"`
		BeginAt        string `json:"begin_at"`
		EndAt          string `json:"end_at"`
		GraduationYear int    `json:"graduation_year"`
		IsCurrent      bool   `json:"is_current"`
	} `json:"groups_history"`
	CurrentGroups []struct {
		URL  string `json:"url"`
		Slug string `json:"slug"`
	} `json:"current_groups"`
	OldAccounts []interface{} `json:"old_accounts"`
	NewAccount  interface{}   `json:"new_account"`
}

// Group struct
type Group struct {
	URL        string        `json:"url"`
	Slug       string        `json:"slug"`
	Gid        int           `json:"gid"`
	Name       string        `json:"name"`
	Kind       string        `json:"kind"`
	MembersURL string        `json:"members_url"`
	HistoryURL string        `json:"history_url"`
	Managers   []interface{} `json:"managers"`
	Private    bool          `json:"private"`
}
