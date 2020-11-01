package chronos

import (
	"encoding/json"

	"gopkg.in/resty.v1"
)

const (
	endpoint = "https://v2ssl.webservices.chronos.epita.net/api/v2"
)

// Client struct
type Client struct {
	httpClient *resty.Client
}

// GetGroupPlanning from Chronos
func (c *Client) GetGroupPlanning(groupSlug string) (*ChronosCalendar, error) {
	resp, err := c.httpClient.R().
		Get("/Planning/GetRangeWeekRecursive/" + groupSlug + "/1")

	if err != nil {
		return nil, err
	}

	var result []ChronosCalendar
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	var days ChronosCalendar
	for _, r := range result {
		days.DayList = append(days.DayList, r.DayList...)
	}

	return &days, err
}

// NewClient constructor
func NewClient(token string, url *string) *Client {
	var e string
	if url == nil {
		e = endpoint
	} else {
		e = *url
	}
	c := Client{
		httpClient: resty.New().SetHostURL(e),
	}

	c.httpClient.Header.Set("Auth-Token", token)

	return &c
}
