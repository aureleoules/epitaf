package cri

import (
	"encoding/json"

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
func (c *Client) SearchUser(email string) (*CRIProfileSearchReq, error) {
	resp, err := c.httpClient.R().
		Get(endpoint + "/users/search/?emails=" + email)

	if err != nil {
		return nil, err
	}
	var result []CRIProfileSearchReq
	err = json.Unmarshal([]byte(resp.Body()), &result)
	if err != nil {
		return nil, err
	}
	return &result[0], err
}

// GetGroup from CRI
func (c *Client) GetGroup(groupSlug string) (*CRIGroup, error) {
	zap.S().Info(groupSlug)
	resp, err := c.httpClient.R().
		Get(endpoint + "/groups/" + groupSlug + "/")

	if err != nil {
		return nil, err
	}
	var result CRIGroup
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

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
