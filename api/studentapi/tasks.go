package studentapi

import "github.com/gin-gonic/gin"

func handleTasks() {
	router.POST("/tasks", createTaskHandler)
	router.GET("/tasks/:id", getTaskHandler)
	router.GET("/tasks", getTasksHandler)
	router.PUT("/tasks/:id", editTaskHandler)
	router.DELETE("/tasks/:id", deleteTaskHandler)

	router.POST("/tasks/:id/complete", completeTaskHandler)
	router.DELETE("/tasks/:id/complete", unCompleteTaskHandler)
}

func createTaskHandler(c *gin.Context) {

}

func getTaskHandler(c *gin.Context) {

}

func getTasksHandler(c *gin.Context) {

}

func editTaskHandler(c *gin.Context) {

}

func deleteTaskHandler(c *gin.Context) {

}

func completeTaskHandler(c *gin.Context) {

}

func unCompleteTaskHandler(c *gin.Context) {

}
