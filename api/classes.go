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
		if formatted[c.Promotion] == nil {
			formatted[c.Promotion] = make(semester)
		}
		if formatted[c.Promotion][c.Semester] == nil {
			formatted[c.Promotion][c.Semester] = make(region)
		}
		formatted[c.Promotion][c.Semester][c.Region] = append(formatted[c.Promotion][c.Semester][c.Region], c.Class)
	}

	c.JSON(http.StatusOK, formatted)
}
