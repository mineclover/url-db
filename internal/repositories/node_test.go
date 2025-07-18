package repositories

import (
	"testing"
	"url-db/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeRepository_Create(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, domainRepo)

	tests := []struct {
		name    string
		node    *models.Node
		wantErr bool
	}{
		{
			name: "성공적인 노드 생성",
			node: &models.Node{
				Content:     "https://example.com",
				DomainID:    domain.ID,
				Title:       "Test Node",
				Description: "Test description",
			},
			wantErr: false,
		},
		{
			name: "중복 노드 생성",
			node: &models.Node{
				Content:     "https://example.com",
				DomainID:    domain.ID,
				Title:       "Duplicate Node",
				Description: "Duplicate description",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := nodeRepo.Create(tt.node)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tt.node.ID)
			assert.NotZero(t, tt.node.CreatedAt)
			assert.NotZero(t, tt.node.UpdatedAt)
		})
	}
}

func TestNodeRepository_GetByID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, domainRepo)
	node := CreateTestNode(t, nodeRepo, domain.ID)

	tests := []struct {
		name    string
		id      int
		want    *models.Node
		wantErr error
	}{
		{
			name: "존재하는 노드 조회",
			id:   node.ID,
			want: node,
		},
		{
			name:    "존재하지 않는 노드 조회",
			id:      9999,
			wantErr: ErrNodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nodeRepo.GetByID(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			AssertNodeEqual(t, tt.want, got)
		})
	}
}

func TestNodeRepository_GetByDomainAndContent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, domainRepo)
	node := CreateTestNode(t, nodeRepo, domain.ID)

	tests := []struct {
		name     string
		domainID int
		content  string
		want     *models.Node
		wantErr  error
	}{
		{
			name:     "존재하는 노드 조회",
			domainID: domain.ID,
			content:  node.Content,
			want:     node,
		},
		{
			name:     "존재하지 않는 노드 조회",
			domainID: domain.ID,
			content:  "https://nonexistent.com",
			wantErr:  ErrNodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nodeRepo.GetByDomainAndContent(tt.domainID, tt.content)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			AssertNodeEqual(t, tt.want, got)
		})
	}
}

func TestNodeRepository_ListByDomain(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인들 생성
	domain1 := CreateTestDomain(t, domainRepo)
	domain2 := NewTestDomainBuilder().WithName("domain2").Build()
	require.NoError(t, domainRepo.Create(domain2))

	// 도메인1에 노드들 생성
	node1 := CreateTestNode(t, nodeRepo, domain1.ID)
	node2 := NewTestNodeBuilder().WithDomainID(domain1.ID).WithContent("https://example2.com").Build()
	require.NoError(t, nodeRepo.Create(node2))

	// 도메인2에 노드 생성
	node3 := NewTestNodeBuilder().WithDomainID(domain2.ID).WithContent("https://example3.com").Build()
	require.NoError(t, nodeRepo.Create(node3))

	tests := []struct {
		name      string
		domainID  int
		offset    int
		limit     int
		wantCount int
		wantTotal int
	}{
		{
			name:      "도메인1의 모든 노드 조회",
			domainID:  domain1.ID,
			offset:    0,
			limit:     10,
			wantCount: 2,
			wantTotal: 2,
		},
		{
			name:      "도메인2의 모든 노드 조회",
			domainID:  domain2.ID,
			offset:    0,
			limit:     10,
			wantCount: 1,
			wantTotal: 1,
		},
		{
			name:      "페이지네이션 테스트",
			domainID:  domain1.ID,
			offset:    0,
			limit:     1,
			wantCount: 1,
			wantTotal: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, total, err := nodeRepo.ListByDomain(tt.domainID, tt.offset, tt.limit)

			require.NoError(t, err)
			assert.Len(t, nodes, tt.wantCount)
			assert.Equal(t, tt.wantTotal, total)

			// 모든 노드가 해당 도메인의 것인지 확인
			for _, node := range nodes {
				assert.Equal(t, tt.domainID, node.DomainID)
			}
		})
	}
}

