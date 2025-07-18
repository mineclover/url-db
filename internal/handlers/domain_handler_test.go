package handlers

import (
	"bytes"
	"encoding/json"
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
	return args.Get(0).(*models.Domain), args.Error(1)
}

func (m *MockDomainService) GetDomains(page, size int) (*models.DomainListResponse, error) {
	args := m.Called(page, size)
	return args.Get(0).(*models.DomainListResponse), args.Error(1)
}

func (m *MockDomainService) GetDomainByID(id int) (*models.Domain, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Domain), args.Error(1)
}

func (m *MockDomainService) UpdateDomain(id int, req *models.UpdateDomainRequest) (*models.Domain, error) {
	args := m.Called(id, req)
	return args.Get(0).(*models.Domain), args.Error(1)
}

func (m *MockDomainService) DeleteDomain(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupDomainHandler() (*DomainHandler, *MockDomainService, *gin.Engine) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDomainService{}
	handler := NewDomainHandler(mockService)

	r := gin.New()
	handler.RegisterRoutes(r)

	return handler, mockService, r
}

func TestDomainHandler_CreateDomain(t *testing.T) {
	handler, mockService, r := setupDomainHandler()

	domain := &models.Domain{
		ID:          1,
		Name:        "test-domain",
		Description: "Test domain",
	}

	mockService.On("CreateDomain", mock.AnythingOfType("*models.CreateDomainRequest")).Return(domain, nil)

	reqBody := `{"name":"test-domain","description":"Test domain"}`
	req := httptest.NewRequest("POST", "/api/domains", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Domain
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, domain.Name, response.Name)
	assert.Equal(t, domain.Description, response.Description)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_CreateDomain_InvalidJSON(t *testing.T) {
	_, _, r := setupDomainHandler()

	reqBody := `{"invalid": json}`
	req := httptest.NewRequest("POST", "/api/domains", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
}

func TestDomainHandler_GetDomains(t *testing.T) {
	handler, mockService, r := setupDomainHandler()

	domains := &models.DomainListResponse{
		Domains: []models.Domain{
			{ID: 1, Name: "domain1", Description: "Description 1"},
			{ID: 2, Name: "domain2", Description: "Description 2"},
		},
		TotalCount: 2,
		Page:       1,
		Size:       20,
		TotalPages: 1,
	}

	mockService.On("GetDomains", 1, 20).Return(domains, nil)

	req := httptest.NewRequest("GET", "/api/domains", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.DomainListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.Domains))
	assert.Equal(t, 2, response.TotalCount)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomains_WithPagination(t *testing.T) {
	handler, mockService, r := setupDomainHandler()

	domains := &models.DomainListResponse{
		Domains:    []models.Domain{},
		TotalCount: 0,
		Page:       2,
		Size:       10,
		TotalPages: 0,
	}

	mockService.On("GetDomains", 2, 10).Return(domains, nil)

	req := httptest.NewRequest("GET", "/api/domains?page=2&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.DomainListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response.Page)
	assert.Equal(t, 10, response.Size)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomain(t *testing.T) {
	handler, mockService, r := setupDomainHandler()

	domain := &models.Domain{
		ID:          1,
		Name:        "test-domain",
		Description: "Test domain",
	}

	mockService.On("GetDomainByID", 1).Return(domain, nil)

	req := httptest.NewRequest("GET", "/api/domains/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Domain
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, domain.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_GetDomain_InvalidID(t *testing.T) {
	_, _, r := setupDomainHandler()

	req := httptest.NewRequest("GET", "/api/domains/invalid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
}

func TestDomainHandler_UpdateDomain(t *testing.T) {
	handler, mockService, r := setupDomainHandler()

	updatedDomain := &models.Domain{
		ID:          1,
		Name:        "updated-domain",
		Description: "Updated description",
	}

	mockService.On("UpdateDomain", 1, mock.AnythingOfType("*models.UpdateDomainRequest")).Return(updatedDomain, nil)

	reqBody := `{"name":"updated-domain","description":"Updated description"}`
	req := httptest.NewRequest("PUT", "/api/domains/1", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Domain
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updatedDomain.Name, response.Name)
	assert.Equal(t, updatedDomain.Description, response.Description)

	mockService.AssertExpectations(t)
}

func TestDomainHandler_DeleteDomain(t *testing.T) {
	handler, mockService, r := setupDomainHandler()

	mockService.On("DeleteDomain", 1).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/domains/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestDomainHandler_DeleteDomain_NotFound(t *testing.T) {
	handler, mockService, r := setupDomainHandler()

	mockService.On("DeleteDomain", 999).Return(NewNotFoundError("Domain not found"))

	req := httptest.NewRequest("DELETE", "/api/domains/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", response.Error)

	mockService.AssertExpectations(t)
}
