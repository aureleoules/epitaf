package adminapi

import (
	"net/http"

	"github.com/aureleoules/epitaf/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/mattn/go-nulltype"
)

func handleUsers() {
	router.GET("/users", getUsersHandler)
	router.POST("/users", createUserHandler)
}

func createUserHandler(c echo.Context) error {
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

	var u models.User
	err = c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	if u.Password.String() == "" {
		u.Password = nulltype.NullString{}
	}

	err = u.Validate()
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	u.RealmID = r.ID

	if u.Password.Valid() {
		u.HashPassword()
	}

	err = u.Insert()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, u.ID)
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

	var filters models.UserFilters
	err = c.Bind(&filters)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	spew.Dump(filters)

	filters.Defaults()

	users, err := models.GetRealmUsers(r.ID, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}
