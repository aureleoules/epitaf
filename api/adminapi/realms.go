package adminapi

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
)

func handleRealms() {
	router.GET("/realms", getCurrentRealmHandler)
	router.GET("/realms/:slug", getRealmHandler)
}

func getCurrentRealmHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	realm, err := models.GetRealmOfUser(userID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, realm)
}

func getRealmHandler(c *gin.Context) {

}
