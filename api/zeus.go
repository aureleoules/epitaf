package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aureleoules/epitaf/lib/zeus"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleZeus() {
	users := router.Group("/zeus")

	users.GET("/feed/:slug", getZeusFeed)
}

// @Summary Get Zeus ICS feed
// @Tags zeus
// @Description Get Zeus ICS feed
// @Param slug query string true "slug"
// @Success 200	"OK"
// @Failure 406	"Not acceptable"
// @Router /zeus/feed/{slug} [GET]
func getZeusFeed(c *gin.Context) {
	client := zeus.NewClient(nil)

	fmt.Println(c.Param("slug"))
	id, ok := zeus.Groups[c.Param("slug")]
	if !ok {
		c.JSON(http.StatusNotAcceptable, "invalid group slug")
		return
	}

	ics, err := client.GetICS(strconv.Itoa(id))
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusNotAcceptable, err)
		return
	}
	c.String(http.StatusOK, ics)

}
