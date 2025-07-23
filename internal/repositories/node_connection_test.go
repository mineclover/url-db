package repositories_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-db/internal/models"
	"url-db/internal/repositories"
)

func TestNodeConnectionRepository_Create(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int)
		connection *models.NodeConnection
		wantErr    bool
		errType    error
	}{
		{
			name: "성공적인_노드연결_생성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain and two nodes
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node1 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example1-%s.com", t.Name()),
					Title:    "Node 1",
					Description: "First node",
				}
				err := nodeRepo.Create(node1)
				require.NoError(t, err)
				
				node2 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example2-%s.com", t.Name()),
					Title:    "Node 2",
					Description: "Second node",
				}
				err = nodeRepo.Create(node2)
				require.NoError(t, err)
				
				return testDB, repo, node1.ID, node2.ID
			},
			connection: &models.NodeConnection{
				RelationshipType: "parent",
				Description:      "Parent-child relationship",
			},
			wantErr: false,
		},
		{
			name: "related_관계타입_노드연결_생성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain and two nodes
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node1 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example1-%s.com", t.Name()),
					Title:    "Node 1",
					Description: "First node",
				}
				err := nodeRepo.Create(node1)
				require.NoError(t, err)
				
				node2 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example2-%s.com", t.Name()),
					Title:    "Node 2",
					Description: "Second node",
				}
				err = nodeRepo.Create(node2)
				require.NoError(t, err)
				
				return testDB, repo, node1.ID, node2.ID
			},
			connection: &models.NodeConnection{
				RelationshipType: "related",
				Description:      "Related nodes",
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_소스노드",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain and one node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, 999, node.ID // Non-existent source node
			},
			connection: &models.NodeConnection{
				RelationshipType: "parent",
				Description:      "Invalid connection",
			},
			wantErr: true,
			errType: repositories.ErrForeignKeyConstraint,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, sourceID, targetID := tt.setup(t)
			defer testDB.Close()

			tt.connection.SourceNodeID = sourceID
			tt.connection.TargetNodeID = targetID

			err := repo.Create(ctx, tt.connection)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.connection.ID)
				assert.False(t, tt.connection.CreatedAt.IsZero())
			}
		})
	}
}

func TestNodeConnectionRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int)
		wantErr bool
		errType error
	}{
		{
			name: "존재하는_노드연결_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain, nodes, and connection
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node1 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example1-%s.com", t.Name()),
					Title:    "Node 1",
					Description: "First node",
				}
				err := nodeRepo.Create(node1)
				require.NoError(t, err)
				
				node2 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example2-%s.com", t.Name()),
					Title:    "Node 2",
					Description: "Second node",
				}
				err = nodeRepo.Create(node2)
				require.NoError(t, err)
				
				connection := &models.NodeConnection{
					SourceNodeID:     node1.ID,
					TargetNodeID:     node2.ID,
					RelationshipType: "parent",
					Description:      "Test connection",
				}
				err = repo.Create(context.Background(), connection)
				require.NoError(t, err)
				
				return testDB, repo, connection.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드연결_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				return testDB, repo, 999
			},
			wantErr: true,
			errType: repositories.ErrNodeConnectionNotFound,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			connection, err := repo.GetByID(ctx, id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, connection)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, connection)
				assert.Equal(t, id, connection.ID)
				assert.Equal(t, "parent", connection.RelationshipType)
			}
		})
	}
}

