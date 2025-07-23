package repositories_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-db/internal/models"
	"url-db/internal/repositories"
)

func TestEventRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64)
		event   *models.NodeEvent
		wantErr bool
	}{
		{
			name: "성공적인_이벤트_생성",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node first
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, int64(node.ID)
			},
			event: &models.NodeEvent{
				EventType: "created",
				EventData: &models.EventData{
					EventID:   "evt-123",
					NodeID:    1,
					EventType: "created",
					Timestamp: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "업데이트_이벤트_생성",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node first
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, int64(node.ID)
			},
			event: &models.NodeEvent{
				EventType: "updated",
				EventData: &models.EventData{
					EventID:   "evt-124",
					NodeID:    1,
					EventType: "updated",
					Timestamp: time.Now(),
					Changes: &models.EventChanges{
						Before: map[string]interface{}{"title": "Old Title"},
						After:  map[string]interface{}{"title": "New Title"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_노드",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				return testDB, repo, 999 // Non-existent node
			},
			event: &models.NodeEvent{
				EventType: "created",
				EventData: &models.EventData{
					EventID:   "evt-125",
					NodeID:    999,
					EventType: "created",
					Timestamp: time.Now(),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeID := tt.setup(t)
			defer testDB.Close()

			tt.event.NodeID = nodeID

			err := repo.Create(tt.event)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.event.ID)
			}
		})
	}
}

func TestEventRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64)
		wantErr bool
	}{
		{
			name: "존재하는_이벤트_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain, node, and event
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				event := &models.NodeEvent{
					NodeID:    int64(node.ID),
					EventType: "created",
					EventData: &models.EventData{
						EventID:   "evt-123",
						NodeID:    int64(node.ID),
						EventType: "created",
						Timestamp: time.Now(),
					},
				}
				err := repo.Create(event)
				require.NoError(t, err)
				
				return testDB, repo, event.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_이벤트_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				return testDB, repo, 999
			},
			wantErr: false, // Returns nil, not error for non-existent
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			event, err := repo.GetByID(id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				if id == 999 {
					assert.Nil(t, event)
				} else {
					assert.NotNil(t, event)
					assert.Equal(t, id, event.ID)
					assert.Equal(t, "created", event.EventType)
				}
			}
		})
	}
}

func TestEventRepository_GetByNode(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64)
		limit       int
		expectedLen int
	}{
		{
			name: "노드의_모든_이벤트_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create multiple events
				eventTypes := []string{"created", "updated", "deleted"}
				for i, eventType := range eventTypes {
					event := &models.NodeEvent{
						NodeID:    int64(node.ID),
						EventType: eventType,
						EventData: &models.EventData{
							EventID:   "evt-" + string(rune('1'+i)),
							NodeID:    int64(node.ID),
							EventType: eventType,
							Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
						},
					}
					err := repo.Create(event)
					require.NoError(t, err)
				}
				
				return testDB, repo, int64(node.ID)
			},
			limit:       10,
			expectedLen: 3,
		},
		{
			name: "제한된_수의_이벤트_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create 5 events
				for i := 1; i <= 5; i++ {
					event := &models.NodeEvent{
						NodeID:    int64(node.ID),
						EventType: "updated",
						EventData: &models.EventData{
							EventID:   "evt-" + string(rune('0'+i)),
							NodeID:    int64(node.ID),
							EventType: "updated",
							Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
						},
					}
					err := repo.Create(event)
					require.NoError(t, err)
				}
				
				return testDB, repo, int64(node.ID)
			},
			limit:       3,
			expectedLen: 3,
		},
		{
			name: "이벤트가_없는_노드",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node only
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				return testDB, repo, int64(node.ID)
			},
			limit:       10,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, nodeID := tt.setup(t)
			defer testDB.Close()

			events, err := repo.GetByNode(nodeID, tt.limit)

			assert.NoError(t, err)
			assert.Len(t, events, tt.expectedLen)
			
			// Check that all returned events belong to the node
			for _, event := range events {
				assert.Equal(t, nodeID, event.NodeID)
			}
		})
	}
}

