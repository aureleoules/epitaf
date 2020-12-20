package studentapi

import (
	"github.com/gin-gonic/gin"
)

func handleUsers() {
	users := router.Group("/users")
	users.GET("/me", getUserHandler)
}

func getUserHandler(c *gin.Context) {
	// TODO: fix
	// claims := jwt.ExtractClaims(c)
	// u, err := models.GetUser(claims["login"].(string))
	// if err != nil {
	// 	zap.S().Error(err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// c.JSON(http.StatusOK, u)
}
