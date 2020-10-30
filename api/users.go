package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
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

	slug := "INFO" + u.Semester + u.Class
	cal, err := getCalendar(slug)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, cal)
}

func getCalendar(class string) (models.ChronosCalendar, error) {
	endpoint := "https://v2ssl.webservices.chronos.epita.net/api/v2/Planning/GetRangeWeekRecursive/" + class + "/0"
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Auth-Token", os.Getenv("CHRONOS_TOKEN"))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.ChronosCalendar{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.ChronosCalendar{}, err
	}

	var result []models.ChronosCalendar
	json.Unmarshal([]byte(body), &result)

	return result[0], nil
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
