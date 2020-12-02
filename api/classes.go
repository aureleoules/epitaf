package api

import (
	"net/http"

	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
)

func handleClasses() {
	router.GET("/classes", getClassesHandler)
}

// @Summary Get classes
// @Tags classes
// @Description Get list of all registered classes
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 500 "Server error"
// @Router /classes [GET]
func getClassesHandler(c *gin.Context) {
	classes, err := models.GetClasses()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type region map[string][]string
	type semester map[string]region
	type f map[int]semester

	var formatted f = make(f)

	for _, c := range classes {
		if formatted[int(c.Promotion.Int64Value())] == nil {
			formatted[int(c.Promotion.Int64Value())] = make(semester)
		}
		if formatted[int(c.Promotion.Int64Value())][c.Semester.String()] == nil {
			formatted[int(c.Promotion.Int64Value())][c.Semester.String()] = make(region)
		}
		formatted[int(c.Promotion.Int64Value())][c.Semester.String()][c.Region.String()] = append(formatted[int(c.Promotion.Int64Value())][c.Semester.String()][c.Region.String()], c.Class.String())
	}

	c.JSON(http.StatusOK, formatted)
}
