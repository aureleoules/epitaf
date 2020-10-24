package api

import (
	"fmt"
	"net/http"
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
}

func getTaskHandler(c *gin.Context) {
	task, err := models.GetTask(c.Param("id"))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, task)
}

func getTasksHandler(c *gin.Context) {
	tasks, err := models.GetTasksRange(time.Now(), time.Now().Add(time.Hour*24*7*2))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func createTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	spew.Dump(claims)
	fmt.Println("UUID", claims["uuid"].(string))
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
	task.Promotion = u.Promotion
	if !task.Global {
		task.Class = u.Class
	}

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
