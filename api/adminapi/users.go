package adminapi

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
)

func handleUsers() {
	router.GET("/users", getUsersHandler)
}

func getUsersHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	r, err := models.GetRealmOfUser(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	users, err := models.GetRealmUsers(r.UUID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, users)
}
