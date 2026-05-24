package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"task-manager/controllers"
	"task-manager/database"
	"task-manager/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTaskTestDB(t *testing.T) {
	t.Helper()
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	database.DB.AutoMigrate(&models.User{}, &models.Task{})
}

func createTestUser(t *testing.T) models.User {
	t.Helper()
	hashed, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	user := models.User{Username: fmt.Sprintf("user_%d", time.Now().UnixNano()), Password: string(hashed)}
	database.DB.Create(&user)
	return user
}

func generateTestToken(t *testing.T, userID uint) string {
	t.Helper()
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("test-secret-key-for-testing"))
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}
	return tokenStr
}

func setupTaskRouter(userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	return r
}

func TestCreateTask_Success(t *testing.T) {
	setupTaskTestDB(t)
	user := createTestUser(t)
	r := setupTaskRouter(user.ID)
	r.POST("/tasks", controllers.CreateTask)

	body := `{"title":"Test Task","description":"A test task","priority":2}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var task models.Task
	json.Unmarshal(w.Body.Bytes(), &task)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, user.ID, task.UserID)
}

func TestCreateTask_MissingTitle(t *testing.T) {
	setupTaskTestDB(t)
	user := createTestUser(t)
	r := setupTaskRouter(user.ID)
	r.POST("/tasks", controllers.CreateTask)

	body := `{"description":"No title here"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTasks_ReturnsOnlyOwn(t *testing.T) {
	setupTaskTestDB(t)
	user1 := createTestUser(t)
	user2 := createTestUser(t)

	database.DB.Create(&models.Task{Title: "User1 Task", UserID: user1.ID})
	database.DB.Create(&models.Task{Title: "User2 Task", UserID: user2.ID})

	r := setupTaskRouter(user1.ID)
	r.GET("/tasks", controllers.GetTasks)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var tasks []models.Task
	json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "User1 Task", tasks[0].Title)
}

func TestGetTaskByID_Forbidden(t *testing.T) {
	setupTaskTestDB(t)
	user1 := createTestUser(t)
	user2 := createTestUser(t)

	task := models.Task{Title: "User2 Task", UserID: user2.ID}
	database.DB.Create(&task)

	r := setupTaskRouter(user1.ID)
	r.GET("/tasks/:id", controllers.GetTaskByID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", task.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateTaskStatus_Success(t *testing.T) {
	setupTaskTestDB(t)
	user := createTestUser(t)

	task := models.Task{Title: "Status Test", UserID: user.ID, Status: "pending"}
	database.DB.Create(&task)

	r := setupTaskRouter(user.ID)
	r.PATCH("/tasks/:id/status", controllers.UpdateTaskStatus)

	body := `{"status":"in_progress"}`
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/tasks/%d/status", task.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
