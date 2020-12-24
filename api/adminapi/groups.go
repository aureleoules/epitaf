package adminapi

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
)

func handleGroups() {
	router.GET("/groups", getGroupsHandler)
	router.GET("/groups/:uuid", getGroupHandler)
}

func getGroupHandler(c *gin.Context) {
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

	groupID, err := models.FromUUID(c.Param("uuid"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	group, err := models.GetGroup(r.UUID, groupID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, group)
}

func getGroupsHandler(c *gin.Context) {
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

	groups, err := models.GetGroupTree(r.UUID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO : build tree from flat architecture server-side
	c.JSON(http.StatusOK, groups)
}
