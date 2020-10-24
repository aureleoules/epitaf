package api

import "github.com/gin-gonic/gin"

func handleTasks() {
	api.POST("/tasks", createTaskHandler)
	api.GET("/tasks", getTasksHandler)
}

func getTasksHandler(c *gin.Context) {

}

func createTaskHandler(c *gin.Context) {

}