func TestNodeConnectionRepository_GetBySourceAndTarget(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int)
		relationshipType string
		expectValue      string
		wantErr          bool
		errType          error
	}{
		{
			name: "존재하는_노드연결_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain, nodes, and connection
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node1 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example1-%s.com", t.Name()),
					Title:    "Node 1",
					Description: "First node",
				}
				err := nodeRepo.Create(node1)
				require.NoError(t, err)
				
				node2 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example2-%s.com", t.Name()),
					Title:    "Node 2",
					Description: "Second node",
				}
				err = nodeRepo.Create(node2)
				require.NoError(t, err)
				
				connection := &models.NodeConnection{
					SourceNodeID:     node1.ID,
					TargetNodeID:     node2.ID,
					RelationshipType: "child",
					Description:      "Child relationship",
				}
				err = repo.Create(context.Background(), connection)
				require.NoError(t, err)
				
				return testDB, repo, node1.ID, node2.ID
			},
			relationshipType: "child",
			expectValue:      "Child relationship",
			wantErr:          false,
		},
		{
			name: "존재하지_않는_노드연결_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				return testDB, repo, 999, 999
			},
			relationshipType: "parent",
			wantErr:          true,
			errType:          repositories.ErrNodeConnectionNotFound,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, sourceID, targetID := tt.setup(t)
			defer testDB.Close()

			connection, err := repo.GetBySourceAndTarget(ctx, sourceID, targetID, tt.relationshipType)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, connection)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, connection)
				assert.Equal(t, sourceID, connection.SourceNodeID)
				assert.Equal(t, targetID, connection.TargetNodeID)
				assert.Equal(t, tt.relationshipType, connection.RelationshipType)
				assert.Equal(t, tt.expectValue, connection.Description)
			}
		})
	}
}

func TestNodeConnectionRepository_ListBySourceNode(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int)
		offset      int
		limit       int 
		expectedLen int
		expectedTotal int
	}{
		{
			name: "소스노드의_모든_연결_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain and multiple nodes
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				sourceNode := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create multiple target nodes and connections
				for i := 1; i <= 3; i++ {
					targetNode := &models.Node{
						DomainID: domain.ID,
						Content:  fmt.Sprintf("http://target%d-%s.com", i, t.Name()),
						Title:    fmt.Sprintf("Target Node %d", i),
						Description: fmt.Sprintf("Target node %d", i),
					}
					err := nodeRepo.Create(targetNode)
					require.NoError(t, err)
					
					connection := &models.NodeConnection{
						SourceNodeID:     sourceNode.ID,
						TargetNodeID:     targetNode.ID,
						RelationshipType: "child",
						Description:      fmt.Sprintf("Connection %d", i),
					}
					err = repo.Create(context.Background(), connection)
					require.NoError(t, err)
				}
				
				return testDB, repo, sourceNode.ID
			},
			offset:        0,
			limit:         10,
			expectedLen:   3,
			expectedTotal: 3,
		},
		{
			name: "페이지네이션_테스트",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain and multiple nodes
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				sourceNode := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create 5 connections
				for i := 1; i <= 5; i++ {
					targetNode := &models.Node{
						DomainID: domain.ID,
						Content:  fmt.Sprintf("http://target%d-%s.com", i, t.Name()),
						Title:    fmt.Sprintf("Target Node %d", i),
						Description: fmt.Sprintf("Target node %d", i),
					}
					err := nodeRepo.Create(targetNode)
					require.NoError(t, err)
					
					connection := &models.NodeConnection{
						SourceNodeID:     sourceNode.ID,
						TargetNodeID:     targetNode.ID,
						RelationshipType: "child",
						Description:      fmt.Sprintf("Connection %d", i),
					}
					err = repo.Create(context.Background(), connection)
					require.NoError(t, err)
				}
				
				return testDB, repo, sourceNode.ID
			},
			offset:        2,
			limit:         2,
			expectedLen:   2,
			expectedTotal: 5,
		},
		{
			name: "연결이_없는_노드",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain and node only
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, node.ID
			},
			offset:        0,
			limit:         10,
			expectedLen:   0,
			expectedTotal: 0,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, sourceNodeID := tt.setup(t)
			defer testDB.Close()

			connections, total, err := repo.ListBySourceNode(ctx, sourceNodeID, tt.offset, tt.limit)

			assert.NoError(t, err)
			assert.Len(t, connections, tt.expectedLen)
			assert.Equal(t, tt.expectedTotal, total)
			
			// Check that all returned connections have the correct source node
			for _, conn := range connections {
				assert.Equal(t, sourceNodeID, conn.SourceNodeID)
			}
		})
	}
}

