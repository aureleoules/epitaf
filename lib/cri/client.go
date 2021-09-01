package cri

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/resty.v1"
)

const (
	endpoint = "https://cri.epita.fr/api/v2"
)

// Client struct
type Client struct {
	httpClient *resty.Client
}

// SearchUser in CRI
func (c *Client) SearchUser(email string) (*ProfileResult, error) {
	zap.S().Info("Searching CRI user: ", email)

	resp, err := c.httpClient.R().
		Get(endpoint + "/users/search/?emails=" + email)

	if err != nil {
		return nil, err
	}
	var result ProfileSearchReq
	err = json.Unmarshal([]byte(resp.Body()), &result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	if strings.Contains(result.Detail, "Request was throttled") {
		time.Sleep(time.Second * 10)
		// Retry
		return c.SearchUser(email)
	}

	zap.S().Info("Fetched CRI user.")
	zap.S().Info(result)
	zap.S().Warn(string(resp.Body()))
	if len(result.Results) == 0 {
		return nil, errors.New("not found")
	}
	return &result.Results[0], err
}

// GetGroup from CRI
func (c *Client) GetGroup(groupSlug string) (*Group, error) {
	zap.S().Info("Fetching CRI group...")

	zap.S().Info(groupSlug)
	resp, err := c.httpClient.R().
		Get(endpoint + "/groups/" + groupSlug + "/")

	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	var result Group
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	zap.S().Info("Fetched CRI group.")
	return &result, err
}

// NewClient constructor
func NewClient(username string, password string, url *string) *Client {
	var e string
	if url == nil {
		e = endpoint
	} else {
		e = *url
	}
	c := Client{
		httpClient: resty.New().SetHostURL(e),
	}

	c.httpClient.SetBasicAuth(username, password)

	return &c
}
