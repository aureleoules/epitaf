package adminapi

import (
	"net/http"

	"github.com/aureleoules/epitaf/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func handleTasks() {
	router.GET("/tasks", getTasksHandler)
}

func getTasksHandler(c echo.Context) error {
	return nil
}

func handleCreateTask(c echo.Context) error {
	groupID, err := models.FromUUID(c.Param("id"))
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	var task models.Task
	err = c.Bind(&task)
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})
	}

	task.GroupID = groupID

	err = task.Validate()
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": err.Error()})

	}

	err = task.Insert()
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, task.ID)
}
