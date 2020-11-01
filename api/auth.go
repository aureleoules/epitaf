package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// Endpoint api url
	Endpoint = "https://login.microsoftonline.com/3534b3d7-316c-4bc9-9ede-605c860f49d2/oauth2/v2.0"
)

func handleAuth() {
	users := api.Group("/users")

	users.POST("/authenticate", authenticateHandler)
	users.POST("/callback", auth.LoginHandler)
}

// AuthMiddleware handles JWT authentications
func AuthMiddleware() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "epitaf",
		Key:        []byte(os.Getenv("JWT_SECRET")),
		Timeout:    time.Hour * 48,
		MaxRefresh: time.Hour * 48,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			u := data.(*models.User)
			// What we put in the JWT claims
			return jwt.MapClaims{
				"uuid":      u.UUID.String(),
				"email":     u.Email,
				"name":      u.Name,
				"promotion": u.Promotion,
				"class":     u.Class,
				"region":    u.Region,
				"semester":  u.Semester,
				"teacher":   u.Teacher,
			}
		},
		Authenticator: callbackHandler,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.AbortWithError(code, errors.New(message))
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	// Authmiddleware must be active
	if err != nil {
		panic(err)
	}

	return authMiddleware
}

func authenticateHandler(c *gin.Context) {
	// Prepare microsoft query
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
	// Return URL that the user must go to
	c.JSON(http.StatusOK, req.URL.String())
}

func getAccessToken(code string) (string, error) {

	// Prepare microsoft query
	var uri string
	if os.Getenv("DEV") == "true" {
		uri = "http://localhost:3000/callback"
	} else {
		uri = "https://epitaf.aureleoules.com/callback"
	}

	resp, err := http.PostForm(Endpoint+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"redirect_uri":  {uri},
	})

	if err != nil {
		zap.S().Error(err)
		return "", jwt.ErrFailedAuthentication
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.S().Error(err)
		return "", jwt.ErrFailedAuthentication
	}

	var result map[string]string
	json.Unmarshal([]byte(body), &result)

	// Return access token
	token := result["access_token"]
	if token == "" {
		zap.S().Error("no access token")
		return "", jwt.ErrFailedAuthentication
	}
	return token, nil
}

func callbackHandler(c *gin.Context) (interface{}, error) {
	var m map[string]string
	c.Bind(&m)
	token, err := getAccessToken(m["code"])
	if err != nil {
		return nil, err
	}

	// Retrieve microsoft profile
	profile, err := getProfile(token)
	if err != nil {
		zap.S().Error(err)
		return nil, jwt.ErrFailedAuthentication
	}

	// Check if user exists in database
	u, err := models.GetUserByEmail(profile.Mail)
	if err != nil {
		// If the user does not exists, we must create a new one using the CRI.
		user, err := models.PrepareUser(profile.Mail)
		if err != nil {
			zap.S().Error(err)
			return nil, err
		}

		// Insert new user and return user data
		err = user.Insert()
		if err != nil {
			zap.S().Error(err)
			return nil, jwt.ErrFailedAuthentication
		}

		return &user, nil
	}

	// User already exists, return user data
	return u, nil
}

func getProfile(token string) (models.MicrosoftProfile, error) {
	zap.S().Info("Fetching Microsoft profile...")
	// Retrieve microsoft profile using access token
	endpoint := "https://graph.microsoft.com/v1.0/me"
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
		return models.MicrosoftProfile{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.S().Error(err)
		return models.MicrosoftProfile{}, err
	}

	var result models.MicrosoftProfile
	json.Unmarshal([]byte(body), &result)

	if result.Mail == "" {
		return models.MicrosoftProfile{}, errors.New("invalid token")
	}

	zap.S().Info("Fetched Microsoft profile...")
	return result, nil
}
