package zeus

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"gopkg.in/resty.v1"
)

const (
	endpoint = "https://zeus.ionis-it.com/api"
)

var (
	token string
)

func SetToken(t string) {
	token = t
}

// Client struct
type Client struct {
	httpClient *resty.Client
}

// GetGroupPlanning from Chronos
func (c *Client) GetGroupPlanning(groupId int) ([]Calendar, error) {
	zap.S().Info("Fetching Chronos data...")
	c.httpClient.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.R().
		SetBody(map[string]interface{}{
			"groups":    []int{groupId},
			"startDate": time.Now(),
			"endDate":   time.Now().Add(14 * time.Hour * 24),
		}).
		Post("/reservation/filter/displayable")

	if err != nil {
		return nil, err
	}

	var result []Calendar
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		zap.S().Error(string(resp.Body()))
		return nil, err
	}

	zap.S().Info("Fetched Chronos data.")

	return result, err
}

// GetGroupPlanning from Chronos
func (c *Client) GetICS(groupId string) (string, error) {
	c.httpClient.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.R().
		Get("/group/" + groupId + "/ics")

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

// NewClient constructor
func NewClient(url *string) *Client {
	var e string
	if url == nil {
		e = endpoint
	} else {
		e = *url
	}
	c := Client{
		httpClient: resty.New().SetHostURL(e),
	}

	c.httpClient.Header.Set("Authorization", "Bearer "+token)

	return &c
}
