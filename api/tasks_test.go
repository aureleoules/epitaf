package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/aureleoules/epitaf/models"
	"github.com/aureleoules/epitaf/utils"
	"github.com/gin-gonic/gin"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func Test_editTaskHandler(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editTaskHandler(tt.args.c)
		})
	}
}

func Test_deleteTaskHandler(t *testing.T) {
	refreshDB()

	u, token := insertTestUser2024C1()
	c1_2024, tokenc1 := insertTestUser2024C1_2()
	c2_2025, token2 := insertTestUser2024C2()
	_, token3 := insertTestUser2025C1()
	teacher, tokenTeacher := insertTestTeacher()

	/* Promotion task */
	task := models.Task{
		Subject:        "mathematics",
		Content:        "This is a test",
		DueDate:        utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility:     models.PromotionVisibility,
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
		Title:          "Thing to do",
		Promotion:      u.Promotion,
		Semester:       u.Semester,
	}
	task.Insert()

	// Check unauthorized
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/" + task.ShortID).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Check not found
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/").
		Header("Authorization", "Bearer "+token2).
		Expect(t).
		Status(http.StatusNotFound).
		End()

	// Check not found
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/abcd1234").
		Header("Authorization", "Bearer "+token2).
		Expect(t).
		Status(http.StatusNotFound).
		End()

	// Other class cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token2).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Other promo cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token3).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Other class member cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenc1).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// User can delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	// Insert task again
	task.Insert()

	// Check ok delete from teacher
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenTeacher).
		Expect(t).
		Status(http.StatusOK).
		End()

	// Should not delete after delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenTeacher).
		Expect(t).
		Status(http.StatusNotFound).
		End()

	/* Self task */
	task = models.Task{
		Subject:        "mathematics",
		Content:        "This is a test",
		DueDate:        utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility:     models.SelfVisibility,
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
		Title:          "Thing to do",
	}
	task.Insert()

	// Other class cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token2).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Other promo cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token3).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Teacher cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenTeacher).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Other class member cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenc1).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// User can delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	/* Class task */
	task = models.Task{
		Subject:        "mathematics",
		Content:        "This is a test",
		DueDate:        utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility:     models.ClassVisibility,
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
		Title:          "Thing to do",
		Class:          u.Class,
		Promotion:      u.Promotion,
		Semester:       u.Semester,
		Region:         u.Region,
	}
	task.Insert()

	// Other class cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token2).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Other promo cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token3).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Same class member cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenc1).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Teacher can delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenTeacher).
		Expect(t).
		Status(http.StatusOK).
		End()
	task.Insert()

	// Author can delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	/* Students task */
	task = models.Task{
		Subject:        "mathematics",
		Content:        "This is a test",
		DueDate:        utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility:     models.StudentsVisibility,
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
		Title:          "Thing to do",
		Members:        []string{c1_2024.Login, c2_2025.Login},
	}
	task.Insert()

	// Other class cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token2).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Other promo cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token3).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Same class member cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenc1).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Teacher cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenTeacher).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()
	task.Insert()

	// Author can delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	/* Students task with teacher */
	task = models.Task{
		Subject:        "mathematics",
		Content:        "This is a test",
		DueDate:        utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility:     models.StudentsVisibility,
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
		Title:          "Thing to do",
		Members:        []string{c1_2024.Login, c2_2025.Login, teacher.Login},
	}
	task.Insert()

	// Teacher cannot delete
	apitest.New().
		Handler(createRouter()).
		Delete("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenTeacher).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

}

