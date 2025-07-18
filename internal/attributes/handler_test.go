package attributes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"url-db/internal/models"
)

// Mock service
type MockAttributeService struct {
	mock.Mock
}

func (m *MockAttributeService) CreateAttribute(ctx context.Context, domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error) {
	args := m.Called(ctx, domainID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attribute), args.Error(1)
}

func (m *MockAttributeService) GetAttribute(ctx context.Context, id int) (*models.Attribute, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attribute), args.Error(1)
}

func (m *MockAttributeService) ListAttributes(ctx context.Context, domainID int) (*models.AttributeListResponse, error) {
	args := m.Called(ctx, domainID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AttributeListResponse), args.Error(1)
}

func (m *MockAttributeService) UpdateAttribute(ctx context.Context, id int, req *models.UpdateAttributeRequest) (*models.Attribute, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attribute), args.Error(1)
}

func (m *MockAttributeService) DeleteAttribute(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter(handler *AttributeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apiGroup := router.Group("/api")
	handler.RegisterRoutes(apiGroup)
	return router
}

func TestAttributeHandler_CreateAttribute(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	req := &models.CreateAttributeRequest{
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
	}

	expectedAttr := &models.Attribute{
		ID:          1,
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
	}

	mockService.On("CreateAttribute", mock.AnythingOfType("*context.cancelCtx"), 1, req).Return(expectedAttr, nil)

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/domains/1/attributes", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Attribute
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAttr.ID, response.ID)
	assert.Equal(t, expectedAttr.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_CreateAttribute_ValidationError(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	// Invalid request body
	reqBody := []byte(`{"invalid": "json"}`)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/domains/1/attributes", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response["error"])
}

func TestAttributeHandler_CreateAttribute_ServiceError(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	req := &models.CreateAttributeRequest{
		Name: "test-attribute",
		Type: models.AttributeTypeTag,
	}

	mockService.On("CreateAttribute", mock.AnythingOfType("*context.cancelCtx"), 1, req).Return(nil, ErrAttributeAlreadyExists)

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/domains/1/attributes", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "conflict", response["error"])
	assert.Equal(t, "Attribute already exists", response["message"])

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_GetAttribute(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	expectedAttr := &models.Attribute{
		ID:          1,
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
	}

	mockService.On("GetAttribute", mock.AnythingOfType("*context.cancelCtx"), 1).Return(expectedAttr, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/attributes/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Attribute
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAttr.ID, response.ID)
	assert.Equal(t, expectedAttr.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_GetAttribute_NotFound(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	mockService.On("GetAttribute", mock.AnythingOfType("*context.cancelCtx"), 1).Return(nil, ErrAttributeNotFound)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/attributes/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", response["error"])
	assert.Equal(t, "Attribute not found", response["message"])

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_ListAttributes(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	expectedResponse := &models.AttributeListResponse{
		Attributes: []models.Attribute{
			{
				ID:       1,
				DomainID: 1,
				Name:     "attr1",
				Type:     models.AttributeTypeTag,
			},
			{
				ID:       2,
				DomainID: 1,
				Name:     "attr2",
				Type:     models.AttributeTypeString,
			},
		},
	}

	mockService.On("ListAttributes", mock.AnythingOfType("*context.cancelCtx"), 1).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/domains/1/attributes", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AttributeListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Attributes, 2)
	assert.Equal(t, "attr1", response.Attributes[0].Name)
	assert.Equal(t, "attr2", response.Attributes[1].Name)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_UpdateAttribute(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	req := &models.UpdateAttributeRequest{
		Description: "Updated description",
	}

	expectedAttr := &models.Attribute{
		ID:          1,
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Updated description",
	}

	mockService.On("UpdateAttribute", mock.AnythingOfType("*context.cancelCtx"), 1, req).Return(expectedAttr, nil)

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/attributes/1", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Attribute
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAttr.Description, response.Description)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_DeleteAttribute(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	mockService.On("DeleteAttribute", mock.AnythingOfType("*context.cancelCtx"), 1).Return(nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/attributes/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockService.AssertExpectations(t)
}

func TestAttributeHandler_DeleteAttribute_HasValues(t *testing.T) {
	mockService := new(MockAttributeService)
	handler := NewAttributeHandler(mockService)
	router := setupRouter(handler)

	mockService.On("DeleteAttribute", mock.AnythingOfType("*context.cancelCtx"), 1).Return(ErrAttributeHasValues)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/attributes/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "conflict", response["error"])
	assert.Equal(t, "Cannot delete attribute with existing values", response["message"])

	mockService.AssertExpectations(t)
}
