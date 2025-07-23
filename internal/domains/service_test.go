package domains_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"url-db/internal/domains"
	"url-db/internal/models"
)

// Mock DomainRepository
type MockDomainRepository struct {
	mock.Mock
}

func (m *MockDomainRepository) Create(ctx context.Context, domain *models.Domain) error {
	args := m.Called(ctx, domain)
	return args.Error(0)
}

func (m *MockDomainRepository) GetByID(ctx context.Context, id int) (*models.Domain, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Domain), args.Error(1)
}

func (m *MockDomainRepository) GetByName(ctx context.Context, name string) (*models.Domain, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Domain), args.Error(1)
}

func (m *MockDomainRepository) List(ctx context.Context, page, size int) ([]*models.Domain, int, error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).([]*models.Domain), args.Int(1), args.Error(2)
}

func (m *MockDomainRepository) Update(ctx context.Context, domain *models.Domain) error {
	args := m.Called(ctx, domain)
	return args.Error(0)
}

func (m *MockDomainRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDomainRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func TestNewDomainService(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	
	assert.NotNil(t, service)
}

func TestDomainService_CreateDomain_Success(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	req := &models.CreateDomainRequest{
		Name:        "test-domain",
		Description: "Test description",
	}
	
	expectedDomain := &models.Domain{
		ID:          1,
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	mockRepo.On("ExistsByName", ctx, "test-domain").Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Domain")).Return(nil).Run(func(args mock.Arguments) {
		domain := args.Get(1).(*models.Domain)
		domain.ID = 1
		domain.CreatedAt = expectedDomain.CreatedAt
		domain.UpdatedAt = expectedDomain.UpdatedAt
	})
	
	result, err := service.CreateDomain(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-domain", result.Name)
	assert.Equal(t, "Test description", result.Description)
	mockRepo.AssertExpectations(t)
}

func TestDomainService_CreateDomain_InvalidName(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	testCases := []struct {
		name        string
		domainName  string
		expectedErr string
	}{
		{"empty name", "", "domain name is required"},
		{"too long name", string(make([]byte, 256)), "domain name cannot exceed 255 characters"},
		{"invalid characters", "test domain", "domain name can only contain alphanumeric characters and hyphens"},
		{"special characters", "test@domain", "domain name can only contain alphanumeric characters and hyphens"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &models.CreateDomainRequest{
				Name:        tc.domainName,
				Description: "Test description",
			}
			
			result, err := service.CreateDomain(ctx, req)
			
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestDomainService_CreateDomain_InvalidDescription(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	req := &models.CreateDomainRequest{
		Name:        "test-domain",
		Description: string(make([]byte, 1001)), // Too long description
	}
	
	result, err := service.CreateDomain(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "description cannot exceed 1000 characters")
}

func TestDomainService_CreateDomain_AlreadyExists(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	req := &models.CreateDomainRequest{
		Name:        "existing-domain",
		Description: "Test description",
	}
	
	mockRepo.On("ExistsByName", ctx, "existing-domain").Return(true, nil)
	
	result, err := service.CreateDomain(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestDomainService_CreateDomain_RepositoryError(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	req := &models.CreateDomainRequest{
		Name:        "test-domain",
		Description: "Test description",
	}
	
	mockRepo.On("ExistsByName", ctx, "test-domain").Return(false, assert.AnError)
	
	result, err := service.CreateDomain(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestDomainService_GetDomain_Success(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	expectedDomain := &models.Domain{
		ID:          1,
		Name:        "test-domain",
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	mockRepo.On("GetByID", ctx, 1).Return(expectedDomain, nil)
	
	result, err := service.GetDomain(ctx, 1)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedDomain, result)
	mockRepo.AssertExpectations(t)
}

func TestDomainService_GetDomain_NotFound(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("GetByID", ctx, 999).Return(nil, sql.ErrNoRows)
	
	result, err := service.GetDomain(ctx, 999)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestDomainService_ListDomains_Success(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	domains := []*models.Domain{
		{ID: 1, Name: "domain1", Description: "Description 1"},
		{ID: 2, Name: "domain2", Description: "Description 2"},
	}
	
	mockRepo.On("List", ctx, 1, 10).Return(domains, 2, nil)
	
	result, err := service.ListDomains(ctx, 1, 10)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Domains, 2)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Size)
	assert.Equal(t, 1, result.TotalPages)
	mockRepo.AssertExpectations(t)
}

func TestDomainService_ListDomains_InvalidParams(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	testCases := []struct {
		name         string
		page         int
		size         int
		expectedPage int
		expectedSize int
	}{
		{"invalid page", 0, 10, 1, 10},
		{"invalid size", 1, 0, 1, 20},
		{"size too large", 1, 101, 1, 100},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Service normalizes invalid parameters, so we need to set up mock expectations
			mockRepo.On("List", ctx, tc.expectedPage, tc.expectedSize).Return([]*models.Domain{}, 0, nil)
			
			result, err := service.ListDomains(ctx, tc.page, tc.size)
			
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.expectedPage, result.Page)
			assert.Equal(t, tc.expectedSize, result.Size)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDomainService_ListDomains_RepositoryError(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("List", ctx, 1, 10).Return([]*models.Domain{}, 0, assert.AnError)
	
	result, err := service.ListDomains(ctx, 1, 10)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestDomainService_UpdateDomain_Success(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	existingDomain := &models.Domain{
		ID:          1,
		Name:        "test-domain",
		Description: "Old description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	req := &models.UpdateDomainRequest{
		Description: "Updated description",
	}
	
	mockRepo.On("GetByID", ctx, 1).Return(existingDomain, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Domain")).Return(nil)
	
	result, err := service.UpdateDomain(ctx, 1, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated description", result.Description)
	mockRepo.AssertExpectations(t)
}

func TestDomainService_UpdateDomain_NotFound(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	req := &models.UpdateDomainRequest{
		Description: "Updated description",
	}
	
	mockRepo.On("GetByID", ctx, 999).Return(nil, sql.ErrNoRows)
	
	result, err := service.UpdateDomain(ctx, 999, req)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestDomainService_DeleteDomain_Success(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	existingDomain := &models.Domain{
		ID:   1,
		Name: "test-domain",
	}
	
	mockRepo.On("GetByID", ctx, 1).Return(existingDomain, nil)
	mockRepo.On("Delete", ctx, 1).Return(nil)
	
	err := service.DeleteDomain(ctx, 1)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDomainService_DeleteDomain_NotFound(t *testing.T) {
	mockRepo := &MockDomainRepository{}
	service := domains.NewDomainService(mockRepo)
	ctx := context.Background()
	
	mockRepo.On("GetByID", ctx, 999).Return(nil, sql.ErrNoRows)
	
	err := service.DeleteDomain(ctx, 999)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}