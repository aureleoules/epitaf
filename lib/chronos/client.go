package chronos

import (
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	endpoint = "https://v2ssl.webservices.chronos.epita.net/api/v2"
)

// Client struct
type Client struct {
	httpClient *resty.Client
}

// GetGroupPlanning from Chronos
func (c *Client) GetGroupPlanning(groupSlug string) (*Calendar, error) {
	zap.S().Info("Fetching Chronos data...")

	resp, err := c.httpClient.R().
		Get("/Planning/GetRangeWeekRecursive/" + groupSlug + "/1")

	if err != nil {
		return nil, err
	}

	var result []Calendar
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		zap.S().Error(string(resp.Body()))
		return nil, err
	}

	var days Calendar
	for _, r := range result {
		for _, d := range r.DayList {
			if d.DateTime.After(time.Now().Add(-time.Hour * 24)) {
				days.DayList = append(days.DayList, d)
			}
		}
	}

	zap.S().Info("Fetched Chronos data.")

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
