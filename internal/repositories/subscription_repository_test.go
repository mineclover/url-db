package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-db/internal/models"
	"url-db/internal/repositories"
)

func TestSubscriptionRepository_Create(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64)
		subscription *models.NodeSubscription
		wantErr      bool
		errType      error
	}{
		{
			name: "성공적인_구독_생성",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain and node first
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, int64(node.ID)
			},
			subscription: &models.NodeSubscription{
				SubscriberService: "test-service",
				EventTypes:        models.EventTypeList{"created", "updated"},
				IsActive:          true,
			},
			wantErr: false,
		},
		{
			name: "엔드포인트가_있는_구독_생성",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain and node first
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, int64(node.ID)
			},
			subscription: &models.NodeSubscription{
				SubscriberService:  "webhook-service",
				SubscriberEndpoint: stringPtr("https://api.example.com/webhook"),
				EventTypes:         models.EventTypeList{"created", "updated", "deleted"},
				IsActive:           true,
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				return testDB, repo, 999 // Non-existent node
			},
			subscription: &models.NodeSubscription{
				SubscriberService: "test-service",
				EventTypes:        models.EventTypeList{"created"},
				IsActive:          true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeID := tt.setup(t)
			defer testDB.Close()

			tt.subscription.SubscribedNodeID = nodeID

			err := repo.Create(tt.subscription)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.subscription.ID)
			}
		})
	}
}

func TestSubscriptionRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64)
		wantErr bool
		errType error
	}{
		{
			name: "존재하는_구독_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain, node, and subscription
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				subscription := &models.NodeSubscription{
					SubscriberService:  "test-service",
					SubscribedNodeID:   int64(node.ID),
					EventTypes:         models.EventTypeList{"created", "updated"},
					IsActive:           true,
				}
				err := repo.Create(subscription)
				require.NoError(t, err)
				
				return testDB, repo, subscription.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_구독_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				return testDB, repo, 999
			},
			wantErr: false, // Returns nil, not error for non-existent
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			subscription, err := repo.GetByID(id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, subscription)
			} else {
				assert.NoError(t, err)
				if id == 999 {
					assert.Nil(t, subscription)
				} else {
					assert.NotNil(t, subscription)
					assert.Equal(t, id, subscription.ID)
					assert.Equal(t, "test-service", subscription.SubscriberService)
				}
			}
		})
	}
}

func TestSubscriptionRepository_GetByNode(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64)
		expectedLen int
	}{
		{
			name: "노드의_모든_활성_구독_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create multiple subscriptions
				for i := 1; i <= 3; i++ {
					subscription := &models.NodeSubscription{
						SubscriberService: "service-" + string(rune('0'+i)),
						SubscribedNodeID:  int64(node.ID),
						EventTypes:        models.EventTypeList{"created"},
						IsActive:          true,
					}
					err := repo.Create(subscription)
					require.NoError(t, err)
				}
				
				// Create one inactive subscription (should not be returned)
				inactiveSubscription := &models.NodeSubscription{
					SubscriberService: "inactive-service",
					SubscribedNodeID:  int64(node.ID),
					EventTypes:        models.EventTypeList{"created"},
					IsActive:          false,
				}
				err := repo.Create(inactiveSubscription)
				require.NoError(t, err)
				
				return testDB, repo, int64(node.ID)
			},
			expectedLen: 3,
		},
		{
			name: "구독이_없는_노드",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain and node only
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, int64(node.ID)
			},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeID := tt.setup(t)
			defer testDB.Close()

			subscriptions, err := repo.GetByNode(nodeID)

			assert.NoError(t, err)
			assert.Len(t, subscriptions, tt.expectedLen)
			
			// Check that all returned subscriptions are active
			for _, sub := range subscriptions {
				assert.True(t, sub.IsActive)
				assert.Equal(t, nodeID, sub.SubscribedNodeID)
			}
		})
	}
}