func TestNodeConnectionRepository_Update(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, *models.NodeConnection)
		updateFn func(*models.NodeConnection)
		wantErr  bool
		errType  error
	}{
		{
			name: "성공적인_노드연결_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, *models.NodeConnection) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain, nodes, and connection
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node1 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example1-%s.com", t.Name()),
					Title:    "Node 1",
					Description: "First node",
				}
				err := nodeRepo.Create(node1)
				require.NoError(t, err)
				
				node2 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example2-%s.com", t.Name()),
					Title:    "Node 2",
					Description: "Second node",
				}
				err = nodeRepo.Create(node2)
				require.NoError(t, err)
				
				connection := &models.NodeConnection{
					SourceNodeID:     node1.ID,
					TargetNodeID:     node2.ID,
					RelationshipType: "parent",
					Description:      "Original description",
				}
				err = repo.Create(context.Background(), connection)
				require.NoError(t, err)
				
				return testDB, repo, connection
			},
			updateFn: func(connection *models.NodeConnection) {
				connection.RelationshipType = "child"
				connection.Description = "Updated description"
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드연결_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, *models.NodeConnection) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				return testDB, repo, &models.NodeConnection{ID: 999}
			},
			updateFn: func(connection *models.NodeConnection) {
				connection.RelationshipType = "non-existent"
			},
			wantErr: true,
			errType: repositories.ErrNodeConnectionNotFound,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, connection := tt.setup(t)
			defer testDB.Close()

			tt.updateFn(connection)
			err := repo.Update(ctx, connection)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				
				// Verify update
				updated, err := repo.GetByID(ctx, connection.ID)
				assert.NoError(t, err)
				assert.Equal(t, connection.RelationshipType, updated.RelationshipType)
				assert.Equal(t, connection.Description, updated.Description)
			}
		})
	}
}

func TestNodeConnectionRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int)
		wantErr bool
		errType error
	}{
		{
			name: "성공적인_노드연결_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain, nodes, and connection
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node1 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example1-%s.com", t.Name()),
					Title:    "Node 1",
					Description: "First node",
				}
				err := nodeRepo.Create(node1)
				require.NoError(t, err)
				
				node2 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example2-%s.com", t.Name()),
					Title:    "Node 2",
					Description: "Second node",
				}
				err = nodeRepo.Create(node2)
				require.NoError(t, err)
				
				connection := &models.NodeConnection{
					SourceNodeID:     node1.ID,
					TargetNodeID:     node2.ID,
					RelationshipType: "parent",
					Description:      "To be deleted",
				}
				err = repo.Create(context.Background(), connection)
				require.NoError(t, err)
				
				return testDB, repo, connection.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드연결_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				return testDB, repo, 999
			},
			wantErr: true,
			errType: repositories.ErrNodeConnectionNotFound,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			err := repo.Delete(ctx, id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				
				// Verify deletion
				_, err := repo.GetByID(ctx, id)
				assert.ErrorIs(t, err, repositories.ErrNodeConnectionNotFound)
			}
		})
	}
}

