package adminapi

import (
	"net/http"

	"github.com/aureleoules/epitaf/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func handleUsers() {
	router.GET("/users", getUsersHandler)
}

func getUsersHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := models.FromUUID(claims["id"].(string))
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": "invalid jwt"})
	}

	r, err := models.GetRealmOfAdmin(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	users, err := models.GetRealmUsers(r.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}
