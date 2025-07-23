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

type MockNodeService struct {
	mock.Mock
}

func (m *MockNodeService) CreateNode(domainID int, req *models.CreateNodeRequest) (*models.Node, error) {
	args := m.Called(domainID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Node), args.Error(1)
}

func (m *MockNodeService) GetNodeByID(id int) (*models.Node, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Node), args.Error(1)
}

func (m *MockNodeService) GetNodesByDomainID(domainID, page, size int) (*models.NodeListResponse, error) {
	args := m.Called(domainID, page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NodeListResponse), args.Error(1)
}

func (m *MockNodeService) FindNodeByURL(domainID int, req *models.FindNodeByURLRequest) (*models.Node, error) {
	args := m.Called(domainID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Node), args.Error(1)
}

func (m *MockNodeService) UpdateNode(id int, req *models.UpdateNodeRequest) (*models.Node, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Node), args.Error(1)
}

func (m *MockNodeService) DeleteNode(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockNodeService) SearchNodes(domainID int, query string, page, size int) (*models.NodeListResponse, error) {
	args := m.Called(domainID, query, page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NodeListResponse), args.Error(1)
}

func setupNodeHandlerTest() (*gin.Engine, *MockNodeService, *NodeHandler) {
	gin.SetMode(gin.TestMode)

	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)

	router := gin.New()
	handler.RegisterRoutes(router)

	return router, mockService, handler
}

func TestNodeHandler_CreateNode_Success(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	expectedNode := &models.Node{
		ID:          1,
		DomainID:    1,
		Content:     "https://example.com",
		Title:       "Example",
		Description: "Example website",
	}

	mockService.On("CreateNode", 1, mock.AnythingOfType("*models.CreateNodeRequest")).Return(expectedNode, nil)

	reqBody := models.CreateNodeRequest{
		URL:         "https://example.com",
		Title:       "Example",
		Description: "Example website",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/domains/1/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Node
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedNode.ID, response.ID)
	assert.Equal(t, expectedNode.Content, response.Content)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_CreateNode_InvalidDomainID(t *testing.T) {
	router, _, _ := setupNodeHandlerTest()

	reqBody := models.CreateNodeRequest{
		URL: "https://example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/domains/invalid/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_GetNodesByDomain_Success(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	expectedResponse := &models.NodeListResponse{
		Nodes: []models.Node{
			{ID: 1, DomainID: 1, Content: "https://example.com", Title: "Example"},
		},
		TotalCount: 1,
		Page:       1,
		Size:       20,
		TotalPages: 1,
	}

	mockService.On("GetNodesByDomainID", 1, 1, 20).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains/1/urls", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.NodeListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.TotalCount, response.TotalCount)
	assert.Len(t, response.Nodes, 1)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNodesByDomain_WithSearch(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	expectedResponse := &models.NodeListResponse{
		Nodes:      []models.Node{},
		TotalCount: 0,
		Page:       1,
		Size:       20,
		TotalPages: 0,
	}

	mockService.On("SearchNodes", 1, "example", 1, 20).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains/1/urls?search=example", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNodesByDomain_WithPagination(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	expectedResponse := &models.NodeListResponse{
		Nodes:      []models.Node{},
		TotalCount: 0,
		Page:       2,
		Size:       10,
		TotalPages: 0,
	}

	mockService.On("GetNodesByDomainID", 1, 2, 10).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains/1/urls?page=2&size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNodesByDomain_SizeLimitEnforced(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	expectedResponse := &models.NodeListResponse{
		Nodes:      []models.Node{},
		TotalCount: 0,
		Page:       1,
		Size:       100,
		TotalPages: 0,
	}

	// Should limit size to 100 even if requested size is higher
	mockService.On("GetNodesByDomainID", 1, 1, 100).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/domains/1/urls?size=200", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNode_Success(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	expectedNode := &models.Node{
		ID:       1,
		DomainID: 1,
		Content:  "https://example.com",
		Title:    "Example",
	}

	mockService.On("GetNodeByID", 1).Return(expectedNode, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/urls/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Node
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedNode.ID, response.ID)
	assert.Equal(t, expectedNode.Content, response.Content)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNode_InvalidID(t *testing.T) {
	router, _, _ := setupNodeHandlerTest()

	req := httptest.NewRequest(http.MethodGet, "/api/urls/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_GetNode_NotFound(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	notFoundErr := NewNotFoundError("node with ID 999 not found")
	mockService.On("GetNodeByID", 999).Return(nil, notFoundErr)

	req := httptest.NewRequest(http.MethodGet, "/api/urls/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestNodeHandler_UpdateNode_Success(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	expectedNode := &models.Node{
		ID:          1,
		DomainID:    1,
		Content:     "https://example.com",
		Title:       "Updated Example",
		Description: "Updated description",
	}

	mockService.On("UpdateNode", 1, mock.AnythingOfType("*models.UpdateNodeRequest")).Return(expectedNode, nil)

	reqBody := models.UpdateNodeRequest{
		Title:       "Updated Example",
		Description: "Updated description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/urls/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Node
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedNode.Title, response.Title)
	assert.Equal(t, expectedNode.Description, response.Description)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_UpdateNode_InvalidID(t *testing.T) {
	router, _, _ := setupNodeHandlerTest()

	reqBody := models.UpdateNodeRequest{
		Title: "Updated title",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/urls/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_DeleteNode_Success(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	mockService.On("DeleteNode", 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/urls/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestNodeHandler_DeleteNode_InvalidID(t *testing.T) {
	router, _, _ := setupNodeHandlerTest()

	req := httptest.NewRequest(http.MethodDelete, "/api/urls/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_FindNodeByURL_Success(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	expectedNode := &models.Node{
		ID:       1,
		DomainID: 1,
		Content:  "https://example.com",
		Title:    "Example",
	}

	mockService.On("FindNodeByURL", 1, mock.AnythingOfType("*models.FindNodeByURLRequest")).Return(expectedNode, nil)

	reqBody := models.FindNodeByURLRequest{
		URL: "https://example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/domains/1/urls/find", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Node
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedNode.Content, response.Content)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_FindNodeByURL_InvalidDomainID(t *testing.T) {
	router, _, _ := setupNodeHandlerTest()

	reqBody := models.FindNodeByURLRequest{
		URL: "https://example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/domains/invalid/urls/find", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_FindNodeByURL_InvalidJSON(t *testing.T) {
	router, _, _ := setupNodeHandlerTest()

	req := httptest.NewRequest(http.MethodPost, "/api/domains/1/urls/find", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_ServiceError(t *testing.T) {
	router, mockService, _ := setupNodeHandlerTest()

	internalErr := errors.New("internal service error")
	mockService.On("GetNodeByID", 1).Return(nil, internalErr)

	req := httptest.NewRequest(http.MethodGet, "/api/urls/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
