package api

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/davecgh/go-spew/spew"
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
	if !u.Teacher && ((task.Global && task.Promotion != u.Promotion) || (!task.Global && (u.Region != task.Region || u.Class != task.Class || u.Promotion != task.Promotion || u.Semester != task.Semester))) {
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

	if !u.Teacher {
		t.Class = task.Class
		t.Semester = task.Semester
		t.Promotion = task.Promotion
		t.Region = task.Region
	}

	t.Semester = strings.ToUpper(t.Semester)
	t.Class = strings.ToUpper(t.Class)
	t.Region = strings.Title(strings.ToLower(t.Region))

	err = models.UpdateTask(t)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	zap.S().Info("User ", u.Name, " updated task ", task.ShortID)
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
	if !u.Teacher && ((task.Global && task.Promotion != u.Promotion) || (!task.Global && (u.Region != task.Region || u.Class != task.Class || u.Promotion != task.Promotion || u.Semester != task.Semester))) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = models.DeleteTask(task.ShortID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	zap.S().Info("User ", u.Name, " deleted task ", task.ShortID)
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
	if !u.Teacher && ((task.Global && task.Promotion != u.Promotion) || (!task.Global && (u.Region != task.Region || u.Class != task.Class || u.Promotion != task.Promotion || u.Semester != task.Semester))) {
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
		spew.Dump(claims)
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var tasks []models.Task

	start := time.Now().Add(-time.Hour * 24)
	// TODO: client chosen time ranges
	end := time.Now().Add(time.Hour * 24 * 365)
	if u.Teacher {
		tasks, err = models.GetAllTasksRange(start, end)
	} else {
		tasks, err = models.GetTasksRange(u.Promotion, u.Semester, u.Class, u.Region, start, end)
	}

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

	if !u.Teacher {
		task.Region = u.Region
		task.Promotion = u.Promotion
		task.Class = u.Class
		task.Semester = u.Semester
	}

	task.Semester = strings.ToUpper(task.Semester)
	task.Class = strings.ToUpper(task.Class)
	task.Region = strings.Title(strings.ToLower(task.Region))

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

	zap.S().Info("User ", u.Name, " created task ", task.ShortID)

	c.JSON(http.StatusOK, task.ShortID)
}
