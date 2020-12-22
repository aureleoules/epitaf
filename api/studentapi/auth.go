package studentapi

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func handleAuth() {
	router.POST("/auth", registerHandler)
	router.POST("/auth/login", auth.LoginHandler)
}

func registerHandler(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		zap.S().Warn(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	err = user.Validate()
	if err != nil {
		zap.S().Warn(err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, err)
		return
	}

	user.HashPassword()

	err = user.Insert()
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func authenticator(c *gin.Context) (interface{}, error) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		zap.S().Warn(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return nil, jwt.ErrFailedAuthentication
	}

	u, err := models.GetUserByEmail(user.Email)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &models.User{
		UUID:    u.UUID,
		Login:   u.Login,
		RealmID: u.RealmID,
		Email:   u.Email,
	}, nil
}
