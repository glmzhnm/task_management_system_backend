package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"task-manager/controllers"
	"task-manager/database"
	"task-manager/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	database.DB.AutoMigrate(&models.User{}, &models.Task{})
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	r := setupRouter()
	r.POST("/auth/register", controllers.Register)

	body := `{"username":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "User registered successfully", resp["message"])
}

func TestRegister_DuplicateUsername(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()
	r.POST("/auth/register", controllers.Register)

	body := `{"username":"dupuser","password":"password123"}`

	req1 := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	req2 := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestRegister_InvalidInput(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()
	r.POST("/auth/register", controllers.Register)

	body := `{"username":"testuser"}` // no password
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	setupTestDB(t)
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	r := setupRouter()
	r.POST("/auth/register", controllers.Register)
	r.POST("/auth/login", controllers.Login)

	regBody := `{"username":"loginuser","password":"password123"}`
	req1 := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(regBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	loginBody := `{"username":"loginuser","password":"password123"}`
	req2 := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
}

func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB(t)
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	r := setupRouter()
	r.POST("/auth/register", controllers.Register)
	r.POST("/auth/login", controllers.Login)

	regBody := `{"username":"wrongpwuser","password":"correctpassword"}`
	req1 := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(regBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	loginBody := `{"username":"wrongpwuser","password":"wrongpassword"}`
	req2 := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}
