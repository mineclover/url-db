package nodeattributes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"url-db/internal/models"
)

// Mock implementations for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(nodeAttribute *models.NodeAttribute) error {
	args := m.Called(nodeAttribute)
	return args.Error(0)
}

func (m *MockRepository) GetByID(id int) (*models.NodeAttributeWithInfo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NodeAttributeWithInfo), args.Error(1)
}

func (m *MockRepository) GetByNodeID(nodeID int) ([]models.NodeAttributeWithInfo, error) {
	args := m.Called(nodeID)
	return args.Get(0).([]models.NodeAttributeWithInfo), args.Error(1)
}

func (m *MockRepository) Update(id int, req *models.UpdateNodeAttributeRequest) (*models.NodeAttribute, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NodeAttribute), args.Error(1)
}

func (m *MockRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) DeleteByNodeIDAndAttributeID(nodeID, attributeID int) error {
	args := m.Called(nodeID, attributeID)
	return args.Error(0)
}

func (m *MockRepository) GetByNodeIDAndAttributeID(nodeID, attributeID int) (*models.NodeAttributeWithInfo, error) {
	args := m.Called(nodeID, attributeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NodeAttributeWithInfo), args.Error(1)
}

func (m *MockRepository) GetMaxOrderIndex(nodeID, attributeID int) (int, error) {
	args := m.Called(nodeID, attributeID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) ReorderAfterIndex(nodeID, attributeID, afterIndex int) error {
	args := m.Called(nodeID, attributeID, afterIndex)
	return args.Error(0)
}

func (m *MockRepository) ValidateNodeAndAttributeDomain(nodeID, attributeID int) error {
	args := m.Called(nodeID, attributeID)
	return args.Error(0)
}

func (m *MockRepository) GetAttributeType(attributeID int) (models.AttributeType, error) {
	args := m.Called(attributeID)
	return args.Get(0).(models.AttributeType), args.Error(1)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(attributeType models.AttributeType, value string, orderIndex *int) error {
	args := m.Called(attributeType, value, orderIndex)
	return args.Error(0)
}

func (m *MockValidator) ValidateValue(attributeType models.AttributeType, value string) error {
	args := m.Called(attributeType, value)
	return args.Error(0)
}

func (m *MockValidator) ValidateOrderIndex(attributeType models.AttributeType, orderIndex *int) error {
	args := m.Called(attributeType, orderIndex)
	return args.Error(0)
}

type MockOrderManager struct {
	mock.Mock
}

func (m *MockOrderManager) AssignOrderIndex(nodeID, attributeID int, orderIndex *int) (*int, error) {
	args := m.Called(nodeID, attributeID, orderIndex)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockOrderManager) ValidateOrderIndex(nodeID, attributeID int, orderIndex *int, excludeID *int) error {
	args := m.Called(nodeID, attributeID, orderIndex, excludeID)
	return args.Error(0)
}

func (m *MockOrderManager) ReorderAfterDeletion(nodeID, attributeID int, deletedOrderIndex int) error {
	args := m.Called(nodeID, attributeID, deletedOrderIndex)
	return args.Error(0)
}

func TestService_CreateNodeAttribute(t *testing.T) {
	t.Run("should create node attribute successfully", func(t *testing.T) {
		repo := &MockRepository{}
		validator := &MockValidator{}
		orderManager := &MockOrderManager{}
		service := NewService(repo, validator, orderManager)

		req := &models.CreateNodeAttributeRequest{
			AttributeID: 1,
			Value:       "programming",
			OrderIndex:  nil,
		}

		repo.On("ValidateNodeAndAttributeDomain", 1, 1).Return(nil)
		repo.On("GetByNodeIDAndAttributeID", 1, 1).Return(nil, nil)
		repo.On("GetAttributeType", 1).Return(models.AttributeTypeTag, nil)
		validator.On("Validate", models.AttributeTypeTag, "programming", (*int)(nil)).Return(nil)
		repo.On("Create", mock.AnythingOfType("*models.NodeAttribute")).Return(nil)

		result, err := service.CreateNodeAttribute(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "programming", result.Value)
		repo.AssertExpectations(t)
		validator.AssertExpectations(t)
	})

	t.Run("should handle domain mismatch error", func(t *testing.T) {
		repo := &MockRepository{}
		validator := &MockValidator{}
		orderManager := &MockOrderManager{}
		service := NewService(repo, validator, orderManager)

		req := &models.CreateNodeAttributeRequest{
			AttributeID: 1,
			Value:       "programming",
		}

		repo.On("ValidateNodeAndAttributeDomain", 1, 1).Return(ErrNodeAttributeDomainMismatch)

		result, err := service.CreateNodeAttribute(1, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNodeAttributeDomainMismatch, err)
		repo.AssertExpectations(t)
	})
}

func TestService_GetNodeAttributeByID(t *testing.T) {
	t.Run("should get node attribute successfully", func(t *testing.T) {
		repo := &MockRepository{}
		validator := &MockValidator{}
		orderManager := &MockOrderManager{}
		service := NewService(repo, validator, orderManager)

		expected := &models.NodeAttributeWithInfo{
			ID:          1,
			NodeID:      1,
			AttributeID: 1,
			Name:        "category",
			Type:        models.AttributeTypeTag,
			Value:       "programming",
			OrderIndex:  nil,
			CreatedAt:   time.Now(),
		}

		repo.On("GetByID", 1).Return(expected, nil)

		result, err := service.GetNodeAttributeByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("should return error when not found", func(t *testing.T) {
		repo := &MockRepository{}
		validator := &MockValidator{}
		orderManager := &MockOrderManager{}
		service := NewService(repo, validator, orderManager)

		repo.On("GetByID", 1).Return(nil, nil)

		result, err := service.GetNodeAttributeByID(1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNodeAttributeNotFound, err)
		repo.AssertExpectations(t)
	})
}

func TestService_ValidateNodeAttribute(t *testing.T) {
	repo := &MockRepository{}
	validator := &MockValidator{}
	orderManager := &MockOrderManager{}
	service := NewService(repo, validator, orderManager)

	t.Run("should validate successfully", func(t *testing.T) {
		req := &models.CreateNodeAttributeRequest{
			AttributeID: 1,
			Value:       "programming",
		}

		err := service.ValidateNodeAttribute(1, req)
		assert.NoError(t, err)
	})

	t.Run("should return error for invalid attribute_id", func(t *testing.T) {
		req := &models.CreateNodeAttributeRequest{
			AttributeID: 0,
			Value:       "programming",
		}

		err := service.ValidateNodeAttribute(1, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "attribute_id must be a positive integer")
	})

	t.Run("should return error for empty value", func(t *testing.T) {
		req := &models.CreateNodeAttributeRequest{
			AttributeID: 1,
			Value:       "",
		}

		err := service.ValidateNodeAttribute(1, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value cannot be empty")
	})

	t.Run("should return error for value too long", func(t *testing.T) {
		req := &models.CreateNodeAttributeRequest{
			AttributeID: 1,
			Value:       string(make([]byte, 2049)),
		}

		err := service.ValidateNodeAttribute(1, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value must be 2048 characters or less")
	})
}