func TestNodeRepository_Search(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, domainRepo)

	// 검색용 테스트 노드들 생성
	node1 := NewTestNodeBuilder().
		WithDomainID(domain.ID).
		WithContent("https://golang.org").
		WithTitle("Go Programming Language").
		WithDescription("Official Go website").
		Build()
	require.NoError(t, nodeRepo.Create(node1))

	node2 := NewTestNodeBuilder().
		WithDomainID(domain.ID).
		WithContent("https://python.org").
		WithTitle("Python Programming").
		WithDescription("Official Python website").
		Build()
	require.NoError(t, nodeRepo.Create(node2))

	tests := []struct {
		name      string
		domainID  int
		query     string
		offset    int
		limit     int
		wantCount int
		wantTotal int
	}{
		{
			name:      "제목으로 검색",
			domainID:  domain.ID,
			query:     "Go",
			offset:    0,
			limit:     10,
			wantCount: 1,
			wantTotal: 1,
		},
		{
			name:      "콘텐츠로 검색",
			domainID:  domain.ID,
			query:     "golang",
			offset:    0,
			limit:     10,
			wantCount: 1,
			wantTotal: 1,
		},
		{
			name:      "설명으로 검색",
			domainID:  domain.ID,
			query:     "Official",
			offset:    0,
			limit:     10,
			wantCount: 2,
			wantTotal: 2,
		},
		{
			name:      "검색 결과 없음",
			domainID:  domain.ID,
			query:     "nonexistent",
			offset:    0,
			limit:     10,
			wantCount: 0,
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, total, err := nodeRepo.Search(tt.domainID, tt.query, tt.offset, tt.limit)

			require.NoError(t, err)
			assert.Len(t, nodes, tt.wantCount)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}

func TestNodeRepository_Update(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, domainRepo)
	node := CreateTestNode(t, nodeRepo, domain.ID)

	tests := []struct {
		name    string
		node    *models.Node
		wantErr error
	}{
		{
			name: "성공적인 노드 업데이트",
			node: &models.Node{
				ID:          node.ID,
				Content:     "https://updated.com",
				Title:       "Updated Title",
				Description: "Updated description",
			},
		},
		{
			name: "존재하지 않는 노드 업데이트",
			node: &models.Node{
				ID:          9999,
				Content:     "https://nonexistent.com",
				Title:       "Nonexistent",
				Description: "Nonexistent description",
			},
			wantErr: ErrNodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := nodeRepo.Update(tt.node)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, tt.node.UpdatedAt)

			// 업데이트된 노드 조회하여 확인
			updated, err := nodeRepo.GetByID(tt.node.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.node.Content, updated.Content)
			assert.Equal(t, tt.node.Title, updated.Title)
			assert.Equal(t, tt.node.Description, updated.Description)
		})
	}
}

func TestNodeRepository_Delete(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, domainRepo)
	node := CreateTestNode(t, nodeRepo, domain.ID)

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name: "성공적인 노드 삭제",
			id:   node.ID,
		},
		{
			name:    "존재하지 않는 노드 삭제",
			id:      9999,
			wantErr: ErrNodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := nodeRepo.Delete(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)

			// 삭제된 노드 조회 시 에러 확인
			_, err = nodeRepo.GetByID(tt.id)
			assert.Equal(t, ErrNodeNotFound, err)
		})
	}
}

func TestNodeRepository_ExistsByDomainAndContent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, domainRepo)
	node := CreateTestNode(t, nodeRepo, domain.ID)

	tests := []struct {
		name     string
		domainID int
		content  string
		want     bool
	}{
		{
			name:     "존재하는 노드",
			domainID: domain.ID,
			content:  node.Content,
			want:     true,
		},
		{
			name:     "존재하지 않는 노드",
			domainID: domain.ID,
			content:  "https://nonexistent.com",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := nodeRepo.ExistsByDomainAndContent(tt.domainID, tt.content)

			require.NoError(t, err)
			assert.Equal(t, tt.want, exists)
		})
	}
}

func TestNodeRepository_BatchOperations(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close()

	domainRepo := NewSQLiteDomainRepository(testDB.DB)
	nodeRepo := NewSQLiteNodeRepository(testDB.DB)

	// 테스트 도메인 생성
	domain := CreateTestDomain(t, domainRepo)

	t.Run("배치 생성", func(t *testing.T) {
		nodes := []models.Node{
			{
				Content:     "https://batch1.com",
				DomainID:    domain.ID,
				Title:       "Batch Node 1",
				Description: "First batch node",
			},
			{
				Content:     "https://batch2.com",
				DomainID:    domain.ID,
				Title:       "Batch Node 2",
				Description: "Second batch node",
			},
		}

		err := nodeRepo.BatchCreate(nodes)
		require.NoError(t, err)

		// 노드들이 생성되었는지 확인
		for _, node := range nodes {
			exists, err := nodeRepo.ExistsByDomainAndContent(domain.ID, node.Content)
			require.NoError(t, err)
			assert.True(t, exists)
		}
	})

	t.Run("배치 삭제", func(t *testing.T) {
		// 테스트 노드들 생성
		node1 := CreateTestNode(t, nodeRepo, domain.ID)
		node2 := NewTestNodeBuilder().WithDomainID(domain.ID).WithContent("https://delete2.com").Build()
		require.NoError(t, nodeRepo.Create(node2))

		// 배치 삭제
		err := nodeRepo.BatchDelete([]int{node1.ID, node2.ID})
		require.NoError(t, err)

		// 노드들이 삭제되었는지 확인
		_, err = nodeRepo.GetByID(node1.ID)
		assert.Equal(t, ErrNodeNotFound, err)
		_, err = nodeRepo.GetByID(node2.ID)
		assert.Equal(t, ErrNodeNotFound, err)
	})
}
