package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"url-db/internal/models"
)

// MockDomainService is a mock implementation of DomainService
type MockDomainService struct {
	mock.Mock
}

func (m *MockDomainService) CreateDomain(req *models.CreateDomainRequest) (*models.Domain, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Domain), args.Error(1)
}

func (m *MockDomainService) GetDomains(page, size int) (*models.DomainListResponse, error) {
	args := m.Called(page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DomainListResponse), args.Error(1)
}

func (m *MockDomainService) GetDomainByID(id int) (*models.Domain, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Domain), args.Error(1)
}

func (m *MockDomainService) UpdateDomain(id int, req *models.UpdateDomainRequest) (*models.Domain, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Domain), args.Error(1)
}

func (m *MockDomainService) DeleteDomain(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewDomainHandler(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.BaseHandler)
	assert.Equal(t, mockService, handler.domainService)
}

func TestDomainHandler_CreateDomain_Success(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	createReq := &models.CreateDomainRequest{
		Name:        "test-domain",
		Description: "Test domain description",
	}

	expectedDomain := &models.Domain{
		ID:          1,
		Name:        "test-domain",
		Description: "Test domain description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("CreateDomain", createReq).Return(expectedDomain, nil)

	router.POST("/domains", handler.CreateDomain)

	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/domains", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.Domain
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedDomain.ID, response.ID)
	assert.Equal(t, expectedDomain.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_CreateDomain_InvalidJSON(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	router.POST("/domains", handler.CreateDomain)

	req := httptest.NewRequest("POST", "/domains", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_CreateDomain_ServiceError(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	createReq := &models.CreateDomainRequest{
		Name:        "test-domain",
		Description: "Test domain description",
	}

	mockService.On("CreateDomain", createReq).Return(nil, NewConflictError("Domain already exists"))

	router.POST("/domains", handler.CreateDomain)

	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/domains", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomains_Success(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	expectedResponse := &models.DomainListResponse{
		Domains: []models.Domain{
			{ID: 1, Name: "domain1", Description: "Description 1"},
			{ID: 2, Name: "domain2", Description: "Description 2"},
		},
		TotalCount: 2,
		Page:       1,
		Size:       20,
		TotalPages: 1,
	}

	mockService.On("GetDomains", 1, 20).Return(expectedResponse, nil)

	router.GET("/domains", handler.GetDomains)

	req := httptest.NewRequest("GET", "/domains", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.DomainListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.TotalCount, response.TotalCount)
	assert.Len(t, response.Domains, 2)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomains_WithQuery(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	expectedResponse := &models.DomainListResponse{
		Domains:    []models.Domain{},
		TotalCount: 0,
		Page:       2,
		Size:       10,
		TotalPages: 0,
	}

	mockService.On("GetDomains", 2, 10).Return(expectedResponse, nil)

	router.GET("/domains", handler.GetDomains)

	req := httptest.NewRequest("GET", "/domains?page=2&size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomains_LargeSizeLimit(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	expectedResponse := &models.DomainListResponse{
		Domains:    []models.Domain{},
		TotalCount: 0,
		Page:       1,
		Size:       100, // Should be capped at 100
		TotalPages: 0,
	}

	// Should call with size=100 even though we requested 200
	mockService.On("GetDomains", 1, 100).Return(expectedResponse, nil)

	router.GET("/domains", handler.GetDomains)

	req := httptest.NewRequest("GET", "/domains?size=200", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomain_Success(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	expectedDomain := &models.Domain{
		ID:          1,
		Name:        "test-domain",
		Description: "Test domain description",
	}

	mockService.On("GetDomainByID", 1).Return(expectedDomain, nil)

	router.GET("/domains/:id", handler.GetDomain)

	req := httptest.NewRequest("GET", "/domains/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.Domain
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedDomain.ID, response.ID)
	assert.Equal(t, expectedDomain.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomain_InvalidID(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	router.GET("/domains/:id", handler.GetDomain)

	req := httptest.NewRequest("GET", "/domains/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_GetDomain_NotFound(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	mockService.On("GetDomainByID", 999).Return(nil, NewNotFoundError("Domain not found"))

	router.GET("/domains/:id", handler.GetDomain)

	req := httptest.NewRequest("GET", "/domains/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_UpdateDomain_Success(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	updateReq := &models.UpdateDomainRequest{
		Description: "Updated description",
	}

	expectedDomain := &models.Domain{
		ID:          1,
		Name:        "test-domain", // Name doesn't change in update
		Description: "Updated description",
	}

	mockService.On("UpdateDomain", 1, updateReq).Return(expectedDomain, nil)

	router.PUT("/domains/:id", handler.UpdateDomain)

	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/domains/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.Domain
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedDomain.Description, response.Description)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_UpdateDomain_InvalidID(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	router.PUT("/domains/:id", handler.UpdateDomain)

	updateReq := &models.UpdateDomainRequest{Description: "test"}
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/domains/invalid", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_DeleteDomain_Success(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	mockService.On("DeleteDomain", 1).Return(nil)

	router.DELETE("/domains/:id", handler.DeleteDomain)

	req := httptest.NewRequest("DELETE", "/domains/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_DeleteDomain_InvalidID(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	router.DELETE("/domains/:id", handler.DeleteDomain)

	req := httptest.NewRequest("DELETE", "/domains/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_DeleteDomain_NotFound(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	mockService.On("DeleteDomain", 999).Return(NewNotFoundError("Domain not found"))

	router.DELETE("/domains/:id", handler.DeleteDomain)

	req := httptest.NewRequest("DELETE", "/domains/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_RegisterRoutes(t *testing.T) {
	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)
	router := setupTestRouter()

	// This method registers routes - we just test it doesn't panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(router)
	})
}