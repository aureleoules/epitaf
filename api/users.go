package api

import (
	"net/http"
	"os"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/lib/chronos"
	"github.com/aureleoules/epitaf/models"
	"github.com/davecgh/go-spew/spew"
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

	// Class mapping
	// TODO clean
	var slug string
	if strings.HasPrefix(u.Semester, "S1") || strings.HasPrefix(u.Semester, "S2") || strings.HasPrefix(u.Semester, "S3") || strings.HasPrefix(u.Semester, "S4") {
		slug += "INFO" + u.Semester
		if strings.Contains(u.Class, "#") {
			slug += "#" + strings.Replace(u.Class, "#", "", -1)
		} else {
			slug = "INFO" + u.Semester + u.Class
		}
	} else {
		if u.Class == "BING" {
			slug = "BING B"
		} else if strings.HasPrefix(u.Class, "A") {
			slug = "RIEMANN " + u.Class
		} else if strings.HasPrefix(u.Class, "C") {
			slug = "SHANNON " + u.Class
		} else if strings.HasPrefix(u.Class, "D") {
			slug = "TANENBAUM " + u.Class
		}
	}

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
		spew.Dump(claims)
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, u)
}