func TestEventRepository_GetPendingEvents(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository)
		limit       int
		expectedLen int
	}{
		{
			name: "대기_중인_이벤트_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create 3 pending events and 2 processed events
				for i := 1; i <= 5; i++ {
					event := &models.NodeEvent{
						NodeID:    int64(node.ID),
						EventType: "created",
						EventData: &models.EventData{
							EventID:   "evt-" + string(rune('0'+i)),
							NodeID:    int64(node.ID),
							EventType: "created",
							Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
						},
					}
					err := repo.Create(event)
					require.NoError(t, err)
					
					// Mark first 2 as processed
					if i <= 2 {
						err = repo.MarkAsProcessed(event.ID)
						require.NoError(t, err)
					}
				}
				
				return testDB, repo
			},
			limit:       10,
			expectedLen: 3,
		},
		{
			name: "대기_이벤트_없음",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				return testDB, repo
			},
			limit:       10,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			events, err := repo.GetPendingEvents(tt.limit)

			assert.NoError(t, err)
			assert.Len(t, events, tt.expectedLen)
			
			// Check that all returned events are unprocessed
			for _, event := range events {
				assert.Nil(t, event.ProcessedAt)
			}
		})
	}
}

func TestEventRepository_MarkAsProcessed(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64)
		wantErr bool
	}{
		{
			name: "성공적인_처리_완료_표시",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain, node, and event
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				event := &models.NodeEvent{
					NodeID:    int64(node.ID),
					EventType: "created",
					EventData: &models.EventData{
						EventID:   "evt-123",
						NodeID:    int64(node.ID),
						EventType: "created",
						Timestamp: time.Now(),
					},
				}
				err := repo.Create(event)
				require.NoError(t, err)
				
				return testDB, repo, event.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_이벤트_처리",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository, int64) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				return testDB, repo, 999
			},
			wantErr: false, // SQLite doesn't error on UPDATE with no matching rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			err := repo.MarkAsProcessed(id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// For valid IDs, verify the event was marked as processed
				if id != 999 {
					event, err := repo.GetByID(id)
					assert.NoError(t, err)
					assert.NotNil(t, event)
					assert.NotNil(t, event.ProcessedAt)
				}
			}
		})
	}
}

func TestEventRepository_GetByTypeAndDateRange(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository)
		eventType   string
		start       time.Time
		end         time.Time
		expectedLen int
	}{
		{
			name: "타입과_날짜범위로_이벤트_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				baseTime := time.Now().Truncate(time.Hour)
				
				// Create events with different types and times
				events := []struct {
					eventType string
					when      time.Time
				}{
					{"created", baseTime.Add(-2 * time.Hour)},
					{"updated", baseTime.Add(-1 * time.Hour)},
					{"created", baseTime.Add(-30 * time.Minute)},
					{"deleted", baseTime.Add(-15 * time.Minute)},
					{"created", baseTime.Add(1 * time.Hour)}, // Outside range
				}
				
				for i, e := range events {
					event := &models.NodeEvent{
						NodeID:    int64(node.ID),
						EventType: e.eventType,
						EventData: &models.EventData{
							EventID:   "evt-" + string(rune('1'+i)),
							NodeID:    int64(node.ID),
							EventType: e.eventType,
							Timestamp: e.when,
						},
					}
					
					// Manually set the occurred_at time for testing
					query := `
						INSERT INTO node_events (node_id, event_type, event_data, occurred_at)
						VALUES (?, ?, ?, ?)
					`
					result, err := testDB.SqlxDB().Exec(query, event.NodeID, event.EventType, event.EventData, e.when)
					require.NoError(t, err)
					
					id, err := result.LastInsertId()
					require.NoError(t, err)
					event.ID = id
				}
				
				return testDB, repo
			},
			eventType:   "created",
			start:       time.Now().Truncate(time.Hour).Add(-3 * time.Hour),
			end:         time.Now().Truncate(time.Hour),
			expectedLen: 2, // Two 'created' events within the range
		},
		{
			name: "해당_타입_이벤트_없음",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				return testDB, repo
			},
			eventType:   "created",
			start:       time.Now().Add(-1 * time.Hour),
			end:         time.Now(),
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			events, err := repo.GetByTypeAndDateRange(tt.eventType, tt.start, tt.end)

			assert.NoError(t, err)
			assert.Len(t, events, tt.expectedLen)
			
			// Check that all returned events are of the correct type
			for _, event := range events {
				assert.Equal(t, tt.eventType, event.EventType)
			}
		})
	}
}

