package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHealthService is a mock implementation of HealthService
type MockHealthService struct {
	mock.Mock
}

func (m *MockHealthService) CheckDatabaseConnection() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockHealthService) GetSystemInfo() (*HealthInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HealthInfo), args.Error(1)
}

func TestNewHealthHandler(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.BaseHandler)
	assert.Equal(t, mockService, handler.healthService)
	assert.False(t, handler.startTime.IsZero())
}

func TestHealthHandler_GetHealth_Healthy(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	router := setupTestRouter()

	mockService.On("CheckDatabaseConnection").Return(nil)

	router.GET("/health", handler.GetHealth)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response HealthInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "1.0.0", response.Version)
	assert.Equal(t, "healthy", response.Database.Status)
	assert.Equal(t, "ok", response.Checks["database"])

	mockService.AssertExpectations(t)
}

func TestHealthHandler_GetHealth_Unhealthy(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	router := setupTestRouter()

	mockService.On("CheckDatabaseConnection").Return(errors.New("database connection failed"))

	router.GET("/health", handler.GetHealth)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	
	var response HealthInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "unhealthy", response.Status)
	assert.Equal(t, "unhealthy", response.Database.Status)
	assert.Equal(t, "database connection failed", response.Checks["database"])

	mockService.AssertExpectations(t)
}

func TestHealthHandler_GetLiveness(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	router := setupTestRouter()

	router.GET("/health/live", handler.GetLiveness)

	req := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "alive", response["status"])
	assert.NotNil(t, response["timestamp"])
}

func TestHealthHandler_GetReadiness_Ready(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	router := setupTestRouter()

	mockService.On("CheckDatabaseConnection").Return(nil)

	router.GET("/health/ready", handler.GetReadiness)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ready", response["status"])
	assert.NotNil(t, response["timestamp"])

	mockService.AssertExpectations(t)
}

func TestHealthHandler_GetReadiness_NotReady(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	router := setupTestRouter()

	dbError := errors.New("database connection failed")
	mockService.On("CheckDatabaseConnection").Return(dbError)

	router.GET("/health/ready", handler.GetReadiness)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not ready", response["status"])
	assert.Equal(t, "database connection failed", response["reason"])
	assert.Equal(t, "database connection failed", response["error"])
	assert.NotNil(t, response["timestamp"])

	mockService.AssertExpectations(t)
}

func TestHealthHandler_GetVersion(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	router := setupTestRouter()

	router.GET("/version", handler.GetVersion)

	req := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", response["version"])
	assert.Equal(t, "2024-01-01T00:00:00Z", response["build_time"])
	assert.Equal(t, "abc123", response["commit"])
	assert.Equal(t, "go1.21", response["go_version"])
}

func TestHealthHandler_RegisterRoutes(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	router := setupTestRouter()

	// This method registers routes - we just test it doesn't panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(router)
	})
}

func TestHealthInfo_SystemUptimeCalculation(t *testing.T) {
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	router := setupTestRouter()

	// Wait a small amount to ensure uptime is > 0
	time.Sleep(1 * time.Millisecond)

	mockService.On("CheckDatabaseConnection").Return(nil)

	router.GET("/health", handler.GetHealth)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response HealthInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Uptime should be greater than 0
	assert.True(t, response.System.Uptime > 0)
	assert.Equal(t, "go1.21", response.System.GoVersion)
	assert.Equal(t, "linux/amd64", response.System.Platform)

	mockService.AssertExpectations(t)
}