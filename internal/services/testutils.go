package services

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/url-db/internal/models"
)

type MockDomainRepository struct {
	domains     map[int]*models.Domain
	nextID      int
	nameToID    map[string]int
	shouldError bool
}

func NewMockDomainRepository() *MockDomainRepository {
	return &MockDomainRepository{
		domains:  make(map[int]*models.Domain),
		nextID:   1,
		nameToID: make(map[string]int),
	}
}

func (m *MockDomainRepository) SetShouldError(shouldError bool) {
	m.shouldError = shouldError
}

func (m *MockDomainRepository) Create(ctx context.Context, domain *models.Domain) error {
	if m.shouldError {
		return sql.ErrConnDone
	}
	
	domain.ID = m.nextID
	m.nextID++
	domain.CreatedAt = time.Now()
	domain.UpdatedAt = time.Now()
	
	m.domains[domain.ID] = domain
	m.nameToID[domain.Name] = domain.ID
	
	return nil
}

func (m *MockDomainRepository) GetByID(ctx context.Context, id int) (*models.Domain, error) {
	if m.shouldError {
		return nil, sql.ErrConnDone
	}
	
	domain, exists := m.domains[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	
	return domain, nil
}

func (m *MockDomainRepository) GetByName(ctx context.Context, name string) (*models.Domain, error) {
	if m.shouldError {
		return nil, sql.ErrConnDone
	}
	
	id, exists := m.nameToID[name]
	if !exists {
		return nil, sql.ErrNoRows
	}
	
	return m.domains[id], nil
}

func (m *MockDomainRepository) List(ctx context.Context, page, size int) ([]*models.Domain, int, error) {
	if m.shouldError {
		return nil, 0, sql.ErrConnDone
	}
	
	var domains []*models.Domain
	for _, domain := range m.domains {
		domains = append(domains, domain)
	}
	
	offset := (page - 1) * size
	totalCount := len(domains)
	
	if offset >= totalCount {
		return []*models.Domain{}, totalCount, nil
	}
	
	end := offset + size
	if end > totalCount {
		end = totalCount
	}
	
	return domains[offset:end], totalCount, nil
}

func (m *MockDomainRepository) Update(ctx context.Context, domain *models.Domain) error {
	if m.shouldError {
		return sql.ErrConnDone
	}
	
	_, exists := m.domains[domain.ID]
	if !exists {
		return sql.ErrNoRows
	}
	
	domain.UpdatedAt = time.Now()
	m.domains[domain.ID] = domain
	
	return nil
}

func (m *MockDomainRepository) Delete(ctx context.Context, id int) error {
	if m.shouldError {
		return sql.ErrConnDone
	}
	
	domain, exists := m.domains[id]
	if !exists {
		return sql.ErrNoRows
	}
	
	delete(m.domains, id)
	delete(m.nameToID, domain.Name)
	
	return nil
}

func (m *MockDomainRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.shouldError {
		return false, sql.ErrConnDone
	}
	
	_, exists := m.nameToID[name]
	return exists, nil
}

type MockNodeRepository struct {
	nodes       map[int]*models.Node
	nextID      int
	shouldError bool
}

func NewMockNodeRepository() *MockNodeRepository {
	return &MockNodeRepository{
		nodes:  make(map[int]*models.Node),
		nextID: 1,
	}
}

func (m *MockNodeRepository) SetShouldError(shouldError bool) {
	m.shouldError = shouldError
}

func (m *MockNodeRepository) Create(ctx context.Context, node *models.Node) error {
	if m.shouldError {
		return sql.ErrConnDone
	}
	
	node.ID = m.nextID
	m.nextID++
	node.CreatedAt = time.Now()
	node.UpdatedAt = time.Now()
	
	m.nodes[node.ID] = node
	
	return nil
}

func (m *MockNodeRepository) GetByID(ctx context.Context, id int) (*models.Node, error) {
	if m.shouldError {
		return nil, sql.ErrConnDone
	}
	
	node, exists := m.nodes[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	
	return node, nil
}

func (m *MockNodeRepository) GetByDomainAndContent(ctx context.Context, domainID int, content string) (*models.Node, error) {
	if m.shouldError {
		return nil, sql.ErrConnDone
	}
	
	for _, node := range m.nodes {
		if node.DomainID == domainID && node.Content == content {
			return node, nil
		}
	}
	
	return nil, sql.ErrNoRows
}

func (m *MockNodeRepository) ListByDomain(ctx context.Context, domainID int, page, size int, search string) ([]*models.Node, int, error) {
	if m.shouldError {
		return nil, 0, sql.ErrConnDone
	}
	
	var nodes []*models.Node
	for _, node := range m.nodes {
		if node.DomainID == domainID {
			nodes = append(nodes, node)
		}
	}
	
	offset := (page - 1) * size
	totalCount := len(nodes)
	
	if offset >= totalCount {
		return []*models.Node{}, totalCount, nil
	}
	
	end := offset + size
	if end > totalCount {
		end = totalCount
	}
	
	return nodes[offset:end], totalCount, nil
}

func (m *MockNodeRepository) Update(ctx context.Context, node *models.Node) error {
	if m.shouldError {
		return sql.ErrConnDone
	}
	
	_, exists := m.nodes[node.ID]
	if !exists {
		return sql.ErrNoRows
	}
	
	node.UpdatedAt = time.Now()
	m.nodes[node.ID] = node
	
	return nil
}

func (m *MockNodeRepository) Delete(ctx context.Context, id int) error {
	if m.shouldError {
		return sql.ErrConnDone
	}
	
	_, exists := m.nodes[id]
	if !exists {
		return sql.ErrNoRows
	}
	
	delete(m.nodes, id)
	
	return nil
}

func (m *MockNodeRepository) ExistsByDomainAndContent(ctx context.Context, domainID int, content string) (bool, error) {
	if m.shouldError {
		return false, sql.ErrConnDone
	}
	
	for _, node := range m.nodes {
		if node.DomainID == domainID && node.Content == content {
			return true, nil
		}
	}
	
	return false, nil
}

func CreateTestDomainService(t *testing.T) (DomainService, *MockDomainRepository) {
	mockRepo := NewMockDomainRepository()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewDomainService(mockRepo, logger)
	
	return service, mockRepo
}

func CreateTestNodeService(t *testing.T) (NodeService, *MockNodeRepository, *MockDomainRepository) {
	mockNodeRepo := NewMockNodeRepository()
	mockDomainRepo := NewMockDomainRepository()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewNodeService(mockNodeRepo, mockDomainRepo, logger)
	
	return service, mockNodeRepo, mockDomainRepo
}

func CreateTestCompositeKeyService(t *testing.T) CompositeKeyService {
	return NewCompositeKeyService("test-tool")
}

func CreateTestDomain(name, description string) *models.Domain {
	return &models.Domain{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func CreateTestNode(domainID int, url, title, description string) *models.Node {
	return &models.Node{
		DomainID:    domainID,
		Content:     url,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func CreateTestAttribute(domainID int, name, attributeType, description string) *models.Attribute {
	return &models.Attribute{
		DomainID:    domainID,
		Name:        name,
		Type:        attributeType,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

func CreateTestContext() context.Context {
	return context.Background()
}