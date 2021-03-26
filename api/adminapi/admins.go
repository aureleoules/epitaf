package adminapi

import (
	"net/http"

	"github.com/aureleoules/epitaf/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func handleAdmins() {
	router.GET("/admins/me", getSelfAdminHandler)
}

func getSelfAdminHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := models.FromUUID(claims["id"].(string))
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": "invalid jwt"})
	}

	u, err := models.GetAdmin(userID)
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, u)
}
