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

type MockAttributeService struct {
	mock.Mock
}

func (m *MockAttributeService) CreateAttribute(domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error) {
	args := m.Called(domainID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attribute), args.Error(1)
}

func (m *MockAttributeService) GetAttributesByDomainID(domainID int) ([]models.Attribute, error) {
	args := m.Called(domainID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Attribute), args.Error(1)
}

func (m *MockAttributeService) GetAttributeByID(id int) (*models.Attribute, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attribute), args.Error(1)
}

func (m *MockAttributeService) UpdateAttribute(id int, req *models.UpdateAttributeRequest) (*models.Attribute, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attribute), args.Error(1)
}

func (m *MockAttributeService) DeleteAttribute(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupAttributeHandlerTest() (*gin.Engine, *MockAttributeService, *AttributeHandler) {
	gin.SetMode(gin.TestMode)

	mockService := &MockAttributeService{}
	handler := NewAttributeHandler(mockService)

	router := gin.New()
	handler.RegisterRoutes(router)

	return router, mockService, handler
}

func TestAttributeHandler_CreateAttribute_Success(t *testing.T) {
	router, mockService, _ := setupAttributeHandlerTest()

	expectedAttribute := &models.Attribute{
		ID:          1,
		DomainID:    1,
		Name:        "category",
		Type:        "tag",
		Description: "Category tag",
	}

	mockService.On("CreateAttribute", 1, mock.AnythingOfType("*models.CreateAttributeRequest")).Return(expectedAttribute, nil)

	reqBody := models.CreateAttributeRequest{
		Name:        "category",
		Type:        "tag",
		Description: "Category tag",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/domains/1/attributes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Attribute
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAttribute.ID, response.ID)
	assert.Equal(t, expectedAttribute.Name, response.Name)
	assert.Equal(t, expectedAttribute.Type, response.Type)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_CreateAttribute_InvalidDomainID(t *testing.T) {
	router, _, _ := setupAttributeHandlerTest()

	reqBody := models.CreateAttributeRequest{
		Name: "category",
		Type: "tag",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/domains/invalid/attributes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAttributeHandler_GetAttributesByDomain_Success(t *testing.T) {
	router, mockService, _ := setupAttributeHandlerTest()

	expectedAttributes := []models.Attribute{
		{ID: 1, DomainID: 1, Name: "category", Type: "tag"},
		{ID: 2, DomainID: 1, Name: "priority", Type: "number"},
	}

	mockService.On("GetAttributesByDomainID", 1).Return(expectedAttributes, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains/1/attributes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]models.Attribute
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response["attributes"], 2)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_GetAttributesByDomain_InvalidDomainID(t *testing.T) {
	router, _, _ := setupAttributeHandlerTest()

	req := httptest.NewRequest(http.MethodGet, "/api/domains/invalid/attributes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAttributeHandler_GetAttribute_Success(t *testing.T) {
	router, mockService, _ := setupAttributeHandlerTest()

	expectedAttribute := &models.Attribute{
		ID:       1,
		DomainID: 1,
		Name:     "category",
		Type:     "tag",
	}

	mockService.On("GetAttributeByID", 1).Return(expectedAttribute, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/attributes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Attribute
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAttribute.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_GetAttribute_InvalidID(t *testing.T) {
	router, _, _ := setupAttributeHandlerTest()

	req := httptest.NewRequest(http.MethodGet, "/api/attributes/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAttributeHandler_UpdateAttribute_Success(t *testing.T) {
	router, mockService, _ := setupAttributeHandlerTest()

	expectedAttribute := &models.Attribute{
		ID:          1,
		DomainID:    1,
		Name:        "category",
		Type:        "tag",
		Description: "Updated description",
	}

	mockService.On("UpdateAttribute", 1, mock.AnythingOfType("*models.UpdateAttributeRequest")).Return(expectedAttribute, nil)

	reqBody := models.UpdateAttributeRequest{
		Description: "Updated description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/attributes/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Attribute
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAttribute.Description, response.Description)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_DeleteAttribute_Success(t *testing.T) {
	router, mockService, _ := setupAttributeHandlerTest()

	mockService.On("DeleteAttribute", 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/attributes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_ServiceError(t *testing.T) {
	router, mockService, _ := setupAttributeHandlerTest()

	internalErr := errors.New("internal service error")
	mockService.On("GetAttributeByID", 1).Return(nil, internalErr)

	req := httptest.NewRequest(http.MethodGet, "/api/attributes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
