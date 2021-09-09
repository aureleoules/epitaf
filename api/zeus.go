package api

import (
	"net/http"

	"github.com/aureleoules/epitaf/lib/zeus"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleZeus() {
	users := router.Group("/zeus")

	users.GET("/feed/:group_id", getZeusFeed)
}

// @Summary Get Zeus ICS feed
// @Tags zeus
// @Description Get Zeus ICS feed
// @Param group_id body string true "group_id"
// @Success 200	"OK"
// @Failure 406	"Not acceptable"
// @Router /zeus/feed/{group_id} [GET]
func getZeusFeed(c *gin.Context) {
	client := zeus.NewClient(nil)

	ics, err := client.GetICS(c.Param("group_id"))
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusNotAcceptable, err)
		return
	}
	c.String(http.StatusOK, ics)

}
