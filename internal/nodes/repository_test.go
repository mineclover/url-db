package nodes

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-db/internal/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	schema := `
		CREATE TABLE domains (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE nodes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			domain_id INTEGER NOT NULL,
			title TEXT,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE,
			UNIQUE(content, domain_id)
		);
		
		CREATE INDEX idx_nodes_domain ON nodes(domain_id);
		CREATE INDEX idx_nodes_content ON nodes(content);
		
		INSERT INTO domains (name, description) VALUES ('test-domain', 'Test domain');
	`

	_, err = db.Exec(schema)
	require.NoError(t, err)

	return db
}

func TestSQLiteNodeRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	node := &models.Node{
		Content:     "https://example.com",
		DomainID:    1,
		Title:       "Example",
		Description: "Test node",
	}

	err := repo.Create(node)
	assert.NoError(t, err)
	assert.NotZero(t, node.ID)
	assert.NotZero(t, node.CreatedAt)
	assert.NotZero(t, node.UpdatedAt)
}

func TestSQLiteNodeRepository_Create_DuplicateURL(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	node1 := &models.Node{
		Content:  "https://example.com",
		DomainID: 1,
		Title:    "Example 1",
	}

	node2 := &models.Node{
		Content:  "https://example.com",
		DomainID: 1,
		Title:    "Example 2",
	}

	err := repo.Create(node1)
	assert.NoError(t, err)

	err = repo.Create(node2)
	assert.Equal(t, ErrNodeAlreadyExists, err)
}

func TestSQLiteNodeRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	original := &models.Node{
		Content:     "https://example.com",
		DomainID:    1,
		Title:       "Example",
		Description: "Test node",
	}

	err := repo.Create(original)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(original.ID)
	assert.NoError(t, err)
	assert.Equal(t, original.Content, retrieved.Content)
	assert.Equal(t, original.DomainID, retrieved.DomainID)
	assert.Equal(t, original.Title, retrieved.Title)
	assert.Equal(t, original.Description, retrieved.Description)
}

func TestSQLiteNodeRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	_, err := repo.GetByID(999)
	assert.Equal(t, ErrNodeNotFound, err)
}

func TestSQLiteNodeRepository_GetByURL(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	original := &models.Node{
		Content:  "https://example.com",
		DomainID: 1,
		Title:    "Example",
	}

	err := repo.Create(original)
	require.NoError(t, err)

	retrieved, err := repo.GetByURL(1, "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, original.Content, retrieved.Content)
	assert.Equal(t, original.DomainID, retrieved.DomainID)
}

func TestSQLiteNodeRepository_GetByURL_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	_, err := repo.GetByURL(1, "https://nonexistent.com")
	assert.Equal(t, ErrNodeNotFound, err)
}

func TestSQLiteNodeRepository_GetByDomainID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	// Create multiple nodes
	for i := 0; i < 5; i++ {
		node := &models.Node{
			Content:  "https://example.com/" + string(rune('a'+i)),
			DomainID: 1,
			Title:    "Example " + string(rune('A'+i)),
		}
		err := repo.Create(node)
		require.NoError(t, err)
	}

	nodes, totalCount, err := repo.GetByDomainID(1, 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, 5, totalCount)
	assert.Len(t, nodes, 3)
}

func TestSQLiteNodeRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	original := &models.Node{
		Content:     "https://example.com",
		DomainID:    1,
		Title:       "Original Title",
		Description: "Original Description",
	}

	err := repo.Create(original)
	require.NoError(t, err)

	original.Title = "Updated Title"
	original.Description = "Updated Description"

	err = repo.Update(original)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(original.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", retrieved.Title)
	assert.Equal(t, "Updated Description", retrieved.Description)
}

func TestSQLiteNodeRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	node := &models.Node{
		ID:          999,
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	err := repo.Update(node)
	assert.Equal(t, ErrNodeNotFound, err)
}

func TestSQLiteNodeRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	node := &models.Node{
		Content:  "https://example.com",
		DomainID: 1,
		Title:    "Example",
	}

	err := repo.Create(node)
	require.NoError(t, err)

	err = repo.Delete(node.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(node.ID)
	assert.Equal(t, ErrNodeNotFound, err)
}

func TestSQLiteNodeRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	err := repo.Delete(999)
	assert.Equal(t, ErrNodeNotFound, err)
}

func TestSQLiteNodeRepository_Search(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	// Create test nodes
	testNodes := []models.Node{
		{Content: "https://golang.org", DomainID: 1, Title: "Go Programming", Description: "Go language docs"},
		{Content: "https://rust-lang.org", DomainID: 1, Title: "Rust Programming", Description: "Rust language docs"},
		{Content: "https://python.org", DomainID: 1, Title: "Python Programming", Description: "Python language docs"},
	}

	for _, node := range testNodes {
		n := node
		err := repo.Create(&n)
		require.NoError(t, err)
	}

	// Search for "Go"
	nodes, totalCount, err := repo.Search(1, "Go", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, totalCount)
	assert.Len(t, nodes, 1)
	assert.Equal(t, "Go Programming", nodes[0].Title)

	// Search for "Programming"
	nodes, totalCount, err = repo.Search(1, "Programming", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 3, totalCount)
	assert.Len(t, nodes, 3)
}

func TestSQLiteNodeRepository_CheckDomainExists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteNodeRepository(db)

	exists, err := repo.CheckDomainExists(1)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.CheckDomainExists(999)
	assert.NoError(t, err)
	assert.False(t, exists)
}
