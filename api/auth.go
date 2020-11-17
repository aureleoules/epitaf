package api

import (
	"errors"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/lib/microsoft"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleAuth() {
	users := router.Group("/users")

	users.POST("/authenticate", authenticateHandler)
	users.POST("/callback", auth.LoginHandler)
}

func authenticateHandler(c *gin.Context) {
	var data struct {
		RedirectURI string `json:"redirect_uri"`
	}
	err := c.BindJSON(&data)
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	c.JSON(http.StatusOK, microsoft.SignInURL(data.RedirectURI))
}

func callbackHandler(c *gin.Context) (interface{}, error) {
	var m map[string]string
	err := c.Bind(&m)
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return nil, jwt.ErrFailedAuthentication
	}
	if m["code"] == "" {
		_ = c.AbortWithError(http.StatusNotAcceptable, errors.New("missing code"))
		return nil, jwt.ErrFailedAuthentication
	}
	if m["redirect_uri"] == "" {
		_ = c.AbortWithError(http.StatusNotAcceptable, errors.New("missing code"))
		return nil, jwt.ErrFailedAuthentication
	}

	token, err := microsoft.GetAccessToken(m["code"], m["redirect_uri"])
	if err != nil {
		return nil, err
	}

	client := microsoft.NewClient(token, nil)
	// Retrieve microsoft profile
	profile, err := client.GetProfile()
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
