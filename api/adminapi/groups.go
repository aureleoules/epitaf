package adminapi

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/aureleoules/epitaf/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func handleGroups() {
	router.GET("/groups", getGroupsHandler)
	router.GET("/groups/:id", getGroupHandler)
	router.POST("/groups/:id", createSubGroupHandler)

	router.POST("/groups/:id/subjects", addGroupSubjectHandler)
	router.DELETE("/groups/:id/subjects/:subject_id", archiveGroupSubjectHandler)

	router.POST("/groups/:id/users", addGroupUsersHandler)
	router.DELETE("/groups/:id", deleteGroupHandler)
}

func archiveGroupSubjectHandler(c echo.Context) error {
	subjectID, err := models.FromUUID(c.Param("subject_id"))
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	err = models.ArchiveSubject(subjectID)
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func addGroupSubjectHandler(c echo.Context) error {
	groupID, err := models.FromUUID(c.Param("id"))
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	var subject models.Subject
	err = c.Bind(&subject)
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	subject.GroupID = groupID

	err = subject.Insert()
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, subject.ID)
}

func addGroupUsersHandler(c echo.Context) error {
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// userID, err := models.FromUUID(claims["id"].(string))
	// if err != nil {
	// 	return c.JSON(http.StatusNotAcceptable, resp{"error": "invalid jwt"})
	// }

	// r, err := models.GetRealmOfAdmin(userID)
	// if err != nil {
	// 	zap.S().Error(err)
	// 	return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	// }

	groupID, err := models.FromUUID(c.Param("id"))
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	req := struct {
		UserIDs string `json:"user_ids"`
	}{}

	err = c.Bind(&req)
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	var l []models.UUID
	for _, id := range strings.Split(req.UserIDs, ",") {
		i, err := models.FromUUID(id)
		if err != nil {
			zap.S().Warn(err)
			return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
		}

		inGroup, err := models.IsUserInGroup(groupID, i)
		if err != nil {
			zap.S().Error(err)
			return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
		}

		if inGroup {
			zap.S().Warn(err)
			return c.JSON(http.StatusNotAcceptable, resp{"error": "user " + id + " is already in group"})
		}

		l = append(l, i)
	}

	err = models.AddGroupUsers(groupID, l)
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func deleteGroupHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := models.FromUUID(claims["id"].(string))
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": "invalid jwt"})
	}

	r, err := models.GetRealmOfAdmin(userID)
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	id, err := models.FromUUID(c.Param("id"))
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	g, err := models.GetGroup(r.ID, id)
	if err != nil {
		zap.S().Error(err)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, resp{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	if g.ParentID == nil {
		return c.JSON(http.StatusUnauthorized, resp{"error": "cannot delete root group"})
	}

	err = models.DeleteGroup(r.ID, id)
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)

}

func createSubGroupHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := models.FromUUID(claims["id"].(string))
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": "invalid jwt"})
	}

	r, err := models.GetRealmOfAdmin(userID)
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	parentID, err := models.FromUUID(c.Param("id"))
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	var group models.Group
	err = c.Bind(&group)
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	group.ParentID = &parentID
	group.RealmID = r.ID
	group.Slug = slug.Make(group.Name)

	err = group.Insert()
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, group.ID)
}

func getGroupHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := models.FromUUID(claims["id"].(string))
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": "invalid jwt"})
	}

	r, err := models.GetRealmOfAdmin(userID)
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	groupID, err := models.FromUUID(c.Param("id"))
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	group, err := models.GetGroup(r.ID, groupID)
	if err != nil {
		zap.S().Error(err)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, resp{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, group)
}

func getGroupsHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, err := models.FromUUID(claims["id"].(string))
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, resp{"error": "invalid jwt"})
	}

	r, err := models.GetRealmOfAdmin(userID)
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	groups, err := models.GetGroupTree(r.ID)
	if err != nil {
		zap.S().Error(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, groups)
}