func TestEventRepository_DeleteOldEvents(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository)
		olderThan      time.Duration
		expectedDeleted int64
	}{
		{
			name: "오래된_처리완료_이벤트_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				baseTime := time.Now()
				
				// Create old processed events, old unprocessed events, and recent events
				events := []struct {
					when      time.Time
					processed bool
				}{
					{baseTime.Add(-3 * time.Hour), true},  // Should be deleted
					{baseTime.Add(-2 * time.Hour), true},  // Should be deleted
					{baseTime.Add(-3 * time.Hour), false}, // Should NOT be deleted (unprocessed)
					{baseTime.Add(-30 * time.Minute), true}, // Should NOT be deleted (recent)
				}
				
				for i, e := range events {
					// Insert event with specific occurred_at time
					query := `
						INSERT INTO node_events (node_id, event_type, event_data, occurred_at, processed_at)
						VALUES (?, ?, ?, ?, ?)
					`
					var processedAt interface{}
					if e.processed {
						processedAt = e.when.Add(10 * time.Minute)
					}
					
					_, err := testDB.SqlxDB().Exec(query, 
						int64(node.ID), 
						"created", 
						&models.EventData{
							EventID:   "evt-" + string(rune('1'+i)),
							NodeID:    int64(node.ID),
							EventType: "created",
							Timestamp: e.when,
						}, 
						e.when,
						processedAt,
					)
					require.NoError(t, err)
				}
				
				return testDB, repo
			},
			olderThan:       time.Hour,
			expectedDeleted: 2,
		},
		{
			name: "삭제할_이벤트_없음",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				return testDB, repo
			},
			olderThan:       time.Hour,
			expectedDeleted: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			deleted, err := repo.DeleteOldEvents(tt.olderThan)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedDeleted, deleted)
		})
	}
}

func TestEventRepository_GetEventStats(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository)
	}{
		{
			name: "이벤트_통계_조회",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				
				// Create domain and node
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				
				nodeRepo := repositories.NewSQLiteNodeRepository(testDB.DB)
				node := repositories.CreateTestNode(t, nodeRepo, domain.ID)
				
				// Create various events
				eventTypes := []string{"created", "updated", "created", "deleted", "updated"}
				for i, eventType := range eventTypes {
					event := &models.NodeEvent{
						NodeID:    int64(node.ID),
						EventType: eventType,
						EventData: &models.EventData{
							EventID:   "evt-" + string(rune('1'+i)),
							NodeID:    int64(node.ID),
							EventType: eventType,
							Timestamp: time.Now(),
						},
					}
					err := repo.Create(event)
					require.NoError(t, err)
					
					// Mark first 3 as processed
					if i < 3 {
						err = repo.MarkAsProcessed(event.ID)
						require.NoError(t, err)
					}
				}
				
				return testDB, repo
			},
		},
		{
			name: "이벤트_없음_통계",
			setup: func(t *testing.T) (*repositories.TestDB, *repositories.EventRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewEventRepository(testDB.SqlxDB())
				return testDB, repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			stats, err := repo.GetEventStats()

			assert.NoError(t, err)
			assert.NotNil(t, stats)
			
			// Check required fields
			assert.Contains(t, stats, "total_events")
			assert.Contains(t, stats, "pending_events")
			assert.Contains(t, stats, "events_by_type")
			
			totalEvents := stats["total_events"].(int)
			pendingEvents := stats["pending_events"].(int)
			eventsByType := stats["events_by_type"].(map[string]int)
			
			if tt.name == "이벤트_통계_조회" {
				assert.Equal(t, 5, totalEvents)
				assert.Equal(t, 2, pendingEvents)
				assert.Equal(t, 2, eventsByType["created"])
				assert.Equal(t, 2, eventsByType["updated"])
				assert.Equal(t, 1, eventsByType["deleted"])
			} else {
				assert.Equal(t, 0, totalEvents)
				assert.Equal(t, 0, pendingEvents)
				assert.Empty(t, eventsByType)
			}
		})
	}
}