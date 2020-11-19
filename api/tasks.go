package api

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/aureleoules/epitaf/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleTasks() {
	router.POST("/tasks", createTaskHandler)
	router.GET("/tasks/:id", getTaskHandler)
	router.GET("/tasks", getTasksHandler)
	router.PUT("/tasks/:id", editTaskHandler)
	router.DELETE("/tasks/:id", deleteTaskHandler)
}

func editTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
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

	var t models.Task
	err = c.BindJSON(&t)
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if task.CreatedByLogin != u.Login && task.Visibility != t.Visibility {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if task.CreatedByLogin == u.Login ||
		(task.Visibility == models.StudentsVisibility && task.Members.Includes(u.Login)) ||
		(u.Teacher && (task.Visibility == models.ClassVisibility || task.Visibility == models.PromotionVisibility)) ||
		(task.Visibility == models.PromotionVisibility && u.Promotion == task.Promotion && u.Semester == task.Semester) ||
		(task.Visibility == models.ClassVisibility && u.Promotion == task.Promotion && u.Semester == task.Semester && task.Class == u.Class && task.Region == u.Region) {

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
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)

}

func deleteTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
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

	// Only author can delete, or teacher is task is promo or class
	if task.CreatedByLogin == u.Login || ((task.Visibility == models.ClassVisibility || task.Visibility == models.PromotionVisibility) && u.Teacher) {
		err = models.DeleteTask(task.ShortID)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		zap.S().Info("User ", u.Name, " deleted task ", task.ShortID)
		c.Status(http.StatusOK)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func getTaskHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, err := models.GetUser(claims["login"].(string))
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

	if task.CreatedByLogin == u.Login ||
		(task.Visibility == models.StudentsVisibility && task.Members.Includes(u.Login)) ||
		(u.Teacher && (task.Visibility == models.ClassVisibility || task.Visibility == models.PromotionVisibility)) ||
		(task.Visibility == models.PromotionVisibility && u.Promotion == task.Promotion && u.Semester == task.Semester) ||
		(task.Visibility == models.ClassVisibility && u.Promotion == task.Promotion && u.Semester == task.Semester && task.Class == u.Class && task.Region == u.Region) {

		c.JSON(http.StatusOK, task)
		return
	}

	c.AbortWithStatus(http.StatusNotFound)
}

func getTasksHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	u, err := models.GetUser(claims["login"].(string))
	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var tasks []models.Task

	start := utils.TruncateDate(time.Now())
	// TODO: client chosen time ranges
	end := utils.TruncateDate(time.Now().Add(time.Hour * 24 * 365))
	if u.Teacher {
		tasks, err = models.GetTeacherTasksRange(start, end)
	} else {
		tasks, err = models.GetTasksRange(*u, start, end)
	}

	if err != nil {
		zap.S().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	zap.S().Info(u.Name + " fetched tasks.")
	c.JSON(http.StatusOK, tasks)
}

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
		DueDate:        t.DueDate,
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
