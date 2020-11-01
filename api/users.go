package api

import (
	"net/http"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/lib/chronos"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleUsers() {
	users := api.Group("/users")
	users.GET("/me", getUserHandler)
	users.GET("/calendar", getCalendarHandler)
}

func getCalendarHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uuid, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	u, err := models.GetUser(uuid)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := chronos.NewClient(os.Getenv("CHRONOS_TOKEN"), nil)
	slug := "INFO" + u.Semester + u.Class
	cal, err := client.GetGroupPlanning(slug)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, cal)
}
func getUserHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uuid, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	u, err := models.GetUser(uuid)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, u)
}
