package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"url-db/internal/application/dto/request"
	"url-db/internal/application/usecase/domain"
)

// DomainHandler handles HTTP requests for domain operations
type DomainHandler struct {
	createUseCase *domain.CreateDomainUseCase
	listUseCase   *domain.ListDomainsUseCase
}

// NewDomainHandler creates a new domain handler
func NewDomainHandler(createUC *domain.CreateDomainUseCase, listUC *domain.ListDomainsUseCase) *DomainHandler {
	return &DomainHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
	}
}

// CreateDomain handles POST /domains
func (h *DomainHandler) CreateDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req request.CreateDomainRequest
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

// ListDomains handles GET /domains
func (h *DomainHandler) ListDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	response, err := h.listUseCase.Execute(r.Context(), page, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
