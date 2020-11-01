package api

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleTasks() {
	api.POST("/tasks", createTaskHandler)
	api.GET("/tasks/:id", getTaskHandler)
	api.GET("/tasks", getTasksHandler)
	api.PUT("/tasks/:id", editTaskHandler)
	api.DELETE("/tasks/:id", deleteTaskHandler)
}

func editTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uuid, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	u, err := models.GetUser(uuid)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task, err := models.GetTask(c.Param("id"))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// Check if user is authorized
	if (task.Global && task.Promotion != u.Promotion) || (!task.Global && (u.Region != task.Region || u.Class != task.Class || u.Promotion != task.Promotion || u.Semester != task.Semester)) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var t models.Task
	err = c.BindJSON(&t)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	t.UpdatedByID = u.UUID

	err = models.UpdateTask(t)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func deleteTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uuid, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	u, err := models.GetUser(uuid)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task, err := models.GetTask(c.Param("id"))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// Check if user is authorized
	if (task.Global && task.Promotion != u.Promotion) || (!task.Global && (u.Region != task.Region || u.Class != task.Class || u.Promotion != task.Promotion || u.Semester != task.Semester)) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = models.DeleteTask(task.ShortID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}
func getTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uuid, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	u, err := models.GetUser(uuid)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task, err := models.GetTask(c.Param("id"))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Check if user is authorized
	if (task.Global && task.Promotion != u.Promotion) || (!task.Global && (u.Region != task.Region || u.Class != task.Class || u.Promotion != task.Promotion || u.Semester != task.Semester)) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, task)
}

func getTasksHandler(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	uuid, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	u, err := models.GetUser(uuid)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tasks, err := models.GetTasksRange(u.Promotion, u.Semester, u.Class, u.Region, time.Now().Add(-time.Hour*24), time.Now().Add(time.Hour*24*7*2))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func createTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uuid, err := models.FromUUID(claims["uuid"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	var task models.Task

	err = c.BindJSON(&task)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	u, err := models.GetUser(uuid)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task.CreatedByID = u.UUID
	task.UpdatedByID = u.UUID
	task.Region = u.Region
	task.Promotion = u.Promotion
	task.Class = u.Class
	task.Semester = u.Semester

	err = task.Validate()
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	err = task.Insert()
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, task.ShortID)
}
