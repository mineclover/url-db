package repositories_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-db/internal/models"
	"url-db/internal/repositories"
)

func TestNodeAttributeRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int)
		nodeAttribute *models.NodeAttribute
		wantErr       bool
		errType       error
	}{
		{
			name: "성공적인_노드속성_생성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain, node, and attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
				
				return testDB, repo, node.ID, attr.ID
			},
			nodeAttribute: &models.NodeAttribute{
				Value:      "test-value",
				OrderIndex: nil,
			},
			wantErr: false,
		},
		{
			name: "순서_인덱스가_있는_노드속성_생성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain, node, and attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
				
				return testDB, repo, node.ID, attr.ID
			},
			nodeAttribute: &models.NodeAttribute{
				Value:      "test-value-with-order",
				OrderIndex: intPtr(1),
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain and attribute only
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
				
				return testDB, repo, 999, attr.ID // Non-existent node
			},
			nodeAttribute: &models.NodeAttribute{
				Value: "test-value",
			},
			wantErr: true,
			errType: repositories.ErrForeignKeyConstraint,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeID, attributeID := tt.setup(t)
			defer testDB.Close()

			tt.nodeAttribute.NodeID = nodeID
			tt.nodeAttribute.AttributeID = attributeID

			err := repo.Create(tt.nodeAttribute)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.nodeAttribute.ID)
				assert.False(t, tt.nodeAttribute.CreatedAt.IsZero())
			}
		})
	}
}

func TestNodeAttributeRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int)
		wantErr bool
		errType error
	}{
		{
			name: "존재하는_노드속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain, node, attribute, and node_attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
				
				nodeAttr := &models.NodeAttribute{
					NodeID:      node.ID,
					AttributeID: attr.ID,
					Value:       "test-value",
					OrderIndex:  nil,
				}
				err := repo.Create(nodeAttr)
				require.NoError(t, err)
				
				return testDB, repo, nodeAttr.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				return testDB, repo, 999
			},
			wantErr: true,
			errType: repositories.ErrNodeAttributeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			nodeAttr, err := repo.GetByID(id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, nodeAttr)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, nodeAttr)
				assert.Equal(t, id, nodeAttr.ID)
				assert.Equal(t, "test-value", nodeAttr.Value)
			}
		})
	}
}

func TestNodeAttributeRepository_GetByNodeAndAttribute(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int)
		wantErr     bool
		errType     error
		expectValue string
	}{
		{
			name: "존재하는_노드속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain, node, attribute, and node_attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
				
				nodeAttr := &models.NodeAttribute{
					NodeID:      node.ID,
					AttributeID: attr.ID,
					Value:       "specific-value",
					OrderIndex:  nil,
				}
				err := repo.Create(nodeAttr)
				require.NoError(t, err)
				
				return testDB, repo, node.ID, attr.ID
			},
			wantErr:     false,
			expectValue: "specific-value",
		},
		{
			name: "존재하지_않는_노드속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				return testDB, repo, 999, 999
			},
			wantErr: true,
			errType: repositories.ErrNodeAttributeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeID, attributeID := tt.setup(t)
			defer testDB.Close()

			nodeAttr, err := repo.GetByNodeAndAttribute(nodeID, attributeID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, nodeAttr)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, nodeAttr)
				assert.Equal(t, nodeID, nodeAttr.NodeID)
				assert.Equal(t, attributeID, nodeAttr.AttributeID)
				assert.Equal(t, tt.expectValue, nodeAttr.Value)
			}
		})
	}
}

func TestNodeAttributeRepository_ListByNode(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int)
		expectedLen int
	}{
		{
			name: "노드의_모든_속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain, node, and multiple attributes
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				// Create multiple attributes
				for i := 1; i <= 3; i++ {
					attr := &models.Attribute{
						DomainID:    domain.ID,
						Name:        fmt.Sprintf("attr%d", i),
						Type:        models.AttributeTypeString,
						Description: fmt.Sprintf("Attribute %d", i),
					}
					err := attrRepo.Create(attr)
					require.NoError(t, err)
					
					// Create node attribute for each
					nodeAttr := &models.NodeAttribute{
						NodeID:      node.ID,
						AttributeID: attr.ID,
						Value:       fmt.Sprintf("value%d", i),
						OrderIndex:  nil,
					}
					err = repo.Create(nodeAttr)
					require.NoError(t, err)
				}
				
				return testDB, repo, node.ID
			},
			expectedLen: 3,
		},
		{
			name: "속성이_없는_노드",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain and node only
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, node.ID
			},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeID := tt.setup(t)
			defer testDB.Close()

			nodeAttrs, err := repo.ListByNode(nodeID)

			assert.NoError(t, err)
			assert.Len(t, nodeAttrs, tt.expectedLen)
			
			// Check that all returned attributes belong to the specified node
			for _, attr := range nodeAttrs {
				assert.Equal(t, nodeID, attr.NodeID)
			}
		})
	}
}

