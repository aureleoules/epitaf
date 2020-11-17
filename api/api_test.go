package api

import (
	"testing"

	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	godotenv.Load("../.env_test")
	gin.SetMode(gin.TestMode)
}

func refreshDB() {
	db.Connect()
	db.Delete()
	db.Init()
	err := models.InjectSQLSchemas()
	if err != nil {
		zap.S().Fatal(err)
	}
}

func insertTestUser2024C1() (models.User, string) {
	u := models.User{
		Name:      "Test C1",
		Login:     "test_user_2024_C1",
		Class:     "C1",
		Email:     "test_user_2024_C1@epita.fr",
		Promotion: 2024,
		Region:    "Paris",
		Semester:  "S3",
		Teacher:   false,
	}
	u.Insert()

	token, _, _ := AuthMiddleware().TokenGenerator(&u)
	return u, token
}

func insertTestUser2024C2() (models.User, string) {
	u := models.User{
		Name:      "Test C2",
		Login:     "test_user_2024_C2",
		Class:     "C2",
		Email:     "test_user_2024_C2@epita.fr",
		Promotion: 2024,
		Region:    "Paris",
		Semester:  "S3",
		Teacher:   false,
	}
	u.Insert()

	token, _, _ := AuthMiddleware().TokenGenerator(&u)
	return u, token
}

func insertTestUser2025C1() (models.User, string) {
	u := models.User{
		Name:      "Test 2025 C1",
		Login:     "test_user_2025_C1",
		Class:     "C1",
		Email:     "test_user_2025_C1@epita.fr",
		Promotion: 2025,
		Region:    "Paris",
		Semester:  "S3",
		Teacher:   false,
	}
	u.Insert()

	token, _, _ := AuthMiddleware().TokenGenerator(&u)
	return u, token
}

func insertTestUser2024Lyon() (models.User, string) {
	u := models.User{
		Name:      "Test 2024 Lyon",
		Login:     "test_user_2024_lyon",
		Class:     "L1",
		Email:     "test_user_2024_lyon@epita.fr",
		Promotion: 2024,
		Region:    "Lyon",
		Semester:  "S3",
		Teacher:   false,
	}
	u.Insert()

	token, _, _ := AuthMiddleware().TokenGenerator(&u)
	return u, token
}

func insertTestTeacher() (models.User, string) {
	u := models.User{
		Name:    "Teacher",
		Login:   "test_teacher",
		Email:   "teacher@epita.fr",
		Teacher: true,
	}
	u.Insert()

	token, _, _ := AuthMiddleware().TokenGenerator(&u)
	return u, token
}

func TestServe(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serve()
		})
	}
}