package microsoft

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"go.uber.org/zap"
	"gopkg.in/resty.v1"
)

const (
	endpoint     = "https://login.microsoftonline.com/3534b3d7-316c-4bc9-9ede-605c860f49d2/oauth2/v2.0"
	userEndpoint = "https://graph.microsoft.com/v1.0"
)

// Client struct
type Client struct {
	httpClient *resty.Client
}

// NewClient constructor
func NewClient(token string, url *string) *Client {
	var e string
	if url == nil {
		e = userEndpoint
	} else {
		e = *url
	}
	c := Client{
		httpClient: resty.New().SetHostURL(e),
	}

	c.httpClient.Header.Set("Authorization", "Bearer "+token)
	return &c
}

// SignInURL returns url which the user must go to
func SignInURL(redirectURI string) string {
	// Prepare microsoft query
	req, _ := http.NewRequest("GET", endpoint+"/authorize", nil)
	q := req.URL.Query()

	q.Add("client_id", os.Getenv("CLIENT_ID"))
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("state", "0000")
	q.Add("scope", "https://graph.microsoft.com/User.Read")
	q.Add("redirect_uri", redirectURI)

	req.URL.RawQuery = q.Encode()
	// Return URL that the user must go to
	return req.URL.String()
}

// GetProfile retrieves MicrosoftProfile of user
func (c *Client) GetProfile() (Profile, error) {
	var result Profile

	resp, err := c.httpClient.R().
		Get(userEndpoint + "/me")

	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(resp.Body()), &result)
	if err != nil {
		return result, err
	}

	if result.Mail == "" {
		return result, errors.New("invalid token")
	}

	zap.S().Info("Fetched Microsoft profile...")
	return result, nil
}

// GetAccessToken retrieves access_token from authentication code
func GetAccessToken(code string, uri string) (string, error) {
	form := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"client_id":    {os.Getenv("CLIENT_ID")},
		"redirect_uri": {uri},
	}

	// Do not include secret for mobile authentication
	if !strings.HasPrefix(uri, "epitaf://") {
		form.Add("client_secret", os.Getenv("CLIENT_SECRET"))
	}

	resp, err := http.PostForm(endpoint+"/token", form)
	if err != nil {
		zap.S().Error(err)
		return "", jwt.ErrFailedAuthentication
	}

	defer func() {
		if resp.Body.Close() != nil {
			zap.S().Warn("could not close body")
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.S().Error(err)
		return "", jwt.ErrFailedAuthentication
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		zap.S().Error(err)
		return "", jwt.ErrFailedAuthentication
	}

	// Return access token
	token := result["access_token"].(string)
	if token == "" {
		zap.S().Error("no access token")
		return "", jwt.ErrFailedAuthentication
	}
	return token, nil
}