func TestNodeAttributeRepository_Update(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, *models.NodeAttribute)
		updateFn func(*models.NodeAttribute)
		wantErr  bool
		errType  error
	}{
		{
			name: "성공적인_노드속성_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, *models.NodeAttribute) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain, node, attribute, and node_attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
				
				nodeAttr := &models.NodeAttribute{
					NodeID:      node.ID,
					AttributeID: attr.ID,
					Value:       "original-value",
					OrderIndex:  nil,
				}
				err := repo.Create(nodeAttr)
				require.NoError(t, err)
				
				return testDB, repo, nodeAttr
			},
			updateFn: func(nodeAttr *models.NodeAttribute) {
				nodeAttr.Value = "updated-value"
				nodeAttr.OrderIndex = intPtr(5)
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드속성_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, *models.NodeAttribute) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				return testDB, repo, &models.NodeAttribute{ID: 999}
			},
			updateFn: func(nodeAttr *models.NodeAttribute) {
				nodeAttr.Value = "non-existent"
			},
			wantErr: true,
			errType: repositories.ErrNodeAttributeNotFound,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeAttr := tt.setup(t)
			defer testDB.Close()

			tt.updateFn(nodeAttr)
			err := repo.Update(nodeAttr)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				
				// Verify update
				updated, err := repo.GetByID(nodeAttr.ID)
				assert.NoError(t, err)
				assert.Equal(t, nodeAttr.Value, updated.Value)
				assert.Equal(t, nodeAttr.OrderIndex, updated.OrderIndex)
			}
		})
	}
}

func TestNodeAttributeRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int)
		wantErr bool
		errType error
	}{
		{
			name: "성공적인_노드속성_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain, node, attribute, and node_attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
				
				nodeAttr := &models.NodeAttribute{
					NodeID:      node.ID,
					AttributeID: attr.ID,
					Value:       "to-be-deleted",
					OrderIndex:  nil,
				}
				err := repo.Create(nodeAttr)
				require.NoError(t, err)
				
				return testDB, repo, nodeAttr.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드속성_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				return testDB, repo, 999
			},
			wantErr: true,
			errType: repositories.ErrNodeAttributeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			err := repo.Delete(id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				
				// Verify deletion
				_, err := repo.GetByID(id)
				assert.ErrorIs(t, err, repositories.ErrNodeAttributeNotFound)
			}
		})
	}
}

func TestNodeAttributeRepository_ExistsByNodeAndAttribute(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int)
		expected    bool
	}{
		{
			name: "존재하는_노드속성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				
				// Create domain, node, attribute, and node_attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
				
				nodeAttr := &models.NodeAttribute{
					NodeID:      node.ID,
					AttributeID: attr.ID,
					Value:       "exists",
					OrderIndex:  nil,
				}
				err := repo.Create(nodeAttr)
				require.NoError(t, err)
				
				return testDB, repo, node.ID, attr.ID
			},
			expected: true,
		},
		{
			name: "존재하지_않는_노드속성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeAttributeRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
				return testDB, repo, 999, 999
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeID, attributeID := tt.setup(t)
			defer testDB.Close()

			exists, err := repo.ExistsByNodeAndAttribute(nodeID, attributeID)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestNodeAttributeRepository_BatchOperations(t *testing.T) {
	t.Run("BatchCreate_성공", func(t *testing.T) {
		testDB := repositories.SetupTestDB(t)
		defer testDB.Close()
		
		repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
		
		// Create domain, node, and attributes
		domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
		domain := repositories.CreateTestDomain(t, domainRepo)
		
		nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
		node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
		
		attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
		attr1 := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
		
		attr2 := &models.Attribute{
			DomainID:    domain.ID,
			Name:        "attr2",
			Type:        models.AttributeTypeString,
			Description: "Attribute 2",
		}
		err := attrRepo.Create(attr2)
		require.NoError(t, err)
		
		// Batch create node attributes
		nodeAttrs := []models.NodeAttribute{
			{NodeID: node.ID, AttributeID: attr1.ID, Value: "batch-value1"},
			{NodeID: node.ID, AttributeID: attr2.ID, Value: "batch-value2"},
		}
		
		err = repo.BatchCreate(nodeAttrs)
		assert.NoError(t, err)
		
		// Verify creation
		result, err := repo.ListByNode(node.ID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})
	
	t.Run("BatchDeleteByNode_성공", func(t *testing.T) {
		testDB := repositories.SetupTestDB(t)
		defer testDB.Close()
		
		repo := repositories.NewSQLiteNodeAttributeRepository(testDB.DB)
		
		// Create domain, node, attribute, and node_attributes
		domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
		domain := repositories.CreateTestDomain(t, domainRepo)
		
		nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
		node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
		
		attrRepo := repositories.NewSQLiteAttributeRepository(testDB.DB)
		attr := repositories.CreateTestAttribute(t, attrRepo, domain.ID)
		
		// Create multiple node attributes
		for i := 1; i <= 3; i++ {
			nodeAttr := &models.NodeAttribute{
				NodeID:      node.ID,
				AttributeID: attr.ID,
				Value:       fmt.Sprintf("value%d", i),
			}
			err := repo.Create(nodeAttr)
			require.NoError(t, err)
		}
		
		// Verify they exist
		before, err := repo.ListByNode(node.ID)
		assert.NoError(t, err)
		assert.Len(t, before, 3)
		
		// Batch delete by node
		err = repo.BatchDeleteByNode(node.ID)
		assert.NoError(t, err)
		
		// Verify deletion
		after, err := repo.ListByNode(node.ID)
		assert.NoError(t, err)
		assert.Len(t, after, 0)
	})
}

// Helper function
func intPtr(i int) *int {
	return &i
}