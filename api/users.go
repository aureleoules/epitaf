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
	users := router.Group("/users")
	users.GET("/me", getUserHandler)
	users.GET("/calendar", getCalendarHandler)
	users.GET("/search", searchUserHandler)
}

// @Summary Search user
// @Tags users
// @Description Search user by name or login
// @Param   query	query	string	true	"query"
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error" "Server error"
// @Router /users/calendar [GET]
func searchUserHandler(c *gin.Context) {
	q := c.Query("query")
	users, err := models.SearchUser(q)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get calendar
// @Tags users
// @Description Get user calendar
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error" "Server error"
// @Router /users/calendar [GET]
func getCalendarHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO
	// Retrieve Teacher's schedule from chronos?
	if u.Teacher {
		c.JSON(http.StatusOK, nil)
		return
	}

	client := chronos.NewClient(os.Getenv("CHRONOS_TOKEN"), nil)

	// Class mapping
	// TODO clean
	var slug string
	if strings.HasPrefix(u.Semester.String(), "S1") || strings.HasPrefix(u.Semester.String(), "S2") || strings.HasPrefix(u.Semester.String(), "S3") || strings.HasPrefix(u.Semester.String(), "S4") {
		slug = "INFO" + strings.ReplaceAll(u.Semester.String(), "#", "%23") + u.Class.String()
	} else {
		if u.Class.String() == "BING" {
			slug = "BING B"
		} else if strings.HasPrefix(u.Class.String(), "A") {
			slug = "RIEMANN " + u.Class.String()
		} else if strings.HasPrefix(u.Class.String(), "C") {
			slug = "SHANNON " + u.Class.String()
		} else if strings.HasPrefix(u.Class.String(), "D") {
			slug = "TANENBAUM " + u.Class.String()
		}
	}

	cal, err := client.GetGroupPlanning(slug)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, models.FormatCalendar(*cal))
}

// @Summary Get self
// @Tags users
// @Description Retrieve data about current user
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error"
// @Router /users/me [GET]
func getUserHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		spew.Dump(claims)
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, u)
}
