package repositories

import (
	"testing"
	"url-db/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainRepository_Create(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	repo := NewSQLiteDomainRepository(testDB.DB)

	tests := []struct {
		name    string
		domain  *models.Domain
		wantErr bool
	}{
		{
			name: "성공적인 도메인 생성",
			domain: &models.Domain{
				Name:        "test-domain",
				Description: "Test domain description",
			},
			wantErr: false,
		},
		{
			name: "중복 도메인 이름",
			domain: &models.Domain{
				Name:        "test-domain",
				Description: "Duplicate domain",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.domain)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tt.domain.ID)
			assert.NotZero(t, tt.domain.CreatedAt)
			assert.NotZero(t, tt.domain.UpdatedAt)
		})
	}
}

func TestDomainRepository_GetByID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	repo := NewSQLiteDomainRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, repo)

	tests := []struct {
		name    string
		id      int
		want    *models.Domain
		wantErr error
	}{
		{
			name: "존재하는 도메인 조회",
			id:   domain.ID,
			want: domain,
		},
		{
			name:    "존재하지 않는 도메인 조회",
			id:      9999,
			wantErr: ErrDomainNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			AssertDomainEqual(t, tt.want, got)
		})
	}
}

func TestDomainRepository_GetByName(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	repo := NewSQLiteDomainRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, repo)

	tests := []struct {
		name       string
		domainName string
		want       *models.Domain
		wantErr    error
	}{
		{
			name:       "존재하는 도메인 이름으로 조회",
			domainName: domain.Name,
			want:       domain,
		},
		{
			name:       "존재하지 않는 도메인 이름으로 조회",
			domainName: "nonexistent-domain",
			wantErr:    ErrDomainNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByName(tt.domainName)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			AssertDomainEqual(t, tt.want, got)
		})
	}
}

func TestDomainRepository_List(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	repo := NewSQLiteDomainRepository(testDB.DB)

	// 테스트 도메인들 생성
	domain1 := CreateTestDomain(t, repo)
	domain2 := NewTestDomainBuilder().WithName("domain2").Build()
	require.NoError(t, repo.Create(domain2))
	domain3 := NewTestDomainBuilder().WithName("domain3").Build()
	require.NoError(t, repo.Create(domain3))

	tests := []struct {
		name      string
		offset    int
		limit     int
		wantCount int
		wantTotal int
	}{
		{
			name:      "모든 도메인 조회",
			offset:    0,
			limit:     10,
			wantCount: 3,
			wantTotal: 3,
		},
		{
			name:      "페이지네이션 - 첫 번째 페이지",
			offset:    0,
			limit:     2,
			wantCount: 2,
			wantTotal: 3,
		},
		{
			name:      "페이지네이션 - 두 번째 페이지",
			offset:    2,
			limit:     2,
			wantCount: 1,
			wantTotal: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domains, total, err := repo.List(tt.offset, tt.limit)

			require.NoError(t, err)
			assert.Len(t, domains, tt.wantCount)
			assert.Equal(t, tt.wantTotal, total)

			// 첫 번째 도메인이 가장 최근 생성된 것인지 확인 (created_at DESC)
			if len(domains) > 0 {
				assert.True(t, domains[0].CreatedAt.After(domain1.CreatedAt) ||
					domains[0].CreatedAt.Equal(domain1.CreatedAt))
			}
		})
	}
}

func TestDomainRepository_Update(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	repo := NewSQLiteDomainRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, repo)

	tests := []struct {
		name    string
		domain  *models.Domain
		wantErr error
	}{
		{
			name: "성공적인 도메인 업데이트",
			domain: &models.Domain{
				ID:          domain.ID,
				Name:        "updated-domain",
				Description: "Updated description",
			},
		},
		{
			name: "존재하지 않는 도메인 업데이트",
			domain: &models.Domain{
				ID:          9999,
				Name:        "nonexistent-domain",
				Description: "Description",
			},
			wantErr: ErrDomainNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.domain)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tt.domain.UpdatedAt)

			// 업데이트된 도메인 조회하여 확인
			updated, err := repo.GetByID(tt.domain.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.domain.Name, updated.Name)
			assert.Equal(t, tt.domain.Description, updated.Description)
		})
	}
}

func TestDomainRepository_Delete(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	repo := NewSQLiteDomainRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, repo)

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name: "성공적인 도메인 삭제",
			id:   domain.ID,
		},
		{
			name:    "존재하지 않는 도메인 삭제",
			id:      9999,
			wantErr: ErrDomainNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)

			// 삭제된 도메인 조회 시 에러 확인
			_, err = repo.GetByID(tt.id)
			assert.Equal(t, ErrDomainNotFound, err)
		})
	}
}

func TestDomainRepository_ExistsByName(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	repo := NewSQLiteDomainRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, repo)

	tests := []struct {
		name       string
		domainName string
		want       bool
	}{
		{
			name:       "존재하는 도메인 이름",
			domainName: domain.Name,
			want:       true,
		},
		{
			name:       "존재하지 않는 도메인 이름",
			domainName: "nonexistent-domain",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistsByName(tt.domainName)

			require.NoError(t, err)
			assert.Equal(t, tt.want, exists)
		})
	}
}

func TestDomainRepository_Transaction(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	repo := NewSQLiteDomainRepository(testDB.DB)

	t.Run("트랜잭션 커밋", func(t *testing.T) {
		domain := &models.Domain{
			Name:        "tx-domain",
			Description: "Transaction test domain",
		}

		err := repo.(*sqliteDomainRepository).WithTransaction(func(tx *sql.Tx) error {
			return repo.CreateTx(tx, domain)
		})

		require.NoError(t, err)
		assert.NotZero(t, domain.ID)

		// 트랜잭션이 커밋되었는지 확인
		found, err := repo.GetByID(domain.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.Name, found.Name)
	})

	t.Run("트랜잭션 롤백", func(t *testing.T) {
		domain := &models.Domain{
			Name:        "rollback-domain",
			Description: "Rollback test domain",
		}

		err := repo.(*sqliteDomainRepository).WithTransaction(func(tx *sql.Tx) error {
			if err := repo.CreateTx(tx, domain); err != nil {
				return err
			}
			// 강제로 에러 발생
			return ErrDuplicateEntry
		})

		require.Error(t, err)
		assert.Equal(t, ErrDuplicateEntry, err)

		// 트랜잭션이 롤백되었는지 확인
		_, err = repo.GetByName(domain.Name)
		assert.Equal(t, ErrDomainNotFound, err)
	})
}
