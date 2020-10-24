package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
)

const (
	// Endpoint api url
	Endpoint = "https://login.microsoftonline.com/3534b3d7-316c-4bc9-9ede-605c860f49d2/oauth2/v2.0"
)

func handleUsers() {
	users := api.Group("/users")

	users.POST("/authenticate", authenticateHandler)
	users.POST("/callback", auth.LoginHandler)
}

func authenticateHandler(c *gin.Context) {
	req, _ := http.NewRequest("GET", Endpoint+"/authorize", nil)
	q := req.URL.Query()

	q.Add("client_id", os.Getenv("CLIENT_ID"))
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("state", "0000")
	q.Add("scope", "https://graph.microsoft.com/User.Read")

	if os.Getenv("DEV") == "true" {
		q.Add("redirect_uri", "http://localhost:3000/callback")
	} else {
		q.Add("redirect_uri", "https://epitaf.aureleoules.com/callback")
	}

	req.URL.RawQuery = q.Encode()
	c.JSON(http.StatusOK, req.URL.String())
}

func callbackHandler(c *gin.Context) (interface{}, error) {
	var m map[string]string
	c.Bind(&m)
	var uri string
	if os.Getenv("DEV") == "true" {
		uri = "http://localhost:3000/callback"
	} else {
		uri = "https://epitaf.aureleoules.com/callback"
	}

	resp, err := http.PostForm(Endpoint+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {m["code"]},
		"client_id":     {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"redirect_uri":  {uri},
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil, jwt.ErrFailedAuthentication
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	var result map[string]string
	json.Unmarshal([]byte(body), &result)

	token := result["access_token"]
	if token == "" {
		return nil, jwt.ErrFailedAuthentication
	}

	profile, err := getProfile(token)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	u, err := models.GetUserByEmail(profile.Mail)
	if err != nil {
		user := models.User{
			Name:  profile.DisplayName,
			Email: profile.Mail,
			// TODO: req on CRI
		}
		err = user.Insert()
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		return user, nil
	}

	return models.User{
		Email:     u.Email,
		Name:      u.Name,
		Class:     u.Class,
		Promotion: u.Promotion,
	}, nil
}

func getProfile(token string) (models.MicrosoftProfile, error) {
	endpoint := "https://graph.microsoft.com/v1.0/me"
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.MicrosoftProfile{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.MicrosoftProfile{}, err
	}

	var result models.MicrosoftProfile
	json.Unmarshal([]byte(body), &result)

	if result.Mail == "" {
		return models.MicrosoftProfile{}, errors.New("invalid token")
	}

	return result, nil
}
