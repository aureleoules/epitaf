package adminapi

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func handleAuth() {
	router.POST("/users", registerHandler)
	router.POST("/users/login", auth.LoginHandler)
}

func registerHandler(c *gin.Context) {
	req := struct {
		Login    string `json:"login"`
		Name     string `json:"name" `
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.BindJSON(&req)
	if err != nil {
		zap.S().Warn(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	admin := models.Admin{
		Login:    req.Login,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err = admin.Validate()
	if err != nil {
		zap.S().Warn(err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, err)
		return
	}

	admin.HashPassword()

	err = admin.Insert()
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func authenticator(c *gin.Context) (interface{}, error) {
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := c.BindJSON(&req)
	if err != nil {
		zap.S().Warn(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return nil, jwt.ErrFailedAuthentication
	}

	u, err := models.GetAdminByEmail(req.Username)
	if err != nil {
		fmt.Println(err)
		return nil, jwt.ErrFailedAuthentication
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
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
