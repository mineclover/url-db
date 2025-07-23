package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"url-db/internal/application/dto/request"
	"url-db/internal/application/usecase/node"
)

// NodeHandler handles HTTP requests for node operations
type NodeHandler struct {
	createUseCase *node.CreateNodeUseCase
	listUseCase   *node.ListNodesUseCase
}

// NewNodeHandler creates a new node handler
func NewNodeHandler(createUC *node.CreateNodeUseCase, listUC *node.ListNodesUseCase) *NodeHandler {
	return &NodeHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
	}
}

// CreateNode handles POST /nodes
func (h *NodeHandler) CreateNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req request.CreateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.createUseCase.Execute(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListNodes handles GET /domains/{domainName}/nodes
func (h *NodeHandler) ListNodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract domain name from URL path
	// This is a simplified version - in a real implementation, you'd use a router
	domainName := r.URL.Query().Get("domain")
	if domainName == "" {
		http.Error(w, "Domain name is required", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if size <= 0 {
		size = 20
	}

	response, err := h.listUseCase.Execute(r.Context(), domainName, page, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
