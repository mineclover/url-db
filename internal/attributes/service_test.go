package attributes

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/url-db/internal/models"
)

// Mock repository
type MockAttributeRepository struct {
	mock.Mock
}

func (m *MockAttributeRepository) Create(ctx context.Context, attribute *models.Attribute) error {
	args := m.Called(ctx, attribute)
	return args.Error(0)
}

func (m *MockAttributeRepository) GetByID(ctx context.Context, id int) (*models.Attribute, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attribute), args.Error(1)
}

func (m *MockAttributeRepository) GetByDomainID(ctx context.Context, domainID int) ([]*models.Attribute, error) {
	args := m.Called(ctx, domainID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Attribute), args.Error(1)
}

func (m *MockAttributeRepository) GetByDomainIDAndName(ctx context.Context, domainID int, name string) (*models.Attribute, error) {
	args := m.Called(ctx, domainID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attribute), args.Error(1)
}

func (m *MockAttributeRepository) Update(ctx context.Context, attribute *models.Attribute) error {
	args := m.Called(ctx, attribute)
	return args.Error(0)
}

func (m *MockAttributeRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAttributeRepository) HasValues(ctx context.Context, attributeID int) (bool, error) {
	args := m.Called(ctx, attributeID)
	return args.Bool(0), args.Error(1)
}

// Mock domain service
type MockDomainService struct {
	mock.Mock
}

func (m *MockDomainService) GetDomain(ctx context.Context, id int) (*models.Domain, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Domain), args.Error(1)
}

func TestAttributeService_CreateAttribute(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	domain := &models.Domain{
		ID:   1,
		Name: "test-domain",
	}

	req := &models.CreateAttributeRequest{
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Test description",
	}

	// Mock domain exists
	mockDomainService.On("GetDomain", ctx, 1).Return(domain, nil)
	
	// Mock attribute doesn't exist
	mockRepo.On("GetByDomainIDAndName", ctx, 1, "test-attribute").Return(nil, ErrAttributeNotFound)
	
	// Mock create success
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Attribute")).Return(nil).Run(func(args mock.Arguments) {
		attr := args.Get(1).(*models.Attribute)
		attr.ID = 1
	})

	result, err := service.CreateAttribute(ctx, 1, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "test-attribute", result.Name)
	assert.Equal(t, models.AttributeTypeTag, result.Type)
	assert.Equal(t, "Test description", result.Description)

	mockRepo.AssertExpectations(t)
	mockDomainService.AssertExpectations(t)
}

func TestAttributeService_CreateAttribute_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	tests := []struct {
		name    string
		req     *models.CreateAttributeRequest
		wantErr error
	}{
		{
			name: "empty name",
			req: &models.CreateAttributeRequest{
				Name: "",
				Type: models.AttributeTypeTag,
			},
			wantErr: ErrAttributeNameRequired,
		},
		{
			name: "long name",
			req: &models.CreateAttributeRequest{
				Name: string(make([]byte, 256)),
				Type: models.AttributeTypeTag,
			},
			wantErr: ErrAttributeNameTooLong,
		},
		{
			name: "empty type",
			req: &models.CreateAttributeRequest{
				Name: "test",
				Type: "",
			},
			wantErr: ErrAttributeTypeRequired,
		},
		{
			name: "invalid type",
			req: &models.CreateAttributeRequest{
				Name: "test",
				Type: "invalid",
			},
			wantErr: ErrAttributeTypeInvalid,
		},
		{
			name: "long description",
			req: &models.CreateAttributeRequest{
				Name:        "test",
				Type:        models.AttributeTypeTag,
				Description: string(make([]byte, 1001)),
			},
			wantErr: ErrDescriptionTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateAttribute(ctx, 1, tt.req)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAttributeService_CreateAttribute_DomainNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	req := &models.CreateAttributeRequest{
		Name: "test-attribute",
		Type: models.AttributeTypeTag,
	}

	mockDomainService.On("GetDomain", ctx, 1).Return(nil, ErrDomainNotFound)

	_, err := service.CreateAttribute(ctx, 1, req)
	assert.ErrorIs(t, err, ErrDomainNotFound)

	mockDomainService.AssertExpectations(t)
}

func TestAttributeService_CreateAttribute_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	domain := &models.Domain{ID: 1, Name: "test-domain"}
	existingAttr := &models.Attribute{
		ID:       1,
		DomainID: 1,
		Name:     "test-attribute",
		Type:     models.AttributeTypeTag,
	}

	req := &models.CreateAttributeRequest{
		Name: "test-attribute",
		Type: models.AttributeTypeTag,
	}

	mockDomainService.On("GetDomain", ctx, 1).Return(domain, nil)
	mockRepo.On("GetByDomainIDAndName", ctx, 1, "test-attribute").Return(existingAttr, nil)

	_, err := service.CreateAttribute(ctx, 1, req)
	assert.ErrorIs(t, err, ErrAttributeAlreadyExists)

	mockRepo.AssertExpectations(t)
	mockDomainService.AssertExpectations(t)
}

