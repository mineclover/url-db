package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-db/internal/models"
	"url-db/internal/repositories"
)

func TestAttributeRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository)
		attribute *models.Attribute
		wantErr   bool
		errType   error
	}{
		{
			name: "성공적인_속성_생성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				// Create domain first
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain // domain ID will be 1
				
				return testDB, repo
			},
			attribute: &models.Attribute{
				DomainID:    1,
				Name:        "test-attr",
				Type:        models.AttributeTypeString,
				Description: "Test attribute",
			},
			wantErr: false,
		},
		{
			name: "중복_속성_이름",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				// Create domain first
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				
				// Create existing attribute
				existing := &models.Attribute{
					DomainID:    1,
					Name:        "test-attr",
					Type:        models.AttributeTypeString,
					Description: "Existing attribute",
				}
				err := repo.Create(existing)
				require.NoError(t, err)
				
				return testDB, repo
			},
			attribute: &models.Attribute{
				DomainID:    1,
				Name:        "test-attr",
				Type:        models.AttributeTypeString,
				Description: "Duplicate attribute",
			},
			wantErr: true,
			errType: repositories.ErrDuplicateEntry,
		},
		{
			name: "존재하지_않는_도메인",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				return testDB, repo
			},
			attribute: &models.Attribute{
				DomainID:    999,
				Name:        "test-attr",
				Type:        models.AttributeTypeString,
				Description: "Test attribute",
			},
			wantErr: true,
			errType: repositories.ErrForeignKeyConstraint,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			err := repo.Create(tt.attribute)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.attribute.ID)
				assert.False(t, tt.attribute.CreatedAt.IsZero())
			}
		})
	}
}

func TestAttributeRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, int)
		id      int
		wantErr bool
		errType error
	}{
		{
			name: "존재하는_속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				// Create domain and attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				attr := &models.Attribute{
					DomainID:    1,
					Name:        "test-attr",
					Type:        models.AttributeTypeString,
					Description: "Test attribute",
				}
				err := repo.Create(attr)
				require.NoError(t, err)
				
				return testDB, repo, attr.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				return testDB, repo, 999
			},
			id:      999,
			wantErr: true,
			errType: repositories.ErrAttributeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, id := tt.setup(t)
			defer testDB.Close()

			if tt.id != 0 {
				id = tt.id
			}

			attribute, err := repo.GetByID(id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, attribute)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, attribute)
				assert.Equal(t, id, attribute.ID)
				assert.Equal(t, "test-attr", attribute.Name)
				assert.Equal(t, models.AttributeTypeString, attribute.Type)
			}
		})
	}
}

func TestAttributeRepository_GetByDomainAndName(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository)
		domainID int
		attrName string
		wantErr  bool
		errType  error
	}{
		{
			name: "존재하는_속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				// Create domain and attribute
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				attr := &models.Attribute{
					DomainID:    1,
					Name:        "test-attr",
					Type:        models.AttributeTypeString,
					Description: "Test attribute",
				}
				err := repo.Create(attr)
				require.NoError(t, err)
				
				return testDB, repo
			},
			domainID: 1,
			attrName: "test-attr",
			wantErr:  false,
		},
		{
			name: "존재하지_않는_속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				return testDB, repo
			},
			domainID: 1,
			attrName: "non-existent",
			wantErr:  true,
			errType:  repositories.ErrAttributeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			attribute, err := repo.GetByDomainAndName(tt.domainID, tt.attrName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, attribute)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, attribute)
				assert.Equal(t, tt.domainID, attribute.DomainID)
				assert.Equal(t, tt.attrName, attribute.Name)
			}
		})
	}
}

func TestAttributeRepository_ListByDomain(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository)
		domainID    int
		expectedLen int
	}{
		{
			name: "도메인의_모든_속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				// Create domain
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				
				// Create multiple attributes
				attrs := []*models.Attribute{
					{DomainID: 1, Name: "attr1", Type: models.AttributeTypeString, Description: "Attribute 1"},
					{DomainID: 1, Name: "attr2", Type: models.AttributeTypeNumber, Description: "Attribute 2"},
					{DomainID: 1, Name: "attr3", Type: models.AttributeTypeTag, Description: "Attribute 3"},
				}
				
				for _, attr := range attrs {
					err := repo.Create(attr)
					require.NoError(t, err)
				}
				
				return testDB, repo
			},
			domainID:    1,
			expectedLen: 3,
		},
		{
			name: "빈_도메인_속성_조회",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				return testDB, repo
			},
			domainID:    1,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			attributes, err := repo.ListByDomain(tt.domainID)

			assert.NoError(t, err)
			assert.Len(t, attributes, tt.expectedLen)
			
			// Check ordering (should be by name)
			if len(attributes) > 1 {
				for i := 0; i < len(attributes)-1; i++ {
					assert.True(t, attributes[i].Name <= attributes[i+1].Name)
				}
			}
		})
	}
}

func TestAttributeRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, *models.Attribute)
		updateFn  func(*models.Attribute)
		wantErr   bool
		errType   error
	}{
		{
			name: "성공적인_속성_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, *models.Attribute) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				attr := &models.Attribute{
					DomainID:    1,
					Name:        "test-attr",
					Type:        models.AttributeTypeString,
					Description: "Original description",
				}
				err := repo.Create(attr)
				require.NoError(t, err)
				
				return testDB, repo, attr
			},
			updateFn: func(attr *models.Attribute) {
				attr.Name = "updated-attr"
				attr.Type = models.AttributeTypeNumber
				attr.Description = "Updated description"
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_속성_업데이트",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, *models.Attribute) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				return testDB, repo, &models.Attribute{ID: 999}
			},
			updateFn: func(attr *models.Attribute) {
				attr.Name = "non-existent"
			},
			wantErr: true,
			errType: repositories.ErrAttributeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo, attr := tt.setup(t)
			defer testDB.Close()

			tt.updateFn(attr)
			err := repo.Update(attr)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				
				// Verify update
				updated, err := repo.GetByID(attr.ID)
				assert.NoError(t, err)
				assert.Equal(t, attr.Name, updated.Name)
				assert.Equal(t, attr.Type, updated.Type)
				assert.Equal(t, attr.Description, updated.Description)
			}
		})
	}
}

func TestAttributeRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, int)
		wantErr bool
		errType error
	}{
		{
			name: "성공적인_속성_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				attr := &models.Attribute{
					DomainID:    1,
					Name:        "test-attr",
					Type:        models.AttributeTypeString,
					Description: "Test attribute",
				}
				err := repo.Create(attr)
				require.NoError(t, err)
				
				return testDB, repo, attr.ID
			},
			wantErr: false,
		},
		{
			name: "존재하지_않는_속성_삭제",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository, int) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				return testDB, repo, 999
			},
			wantErr: true,
			errType: repositories.ErrAttributeNotFound,
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
				assert.ErrorIs(t, err, repositories.ErrAttributeNotFound)
			}
		})
	}
}

func TestAttributeRepository_ExistsByDomainAndName(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository)
		domainID int
		attrName string
		expected bool
	}{
		{
			name: "존재하는_속성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				attr := &models.Attribute{
					DomainID:    1,
					Name:        "test-attr",
					Type:        models.AttributeTypeString,
					Description: "Test attribute",
				}
				err := repo.Create(attr)
				require.NoError(t, err)
				
				return testDB, repo
			},
			domainID: 1,
			attrName: "test-attr",
			expected: true,
		},
		{
			name: "존재하지_않는_속성",
			setup: func(t *testing.T) (*repositories.TestDB, repositories.AttributeRepository) {
				testDB := repositories.SetupTestDB(t)
				repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
				domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
				domain := repositories.CreateTestDomain(t, domainRepo)
				_ = domain
				return testDB, repo
			},
			domainID: 1,
			attrName: "non-existent",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB, repo := tt.setup(t)
			defer testDB.Close()

			exists, err := repo.ExistsByDomainAndName(tt.domainID, tt.attrName)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestAttributeRepository_TransactionMethods(t *testing.T) {
	t.Run("CreateTx_성공", func(t *testing.T) {
		testDB := repositories.SetupTestDB(t)
		defer testDB.Close()
		
		repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
		domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
		domain := repositories.CreateTestDomain(t, domainRepo)
		_ = domain
		
		tx, err := testDB.DB.Begin()
		require.NoError(t, err)
		defer tx.Rollback()
		
		attr := &models.Attribute{
			DomainID:    1,
			Name:        "tx-attr",
			Type:        models.AttributeTypeString,
			Description: "Transaction attribute",
		}
		
		err = repo.CreateTx(tx, attr)
		assert.NoError(t, err)
		assert.NotZero(t, attr.ID)
		assert.False(t, attr.CreatedAt.IsZero())
		
		err = tx.Commit()
		assert.NoError(t, err)
		
		// Verify creation
		created, err := repo.GetByID(attr.ID)
		assert.NoError(t, err)
		assert.Equal(t, attr.Name, created.Name)
	})
	
	t.Run("UpdateTx_성공", func(t *testing.T) {
		testDB := repositories.SetupTestDB(t)
		defer testDB.Close()
		
		repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
		domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
		domain := repositories.CreateTestDomain(t, domainRepo)
		_ = domain
		
		// Create attribute first
		attr := &models.Attribute{
			DomainID:    1,
			Name:        "original-attr",
			Type:        models.AttributeTypeString,
			Description: "Original description",
		}
		err := repo.Create(attr)
		require.NoError(t, err)
		
		// Update in transaction
		tx, err := testDB.DB.Begin()
		require.NoError(t, err)
		defer tx.Rollback()
		
		attr.Name = "updated-attr"
		attr.Description = "Updated description"
		
		err = repo.UpdateTx(tx, attr)
		assert.NoError(t, err)
		
		err = tx.Commit()
		assert.NoError(t, err)
		
		// Verify update
		updated, err := repo.GetByID(attr.ID)
		assert.NoError(t, err)
		assert.Equal(t, "updated-attr", updated.Name)
		assert.Equal(t, "Updated description", updated.Description)
	})
	
	t.Run("DeleteTx_성공", func(t *testing.T) {
		testDB := repositories.SetupTestDB(t)
		defer testDB.Close()
		
		repo := repositories.NewSQLiteAttributeRepository(testDB.DB)
		domainRepo := repositories.NewSQLiteDomainRepository(testDB.DB)
		domain := repositories.CreateTestDomain(t, domainRepo)
		_ = domain
		
		// Create attribute first
		attr := &models.Attribute{
			DomainID:    1,
			Name:        "delete-attr",
			Type:        models.AttributeTypeString,
			Description: "To be deleted",
		}
		err := repo.Create(attr)
		require.NoError(t, err)
		
		// Delete in transaction
		tx, err := testDB.DB.Begin()
		require.NoError(t, err)
		defer tx.Rollback()
		
		err = repo.DeleteTx(tx, attr.ID)
		assert.NoError(t, err)
		
		err = tx.Commit()
		assert.NoError(t, err)
		
		// Verify deletion
		_, err = repo.GetByID(attr.ID)
		assert.ErrorIs(t, err, repositories.ErrAttributeNotFound)
	})
}