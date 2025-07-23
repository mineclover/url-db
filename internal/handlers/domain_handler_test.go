package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"url-db/internal/models"
)

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

func setupDomainHandlerTest() (*gin.Engine, *MockDomainService, *DomainHandler) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)

	router := gin.New()
	handler.RegisterRoutes(router)

	return router, mockService, handler
}

func TestDomainHandler_CreateDomain_Success(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	expectedDomain := &models.Domain{
		ID:          1,
		Name:        "tech",
		Description: "Technology domain",
	}

	mockService.On("CreateDomain", mock.AnythingOfType("*models.CreateDomainRequest")).Return(expectedDomain, nil)

	reqBody := models.CreateDomainRequest{
		Name:        "tech",
		Description: "Technology domain",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/domains", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Domain
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedDomain.ID, response.ID)
	assert.Equal(t, expectedDomain.Name, response.Name)
	assert.Equal(t, expectedDomain.Description, response.Description)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_CreateDomain_ValidationError(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	validationErr := NewValidationError("name", "domain name is required")
	mockService.On("CreateDomain", mock.AnythingOfType("*models.CreateDomainRequest")).Return(nil, validationErr)

	reqBody := models.CreateDomainRequest{
		Name:        "valid-name", // Use a valid name so JSON binding passes
		Description: "Domain that will fail service validation",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/domains", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestDomainHandler_CreateDomain_InvalidJSON(t *testing.T) {
	router, _, _ := setupDomainHandlerTest()

	req := httptest.NewRequest(http.MethodPost, "/api/domains", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_GetDomains_Success(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	expectedResponse := &models.DomainListResponse{
		Domains: []models.Domain{
			{ID: 1, Name: "tech", Description: "Technology"},
			{ID: 2, Name: "sports", Description: "Sports"},
		},
		TotalCount: 2,
		Page:       1,
		Size:       20,
		TotalPages: 1,
	}

	mockService.On("GetDomains", 1, 20).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains", nil)
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

func TestDomainHandler_GetDomains_WithPagination(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	expectedResponse := &models.DomainListResponse{
		Domains:    []models.Domain{},
		TotalCount: 0,
		Page:       2,
		Size:       10,
		TotalPages: 0,
	}

	mockService.On("GetDomains", 2, 10).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains?page=2&size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomains_SizeLimitEnforced(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	expectedResponse := &models.DomainListResponse{
		Domains:    []models.Domain{},
		TotalCount: 0,
		Page:       1,
		Size:       100,
		TotalPages: 0,
	}

	// Should limit size to 100 even if requested size is higher
	mockService.On("GetDomains", 1, 100).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains?size=200", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomain_Success(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	expectedDomain := &models.Domain{
		ID:          1,
		Name:        "tech",
		Description: "Technology domain",
	}

	mockService.On("GetDomainByID", 1).Return(expectedDomain, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains/1", nil)
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
	router, _, _ := setupDomainHandlerTest()

	req := httptest.NewRequest(http.MethodGet, "/api/domains/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_GetDomain_NotFound(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	notFoundErr := NewNotFoundError("domain with ID 999 not found")
	mockService.On("GetDomainByID", 999).Return(nil, notFoundErr)

	req := httptest.NewRequest(http.MethodGet, "/api/domains/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestDomainHandler_UpdateDomain_Success(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	expectedDomain := &models.Domain{
		ID:          1,
		Name:        "tech",
		Description: "Updated technology domain",
	}

	mockService.On("UpdateDomain", 1, mock.AnythingOfType("*models.UpdateDomainRequest")).Return(expectedDomain, nil)

	reqBody := models.UpdateDomainRequest{
		Description: "Updated technology domain",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/domains/1", bytes.NewBuffer(body))
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
	router, _, _ := setupDomainHandlerTest()

	reqBody := models.UpdateDomainRequest{
		Description: "Updated description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/domains/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_DeleteDomain_Success(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	mockService.On("DeleteDomain", 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/domains/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestDomainHandler_DeleteDomain_InvalidID(t *testing.T) {
	router, _, _ := setupDomainHandlerTest()

	req := httptest.NewRequest(http.MethodDelete, "/api/domains/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDomainHandler_DeleteDomain_NotFound(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	notFoundErr := NewNotFoundError("domain with ID 999 not found")
	mockService.On("DeleteDomain", 999).Return(notFoundErr)

	req := httptest.NewRequest(http.MethodDelete, "/api/domains/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// Additional error cases for better coverage
func TestDomainHandler_ServiceError(t *testing.T) {
	router, mockService, _ := setupDomainHandlerTest()

	internalErr := errors.New("internal service error")
	mockService.On("GetDomainByID", 1).Return(nil, internalErr)

	req := httptest.NewRequest(http.MethodGet, "/api/domains/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
