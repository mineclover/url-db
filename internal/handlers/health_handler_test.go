package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func setupHealthHandlerTest() (*gin.Engine, *MockHealthService, *HealthHandler) {
	gin.SetMode(gin.TestMode)
	
	mockService := &MockHealthService{}
	handler := NewHealthHandler(mockService)
	
	router := gin.New()
	handler.RegisterRoutes(router)
	
	return router, mockService, handler
}

func TestHealthHandler_GetHealth_Healthy(t *testing.T) {
	router, mockService, _ := setupHealthHandlerTest()

	mockService.On("CheckDatabaseConnection").Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
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

func TestHealthHandler_GetHealth_DatabaseUnhealthy(t *testing.T) {
	router, mockService, _ := setupHealthHandlerTest()

	dbError := errors.New("database connection failed")
	mockService.On("CheckDatabaseConnection").Return(dbError)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	
	var response HealthInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "unhealthy", response.Status)
	assert.Equal(t, "unhealthy", response.Database.Status)
	assert.Equal(t, dbError.Error(), response.Checks["database"])

	mockService.AssertExpectations(t)
}

func TestHealthHandler_GetLiveness(t *testing.T) {
	router, _, _ := setupHealthHandlerTest()

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "alive", response["status"])
	assert.Contains(t, response, "timestamp")
}

func TestHealthHandler_GetReadiness_Ready(t *testing.T) {
	router, mockService, _ := setupHealthHandlerTest()

	mockService.On("CheckDatabaseConnection").Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ready", response["status"])
	assert.Contains(t, response, "timestamp")

	mockService.AssertExpectations(t)
}

func TestHealthHandler_GetReadiness_NotReady(t *testing.T) {
	router, mockService, _ := setupHealthHandlerTest()

	dbError := errors.New("database connection failed")
	mockService.On("CheckDatabaseConnection").Return(dbError)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not ready", response["status"])
	assert.Equal(t, "database connection failed", response["reason"])
	assert.Equal(t, dbError.Error(), response["error"])

	mockService.AssertExpectations(t)
}

func TestHealthHandler_GetVersion(t *testing.T) {
	router, _, _ := setupHealthHandlerTest()

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
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