package api

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleTasks() {
	router.POST("/tasks", createTaskHandler)
	router.GET("/tasks/:id", getTaskHandler)
	router.GET("/tasks", getTasksHandler)
	router.PUT("/tasks/:id", editTaskHandler)
	router.DELETE("/tasks/:id", deleteTaskHandler)

	router.POST("/tasks/:id/complete", completeTaskHandler)
	router.DELETE("/tasks/:id/complete", unCompleteTaskHandler)
}

// @Summary Complete task
// @Tags tasks
// @Description Mark a specific task as completed
// @Param   short_id	path	string	true	"short_id"
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error"
// @Router /tasks/{short_id}/complete [POST]
func completeTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task, err := models.GetUserTask(c.Param("id"), u.Login)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if task.Completed {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if !u.CanViewTask(*task) {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = task.Mark(u.Login)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Uncomplete task
// @Tags tasks
// @Description Mark a specific task as uncompleted
// @Param   short_id	path	string	true	"short_id"
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error"
// @Router /tasks/{short_id}/complete [DELETE]
func unCompleteTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task, err := models.GetUserTask(c.Param("id"), u.Login)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if !task.Completed {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if !u.CanViewTask(*task) {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = task.Unmark(u.Login)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Update task
// @Tags tasks
// @Description Edit a specific task
// @Accept  json
// @Param   short_id	path	string	true	"short_id"
// @Param 	task		body	models.Task	true "Task"
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error"
// @Router /tasks/{short_id} [PUT]
func editTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task, err := models.GetUserTask(c.Param("id"), u.Login)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var t models.Task
	err = c.BindJSON(&t)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if task.CreatedByLogin != u.Login && t.Visibility != "" && task.Visibility != t.Visibility {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !u.CanEditTask(*task) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	update := t.PrepareUpdate(*task, *u)
	update.UpdatedByLogin = u.Login
	update.ShortID = task.ShortID

	err = update.Validate()
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	err = models.UpdateTask(update)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	zap.S().Info("User ", u.Name, " updated task ", update.ShortID)
	c.Status(http.StatusOK)
}

// @Summary Delete task
// @Tags tasks
// @Description Delete a specific task
// @Param   short_id	path	string	true	"short_id"
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error"
// @Router /tasks/{short_id} [DELETE]
func deleteTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task, err := models.GetUserTask(c.Param("id"), u.Login)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Only author can delete, or teacher is task is promo or class
	if !u.CanDeleteTask(*task) {
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

// @Summary Get task
// @Tags tasks
// @Description Get a specific task
// @Param   short_id	path	string	true	"short_id"
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error" "Server error"
// @Router /tasks/{short_id} [GET]
func getTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task, err := models.GetUserTask(c.Param("id"), u.Login)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if !u.CanViewTask(*task) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary Get tasks
// @Tags tasks
// @Description Get tasks
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error" "Server error"
// @Router /tasks [GET]
func getTasksHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var tasks []models.Task

	var query models.Filters
	err = c.BindQuery(&query)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	err = query.Validate()
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if u.Teacher {
		tasks, err = models.GetTeacherTasksRange(query.StartDate, query.EndDate)
	} else {
		if u.Login != "api_key" {
			tasks, err = models.GetTasksRange(*u, query)
		} else {
			tasks, err = models.GetAllTasks(query.StartDate, query.EndDate)
		}
	}

	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	zap.S().Info(u.Name + " fetched tasks.")
	c.JSON(http.StatusOK, tasks)
}

// @Summary Create task
// @Tags tasks
// @Description Create a new task
// @Accept  json
// @Param   task body    models.Task     true        "Task"
// @Success 200	"OK"
// @Failure 401	"Unauthorized"
// @Failure 404	"Not found"
// @Failure 406	"Not acceptable"
// @Failure 500 "Server error"
// @Router /tasks [POST]
func createTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	var t models.Task
	err := c.BindJSON(&t)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task := models.Task{
		Content:        t.Content,
		DueDate:        t.DueDate.Local(),
		Subject:        t.Subject,
		Title:          t.Title,
		Visibility:     t.Visibility,
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
	}

	task.CreatedByLogin = u.Login
	task.UpdatedByLogin = u.Login

	// If user is a student
	// Retrieve class & promo from user data
	// Prevents classes from adding tasks to other classes
	if !u.Teacher {
		if t.Visibility == models.ClassVisibility {
			task.Region = u.Region
			task.Promotion = u.Promotion
			task.Class = u.Class
			task.Semester = u.Semester
		}
		if t.Visibility == models.PromotionVisibility {
			task.Promotion = u.Promotion
			task.Semester = u.Semester
		}
	} else {
		// If user is teacher, set task class & class to input
		if t.Visibility == models.ClassVisibility {
			task.Region = t.Region
			task.Promotion = t.Promotion
			task.Class = t.Class
			task.Semester = t.Semester
		}
		if t.Visibility == models.PromotionVisibility {
			task.Promotion = t.Promotion
			task.Semester = t.Semester
		}
	}

	if t.Visibility == models.StudentsVisibility {
		task.Members = t.Members
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

	zap.S().Info("User ", u.Name, " created task ", task.ShortID)
	c.JSON(http.StatusOK, task.ShortID)
}
