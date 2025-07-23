package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"url-db/internal/models"
)

// MockNodeService is a mock implementation of NodeService
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

func TestNewNodeHandler(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.BaseHandler)
	assert.Equal(t, mockService, handler.nodeService)
}

func TestNodeHandler_CreateNode_Success(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	createReq := &models.CreateNodeRequest{
		URL:         "https://example.com",
		Title:       "Test Node",
		Description: "Test node description",
	}

	expectedNode := &models.Node{
		ID:          1,
		Content:     "https://example.com",
		DomainID:    1,
		Title:       "Test Node",
		Description: "Test node description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("CreateNode", 1, createReq).Return(expectedNode, nil)

	router.POST("/domains/:domain_id/urls", handler.CreateNode)

	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/domains/1/urls", bytes.NewBuffer(reqBody))
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
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.POST("/domains/:domain_id/urls", handler.CreateNode)

	createReq := &models.CreateNodeRequest{URL: "https://example.com"}
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/domains/invalid/urls", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_CreateNode_InvalidJSON(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.POST("/domains/:domain_id/urls", handler.CreateNode)

	req := httptest.NewRequest("POST", "/domains/1/urls", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_CreateNode_ServiceError(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	createReq := &models.CreateNodeRequest{
		URL:   "https://example.com",
		Title: "Test Node",
	}

	mockService.On("CreateNode", 1, createReq).Return(nil, NewConflictError("URL already exists"))

	router.POST("/domains/:domain_id/urls", handler.CreateNode)

	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/domains/1/urls", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNodesByDomain_Success(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	expectedResponse := &models.NodeListResponse{
		Nodes: []models.Node{
			{ID: 1, Content: "https://example.com", Title: "Node 1"},
			{ID: 2, Content: "https://test.com", Title: "Node 2"},
		},
		TotalCount: 2,
		Page:       1,
		Size:       20,
		TotalPages: 1,
	}

	mockService.On("GetNodesByDomainID", 1, 1, 20).Return(expectedResponse, nil)

	router.GET("/domains/:domain_id/urls", handler.GetNodesByDomain)

	req := httptest.NewRequest("GET", "/domains/1/urls", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.NodeListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.TotalCount, response.TotalCount)
	assert.Len(t, response.Nodes, 2)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNodesByDomain_WithSearch(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	expectedResponse := &models.NodeListResponse{
		Nodes:      []models.Node{},
		TotalCount: 0,
		Page:       1,
		Size:       20,
		TotalPages: 0,
	}

	mockService.On("SearchNodes", 1, "test", 1, 20).Return(expectedResponse, nil)

	router.GET("/domains/:domain_id/urls", handler.GetNodesByDomain)

	req := httptest.NewRequest("GET", "/domains/1/urls?search=test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNodesByDomain_WithPagination(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	expectedResponse := &models.NodeListResponse{
		Nodes:      []models.Node{},
		TotalCount: 0,
		Page:       2,
		Size:       10,
		TotalPages: 0,
	}

	mockService.On("GetNodesByDomainID", 1, 2, 10).Return(expectedResponse, nil)

	router.GET("/domains/:domain_id/urls", handler.GetNodesByDomain)

	req := httptest.NewRequest("GET", "/domains/1/urls?page=2&size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNodesByDomain_LargeSizeLimit(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	expectedResponse := &models.NodeListResponse{
		Nodes:      []models.Node{},
		TotalCount: 0,
		Page:       1,
		Size:       100, // Should be capped at 100
		TotalPages: 0,
	}

	// Should call with size=100 even though we requested 200
	mockService.On("GetNodesByDomainID", 1, 1, 100).Return(expectedResponse, nil)

	router.GET("/domains/:domain_id/urls", handler.GetNodesByDomain)

	req := httptest.NewRequest("GET", "/domains/1/urls?size=200", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_GetNodesByDomain_InvalidDomainID(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.GET("/domains/:domain_id/urls", handler.GetNodesByDomain)

	req := httptest.NewRequest("GET", "/domains/invalid/urls", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_GetNode_Success(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	expectedNode := &models.Node{
		ID:          1,
		Content:     "https://example.com",
		DomainID:    1,
		Title:       "Test Node",
		Description: "Test node description",
	}

	mockService.On("GetNodeByID", 1).Return(expectedNode, nil)

	router.GET("/urls/:id", handler.GetNode)

	req := httptest.NewRequest("GET", "/urls/1", nil)
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
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.GET("/urls/:id", handler.GetNode)

	req := httptest.NewRequest("GET", "/urls/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_GetNode_NotFound(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	mockService.On("GetNodeByID", 999).Return(nil, NewNotFoundError("Node not found"))

	router.GET("/urls/:id", handler.GetNode)

	req := httptest.NewRequest("GET", "/urls/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_UpdateNode_Success(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	updateReq := &models.UpdateNodeRequest{
		Title:       "Updated Node",
		Description: "Updated description",
	}

	expectedNode := &models.Node{
		ID:          1,
		Content:     "https://example.com",
		DomainID:    1,
		Title:       "Updated Node",
		Description: "Updated description",
	}

	mockService.On("UpdateNode", 1, updateReq).Return(expectedNode, nil)

	router.PUT("/urls/:id", handler.UpdateNode)

	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/urls/1", bytes.NewBuffer(reqBody))
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
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.PUT("/urls/:id", handler.UpdateNode)

	updateReq := &models.UpdateNodeRequest{Title: "test"}
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/urls/invalid", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_UpdateNode_InvalidJSON(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.PUT("/urls/:id", handler.UpdateNode)

	req := httptest.NewRequest("PUT", "/urls/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_DeleteNode_Success(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	mockService.On("DeleteNode", 1).Return(nil)

	router.DELETE("/urls/:id", handler.DeleteNode)

	req := httptest.NewRequest("DELETE", "/urls/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_DeleteNode_InvalidID(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.DELETE("/urls/:id", handler.DeleteNode)

	req := httptest.NewRequest("DELETE", "/urls/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_DeleteNode_NotFound(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	mockService.On("DeleteNode", 999).Return(NewNotFoundError("Node not found"))

	router.DELETE("/urls/:id", handler.DeleteNode)

	req := httptest.NewRequest("DELETE", "/urls/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_FindNodeByURL_Success(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	findReq := &models.FindNodeByURLRequest{
		URL: "https://example.com",
	}

	expectedNode := &models.Node{
		ID:       1,
		Content:  "https://example.com",
		DomainID: 1,
		Title:    "Found Node",
	}

	mockService.On("FindNodeByURL", 1, findReq).Return(expectedNode, nil)

	router.POST("/domains/:domain_id/urls/find", handler.FindNodeByURL)

	reqBody, _ := json.Marshal(findReq)
	req := httptest.NewRequest("POST", "/domains/1/urls/find", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
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

func TestNodeHandler_FindNodeByURL_InvalidDomainID(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.POST("/domains/:domain_id/urls/find", handler.FindNodeByURL)

	findReq := &models.FindNodeByURLRequest{URL: "https://example.com"}
	reqBody, _ := json.Marshal(findReq)
	req := httptest.NewRequest("POST", "/domains/invalid/urls/find", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_FindNodeByURL_InvalidJSON(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	router.POST("/domains/:domain_id/urls/find", handler.FindNodeByURL)

	req := httptest.NewRequest("POST", "/domains/1/urls/find", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNodeHandler_FindNodeByURL_NotFound(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	findReq := &models.FindNodeByURLRequest{
		URL: "https://notfound.com",
	}

	mockService.On("FindNodeByURL", 1, findReq).Return(nil, NewNotFoundError("Node not found"))

	router.POST("/domains/:domain_id/urls/find", handler.FindNodeByURL)

	reqBody, _ := json.Marshal(findReq)
	req := httptest.NewRequest("POST", "/domains/1/urls/find", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestNodeHandler_RegisterRoutes(t *testing.T) {
	mockService := &MockNodeService{}
	handler := NewNodeHandler(mockService)
	router := setupTestRouter()

	// This method registers routes - we just test it doesn't panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(router)
	})
}