package api

import (
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/lib/zeus"
	"github.com/aureleoules/epitaf/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleUsers() {
	users := router.Group("/users")
	users.POST("/refresh", auth.RefreshHandler)

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

	client := zeus.NewClient(nil)

	// Class mapping
	// TODO clean
	var id int
	if strings.HasPrefix(u.Semester.String(), "S1") || strings.HasPrefix(u.Semester.String(), "S2") {
		id = zeus.Groups["sup-"+strings.ToLower(u.Class.String())]
	} else if strings.HasPrefix(u.Semester.String(), "S5") || strings.HasPrefix(u.Semester.String(), "S6") {
		id = zeus.Groups["ing1-"+strings.ToLower(u.Class.String())]
	}

	cal, err := client.GetGroupPlanning(id)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, cal)
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
