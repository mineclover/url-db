package nodes

import (
	"database/sql"
	"fmt"
	"time"

	"url-db/internal/models"
)

type NodeRepository interface {
	Create(node *models.Node) error
	GetByID(id int) (*models.Node, error)
	GetByDomainID(domainID, page, size int) ([]models.Node, int, error)
	GetByURL(domainID int, url string) (*models.Node, error)
	Update(node *models.Node) error
	Delete(id int) error
	Search(domainID int, query string, page, size int) ([]models.Node, int, error)
	CheckDomainExists(domainID int) (bool, error)
}

type SQLiteNodeRepository struct {
	db *sql.DB
}

func NewSQLiteNodeRepository(db *sql.DB) NodeRepository {
	return &SQLiteNodeRepository{db: db}
}

func (r *SQLiteNodeRepository) Create(node *models.Node) error {
	query := `
		INSERT INTO nodes (content, domain_id, title, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(query, node.Content, node.DomainID, node.Title, node.Description, now, now)
	if err != nil {
		if isUniqueConstraintError(err) {
			return ErrNodeAlreadyExists
		}
		return fmt.Errorf("failed to create node: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	node.ID = int(id)
	node.CreatedAt = now
	node.UpdatedAt = now

	return nil
}

func (r *SQLiteNodeRepository) GetByID(id int) (*models.Node, error) {
	query := `
		SELECT id, content, domain_id, title, description, created_at, updated_at
		FROM nodes
		WHERE id = ?
	`

	node := &models.Node{}
	err := r.db.QueryRow(query, id).Scan(
		&node.ID,
		&node.Content,
		&node.DomainID,
		&node.Title,
		&node.Description,
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNodeNotFound
		}
		return nil, fmt.Errorf("failed to get node by id: %w", err)
	}

	return node, nil
}

func (r *SQLiteNodeRepository) GetByDomainID(domainID, page, size int) ([]models.Node, int, error) {
	offset := (page - 1) * size

	// Get total count
	countQuery := `SELECT COUNT(*) FROM nodes WHERE domain_id = ?`
	var totalCount int
	err := r.db.QueryRow(countQuery, domainID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get nodes with pagination
	query := `
		SELECT id, content, domain_id, title, description, created_at, updated_at
		FROM nodes
		WHERE domain_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, domainID, size, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get nodes by domain id: %w", err)
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(
			&node.ID,
			&node.Content,
			&node.DomainID,
			&node.Title,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return nodes, totalCount, nil
}

func (r *SQLiteNodeRepository) GetByURL(domainID int, url string) (*models.Node, error) {
	query := `
		SELECT id, content, domain_id, title, description, created_at, updated_at
		FROM nodes
		WHERE domain_id = ? AND content = ?
	`

	node := &models.Node{}
	err := r.db.QueryRow(query, domainID, url).Scan(
		&node.ID,
		&node.Content,
		&node.DomainID,
		&node.Title,
		&node.Description,
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNodeNotFound
		}
		return nil, fmt.Errorf("failed to get node by url: %w", err)
	}

	return node, nil
}

func (r *SQLiteNodeRepository) Update(node *models.Node) error {
	query := `
		UPDATE nodes
		SET title = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := r.db.Exec(query, node.Title, node.Description, now, node.ID)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNodeNotFound
	}

	node.UpdatedAt = now
	return nil
}

func (r *SQLiteNodeRepository) Delete(id int) error {
	query := `DELETE FROM nodes WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNodeNotFound
	}

	return nil
}

func (r *SQLiteNodeRepository) Search(domainID int, query string, page, size int) ([]models.Node, int, error) {
	offset := (page - 1) * size
	searchPattern := "%" + query + "%"

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM nodes 
		WHERE domain_id = ? AND (title LIKE ? OR description LIKE ? OR content LIKE ?)
	`
	var totalCount int
	err := r.db.QueryRow(countQuery, domainID, searchPattern, searchPattern, searchPattern).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get search count: %w", err)
	}

	// Get nodes with search and pagination
	searchQuery := `
		SELECT id, content, domain_id, title, description, created_at, updated_at
		FROM nodes
		WHERE domain_id = ? AND (title LIKE ? OR description LIKE ? OR content LIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(searchQuery, domainID, searchPattern, searchPattern, searchPattern, size, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search nodes: %w", err)
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(
			&node.ID,
			&node.Content,
			&node.DomainID,
			&node.Title,
			&node.Description,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan search result: %w", err)
		}
		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("search rows iteration error: %w", err)
	}

	return nodes, totalCount, nil
}

func (r *SQLiteNodeRepository) CheckDomainExists(domainID int) (bool, error) {
	query := `SELECT COUNT(*) FROM domains WHERE id = ?`
	var count int
	err := r.db.QueryRow(query, domainID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check domain existence: %w", err)
	}
	return count > 0, nil
}

func isUniqueConstraintError(err error) bool {
	return err != nil && (err.Error() == "UNIQUE constraint failed: nodes.content, nodes.domain_id" ||
		err.Error() == "constraint failed: UNIQUE constraint failed: nodes.content, nodes.domain_id")
}