func TestNodeConnectionRepository_ExistsBySourceAndTarget(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int)
		relationshipType string
		expected         bool
	}{
		{
			name: "존재하는_노드연결",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				
				// Create domain, nodes, and connection
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node1 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example1-%s.com", t.Name()),
					Title:    "Node 1",
					Description: "First node",
				}
				err := nodeRepo.Create(node1)
				require.NoError(t, err)
				
				node2 := &models.Node{
					DomainID: domain.ID,
					Content:  fmt.Sprintf("http://example2-%s.com", t.Name()),
					Title:    "Node 2",
					Description: "Second node",
				}
				err = nodeRepo.Create(node2)
				require.NoError(t, err)
				
				connection := &models.NodeConnection{
					SourceNodeID:     node1.ID,
					TargetNodeID:     node2.ID,
					RelationshipType: "related",
					Description:      "Exists",
				}
				err = repo.Create(context.Background(), connection)
				require.NoError(t, err)
				
				return testDB, repo, node1.ID, node2.ID
			},
			relationshipType: "related",
			expected:         true,
		},
		{
			name: "존재하지_않는_노드연결",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.NodeConnectionRepository, int, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
				return testDB, repo, 999, 999
			},
			relationshipType: "parent",
			expected:         false,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, sourceID, targetID := tt.setup(t)
			defer testDB.Close()

			exists, err := repo.ExistsBySourceAndTarget(ctx, sourceID, targetID, tt.relationshipType)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestNodeConnectionRepository_BatchOperations(t *testing.T) {
	ctx := context.Background()
	
	t.Run("BatchCreate_성공", func(t *testing.T) {
		testDB := repositories.SetupTestDB(t)
		defer testDB.Close()
		
		repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
		
		// Create domain and nodes
		domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
		domain := repositories.CreateTestDomain(t, domainRepo)
		
		nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
		node1 := &models.Node{
			DomainID: domain.ID,
			Content:  fmt.Sprintf("http://node1-%s.com", t.Name()),
			Title:    "Node 1",
			Description: "First node",
		}
		err := nodeRepo.Create(node1)
		require.NoError(t, err)
		
		node2 := &models.Node{
			DomainID: domain.ID,
			Content:  fmt.Sprintf("http://node2-%s.com", t.Name()),
			Title:    "Node 2",
			Description: "Second node",
		}
		err = nodeRepo.Create(node2)
		require.NoError(t, err)
		
		node3 := &models.Node{
			DomainID: domain.ID,
			Content:  fmt.Sprintf("http://node3-%s.com", t.Name()),
			Title:    "Node 3",
			Description: "Third node",
		}
		err = nodeRepo.Create(node3)
		require.NoError(t, err)
		
		// Batch create connections
		connections := []models.NodeConnection{
			{SourceNodeID: node1.ID, TargetNodeID: node2.ID, RelationshipType: "parent", Description: "Parent connection"},
			{SourceNodeID: node1.ID, TargetNodeID: node3.ID, RelationshipType: "parent", Description: "Another parent connection"},
		}
		
		err = repo.BatchCreate(ctx, connections)
		assert.NoError(t, err)
		
		// Verify creation
		result, total, err := repo.ListBySourceNode(ctx, node1.ID, 0, 10)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, total)
	})
	
	t.Run("BatchDelete_성공", func(t *testing.T) {
		testDB := repositories.SetupTestDB(t)
		defer testDB.Close()
		
		repo := repositories.NewSQLiteNodeConnectionRepository(testDB.DB)
		
		// Create domain, nodes, and connections
		domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
		domain := repositories.CreateTestDomain(t, domainRepo)
		
		nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
		node1 := &models.Node{
			DomainID: domain.ID,
			Content:  fmt.Sprintf("http://node1-%s.com", t.Name()),
			Title:    "Node 1",
			Description: "First node",
		}
		err := nodeRepo.Create(node1)
		require.NoError(t, err)
		
		node2 := &models.Node{
			DomainID: domain.ID,
			Content:  fmt.Sprintf("http://node2-%s.com", t.Name()),
			Title:    "Node 2",
			Description: "Second node",
		}
		err = nodeRepo.Create(node2)
		require.NoError(t, err)
		
		// Create multiple connections with different relationship types
		var connectionIDs []int
		relationshipTypes := []string{"child", "parent", "related"}
		for i := 0; i < 3; i++ {
			connection := &models.NodeConnection{
				SourceNodeID:     node1.ID,
				TargetNodeID:     node2.ID,
				RelationshipType: relationshipTypes[i],
				Description:      fmt.Sprintf("Connection %d", i),
			}
			err = repo.Create(ctx, connection)
			require.NoError(t, err)
			connectionIDs = append(connectionIDs, connection.ID)
		}
		
		// Verify they exist
		before, total, err := repo.ListBySourceNode(ctx, node1.ID, 0, 10)
		assert.NoError(t, err)
		assert.Len(t, before, 3)
		assert.Equal(t, 3, total)
		
		// Batch delete first two connections
		err = repo.BatchDelete(ctx, connectionIDs[:2])
		assert.NoError(t, err)
		
		// Verify deletion
		after, total, err := repo.ListBySourceNode(ctx, node1.ID, 0, 10)
		assert.NoError(t, err)
		assert.Len(t, after, 1)
		assert.Equal(t, 1, total)
	})
}