func TestSubscriptionRepository_GetByService(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository)
		service     string
		expectedLen int
	}{
		{
			name: "서비스의_모든_구독_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create multiple subscriptions for target service
				for i := 1; i <= 2; i++ {
					subscription := &models.NodeSubscription{
						SubscriberService: "target-service",
						SubscribedNodeID:  int64(node.ID),
						EventTypes:        models.EventTypeList{"created"},
						IsActive:          true,
					}
					err := repo.Create(subscription)
					require.NoError(t, err)
				}
				
				// Create subscription for different service
				otherSubscription := &models.NodeSubscription{
					SubscriberService: "other-service",
					SubscribedNodeID:  int64(node.ID),
					EventTypes:        models.EventTypeList{"created"},
					IsActive:          true,
				}
				err := repo.Create(otherSubscription)
				require.NoError(t, err)
				
				return testDB, repo
			},
			service:     "target-service",
			expectedLen: 2,
		},
		{
			name: "존재하지_않는_서비스",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				return testDB, repo
			},
			service:     "non-existent-service",
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			subscriptions, err := repo.GetByService(tt.service)

			assert.NoError(t, err)
			assert.Len(t, subscriptions, tt.expectedLen)
			
			// Check that all returned subscriptions belong to the service
			for _, sub := range subscriptions {
				assert.Equal(t, tt.service, sub.SubscriberService)
			}
		})
	}
}

func TestSubscriptionRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64)
		updates map[string]interface{}
		wantErr bool
	}{
		{
			name: "성공적인_구독_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain, node, and subscription
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				subscription := &models.NodeSubscription{
					SubscriberService: "test-service",
					SubscribedNodeID:  int64(node.ID),
					EventTypes:        models.EventTypeList{"created"},
					IsActive:          true,
				}
				err := repo.Create(subscription)
				require.NoError(t, err)
				
				return testDB, repo, subscription.ID
			},
			updates: map[string]interface{}{
				"is_active": false,
				"event_types": `["created","updated","deleted"]`,
			},
			wantErr: false,
		},
		{
			name: "빈_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				return testDB, repo, 1
			},
			updates: map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "존재하지_않는_구독_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				return testDB, repo, 999
			},
			updates: map[string]interface{}{
				"is_active": false,
			},
			wantErr: false, // No error, just no rows affected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			err := repo.Update(id, tt.updates)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSubscriptionRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64)
		wantErr bool
		errType error
	}{
		{
			name: "성공적인_구독_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain, node, and subscription
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				subscription := &models.NodeSubscription{
					SubscriberService: "test-service",
					SubscribedNodeID:  int64(node.ID),
					EventTypes:        models.EventTypeList{"created"},
					IsActive:          true,
				}
				err := repo.Create(subscription)
				require.NoError(t, err)
				
				return testDB, repo, subscription.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_구독_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				return testDB, repo, 999
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			err := repo.Delete(id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Verify deletion
				subscription, err := repo.GetByID(id)
				assert.NoError(t, err)
				assert.Nil(t, subscription)
			}
		})
	}
}

func TestSubscriptionRepository_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository)
		offset        int
		limit         int
		expectedLen   int
		expectedTotal int
	}{
		{
			name: "모든_구독_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create multiple subscriptions
				for i := 1; i <= 5; i++ {
					subscription := &models.NodeSubscription{
						SubscriberService: "service-" + string(rune('0'+i)),
						SubscribedNodeID:  int64(node.ID),
						EventTypes:        models.EventTypeList{"created"},
						IsActive:          true,
					}
					err := repo.Create(subscription)
					require.NoError(t, err)
				}
				
				return testDB, repo
			},
			offset:        0,
			limit:         10,
			expectedLen:   5,
			expectedTotal: 5,
		},
		{
			name: "페이지네이션_테스트",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create 7 subscriptions
				for i := 1; i <= 7; i++ {
					subscription := &models.NodeSubscription{
						SubscriberService: "service-" + string(rune('0'+i)),
						SubscribedNodeID:  int64(node.ID),
						EventTypes:        models.EventTypeList{"created"},
						IsActive:          true,
					}
					err := repo.Create(subscription)
					require.NoError(t, err)
				}
				
				return testDB, repo
			},
			offset:        2,
			limit:         3,
			expectedLen:   3,
			expectedTotal: 7,
		},
		{
			name: "구독이_없는_경우",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.SubscriptionRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSubscriptionRepository(testDB.SqlxDB())
				return testDB, repo
			},
			offset:        0,
			limit:         10,
			expectedLen:   0,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			subscriptions, total, err := repo.GetAll(tt.offset, tt.limit)

			assert.NoError(t, err)
			assert.Len(t, subscriptions, tt.expectedLen)
			assert.Equal(t, tt.expectedTotal, total)
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}