func Test_getTaskHandler(t *testing.T) {
	refreshDB()

	u, token := insertTestUser2024C1()
	_, token2 := insertTestUser2024C2()
	_, token3 := insertTestUser2025C1()
	_, tokenTeacher := insertTestTeacher()

	// Insert task
	task := models.Task{
		Subject:        "mathematics",
		Content:        "This is a test",
		DueDate:        utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility:     models.PromotionVisibility,
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
		Title:          "Thing to do",
		Promotion:      u.Promotion,
		Semester:       u.Semester,
	}
	task.Insert()

	// Check unauthorized
	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/" + task.ShortID).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Check not found
	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/1").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusNotFound).
		End()
	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/1234").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusNotFound).
		End()

	r := apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	var result models.Task
	r.JSON(&result)

	assert.Equal(t, task.Subject, result.Subject)
	assert.Equal(t, task.Title, result.Title)
	assert.Equal(t, task.Content, result.Content)
	assert.Equal(t, task.DueDate.Unix(), result.DueDate.Unix())
	assert.Equal(t, u.Login, result.CreatedByLogin)
	assert.Equal(t, u.Login, result.UpdatedByLogin)
	assert.Equal(t, task.Visibility, result.Visibility)
	assert.Equal(t, u.Semester, result.Semester)
	assert.Equal(t, u.Promotion, result.Promotion)

	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token2).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token3).
		Expect(t).
		Status(http.StatusNotFound).
		End()

	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenTeacher).
		Expect(t).
		Status(http.StatusOK).
		End()

	// Insert task
	task = models.Task{
		Subject:        "mathematics",
		Content:        "This is a test",
		DueDate:        utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility:     models.SelfVisibility,
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
		Title:          "Thing to do",
	}
	task.Insert()

	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token2).
		Expect(t).
		Status(http.StatusNotFound).
		End()

	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+token3).
		Expect(t).
		Status(http.StatusNotFound).
		End()

	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks/"+task.ShortID).
		Header("Authorization", "Bearer "+tokenTeacher).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func Test_getTasksHandler(t *testing.T) {
	refreshDB()
	u, token := insertTestUser2024C1()

	// Check unauthorized
	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks").
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Check without any data
	apitest.New().
		Handler(createRouter()).
		Get("/api/tasks").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Body("null").
		Status(http.StatusOK).
		End()

	// Insert task
	task := models.Task{
		Subject:        "mathematics",
		Title:          "Thing to do",
		Content:        "This is a test",
		DueDate:        utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		CreatedByLogin: u.Login,
		UpdatedByLogin: u.Login,
		Visibility:     models.PromotionVisibility,
		Semester:       "S3",
		Promotion:      2024,
	}
	task.Insert()

	// Check if equal
	r := apitest.New().
		Handler(createRouter()).
		Get("/api/tasks").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()
	var tasks []models.Task
	r.JSON(&tasks)

	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, task.Subject, tasks[0].Subject)
	assert.Equal(t, task.Title, tasks[0].Title)
	assert.Equal(t, task.Content, tasks[0].Content)
	assert.Equal(t, task.DueDate.Unix(), tasks[0].DueDate.Unix())
	assert.Equal(t, task.CreatedByLogin, tasks[0].CreatedByLogin)
	assert.Equal(t, task.UpdatedByLogin, tasks[0].UpdatedByLogin)
	assert.Equal(t, task.Visibility, tasks[0].Visibility)
	assert.Equal(t, task.Semester, tasks[0].Semester)
	assert.Equal(t, task.Promotion, tasks[0].Promotion)

	for i := 0; i < 19; i++ {
		task.Insert()
	}

	r = apitest.New().
		Handler(createRouter()).
		Get("/api/tasks").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	r.JSON(&tasks)
	assert.Equal(t, 20, len(tasks))
}

func Test_createTaskHandler(t *testing.T) {
	refreshDB()
	u, token := insertTestUser2024C1()

	// Insert task
	task := models.Task{
		Subject:    "mathematics",
		Content:    "This is a test",
		DueDate:    utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility: models.PromotionVisibility,
		Title:      "Thing to do",
	}
	data, err := json.Marshal(task)
	assert.Nil(t, err)

	// Check unauthorized
	apitest.New().
		Handler(createRouter()).
		Post("/api/tasks").
		Body(string(data)).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	// Check ok
	r := apitest.New().
		Handler(createRouter()).
		Post("/api/tasks").
		Header("Authorization", "Bearer "+token).
		Body(string(data)).
		Expect(t).
		Status(http.StatusOK).
		End()
	var id string
	r.JSON(&id)
	assert.NotEqual(t, "", id)

	ta, err := models.GetTask(id)
	assert.Nil(t, err)

	assert.Equal(t, task.Subject, ta.Subject)
	assert.Equal(t, task.Title, ta.Title)
	assert.Equal(t, task.Content, ta.Content)
	assert.Equal(t, task.DueDate.Unix(), ta.DueDate.Unix())
	assert.Equal(t, u.Login, ta.CreatedByLogin)
	assert.Equal(t, u.Login, ta.UpdatedByLogin)
	assert.Equal(t, task.Visibility, ta.Visibility)
	assert.Equal(t, u.Semester, ta.Semester)
	assert.Equal(t, u.Promotion, ta.Promotion)

	// Insert task
	task = models.Task{
		Subject:    "mathematics",
		Content:    "This is a test",
		DueDate:    utils.TruncateDate(time.Now().Add(-time.Hour * 24)),
		Visibility: models.PromotionVisibility,
		Title:      "Thing to do",
	}
	data, err = json.Marshal(task)
	assert.Nil(t, err)

	// Check impossible due date
	apitest.New().
		Handler(createRouter()).
		Post("/api/tasks").
		Header("Authorization", "Bearer "+token).
		Body(string(data)).
		Expect(t).
		Status(http.StatusNotAcceptable).
		End()

	// Check normal task
	task = models.Task{
		Subject:    "physics",
		Content:    "This is a test",
		DueDate:    utils.TruncateDate(time.Now().Add(time.Hour * 72)),
		Visibility: models.SelfVisibility,
		Title:      "Thing to do",
	}

	data, err = json.Marshal(task)
	assert.Nil(t, err)

	r = apitest.New().
		Handler(createRouter()).
		Post("/api/tasks").
		Header("Authorization", "Bearer "+token).
		Body(string(data)).
		Expect(t).
		Status(http.StatusOK).
		End()

	r.JSON(&id)
	assert.NotEqual(t, "", id)

	ta, err = models.GetTask(id)
	assert.Nil(t, err)

	assert.Equal(t, task.Subject, ta.Subject)
	assert.Equal(t, task.Title, ta.Title)
	assert.Equal(t, task.Content, ta.Content)
	assert.Equal(t, task.DueDate.Unix(), ta.DueDate.Unix())
	assert.Equal(t, u.Login, ta.CreatedByLogin)
	assert.Equal(t, u.Login, ta.UpdatedByLogin)
	assert.Equal(t, task.Visibility, ta.Visibility)
	assert.Equal(t, "", ta.Semester)
	assert.Equal(t, 0, ta.Promotion)
}