func TestAttributeService_GetAttribute(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	expectedAttr := &models.Attribute{
		ID:       1,
		DomainID: 1,
		Name:     "test-attribute",
		Type:     models.AttributeTypeTag,
	}

	mockRepo.On("GetByID", ctx, 1).Return(expectedAttr, nil)

	result, err := service.GetAttribute(ctx, 1)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedAttr, result)

	mockRepo.AssertExpectations(t)
}

func TestAttributeService_ListAttributes(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	domain := &models.Domain{ID: 1, Name: "test-domain"}
	attributes := []*models.Attribute{
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
	}

	mockDomainService.On("GetDomain", ctx, 1).Return(domain, nil)
	mockRepo.On("GetByDomainID", ctx, 1).Return(attributes, nil)

	result, err := service.ListAttributes(ctx, 1)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Attributes, 2)
	assert.Equal(t, "attr1", result.Attributes[0].Name)
	assert.Equal(t, "attr2", result.Attributes[1].Name)

	mockRepo.AssertExpectations(t)
	mockDomainService.AssertExpectations(t)
}

func TestAttributeService_UpdateAttribute(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	existingAttr := &models.Attribute{
		ID:          1,
		DomainID:    1,
		Name:        "test-attribute",
		Type:        models.AttributeTypeTag,
		Description: "Original description",
		CreatedAt:   time.Now(),
	}

	req := &models.UpdateAttributeRequest{
		Description: "Updated description",
	}

	mockRepo.On("GetByID", ctx, 1).Return(existingAttr, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Attribute")).Return(nil)

	result, err := service.UpdateAttribute(ctx, 1, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated description", result.Description)

	mockRepo.AssertExpectations(t)
}

func TestAttributeService_DeleteAttribute(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	existingAttr := &models.Attribute{
		ID:       1,
		DomainID: 1,
		Name:     "test-attribute",
		Type:     models.AttributeTypeTag,
	}

	mockRepo.On("GetByID", ctx, 1).Return(existingAttr, nil)
	mockRepo.On("HasValues", ctx, 1).Return(false, nil)
	mockRepo.On("Delete", ctx, 1).Return(nil)

	err := service.DeleteAttribute(ctx, 1)
	
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAttributeService_DeleteAttribute_HasValues(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAttributeRepository)
	mockDomainService := new(MockDomainService)
	service := NewAttributeService(mockRepo, mockDomainService)

	existingAttr := &models.Attribute{
		ID:       1,
		DomainID: 1,
		Name:     "test-attribute",
		Type:     models.AttributeTypeTag,
	}

	mockRepo.On("GetByID", ctx, 1).Return(existingAttr, nil)
	mockRepo.On("HasValues", ctx, 1).Return(true, nil)

	err := service.DeleteAttribute(ctx, 1)
	
	assert.ErrorIs(t, err, ErrAttributeHasValues)

	mockRepo.AssertExpectations(t)
}