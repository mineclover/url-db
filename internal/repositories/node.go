package repositories

import (
	"context"
	"database/sql"
	"url-db/internal/models"
)

// sqliteNodeRepository 는 SQLite 기반 노드 리포지토리 구현체입니다.
type sqliteNodeRepository struct {
	*BaseRepository
}

// NewSQLiteNodeRepository 는 새로운 SQLite 노드 리포지토리를 생성합니다.
func NewSQLiteNodeRepository(db *sql.DB) NodeRepository {
	return &sqliteNodeRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 는 새로운 노드를 생성합니다.
func (r *sqliteNodeRepository) Create(node *models.Node) error {
	query := `
		INSERT INTO nodes (content, domain_id, title, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
		RETURNING id, created_at, updated_at
	`
	
	err := r.QueryRow(query, node.Content, node.DomainID, node.Title, node.Description).Scan(
		&node.ID, &node.CreatedAt, &node.UpdatedAt,
	)
	
	if err != nil {
		return MapSQLiteError(err)
	}
	
	return nil
}

// GetByID 는 ID로 노드를 조회합니다.
func (r *sqliteNodeRepository) GetByID(id int) (*models.Node, error) {
	query := `
		SELECT id, content, domain_id, title, description, created_at, updated_at
		FROM nodes
		WHERE id = ?
	`
	
	node := &models.Node{}
	err := r.QueryRow(query, id).Scan(
		&node.ID, &node.Content, &node.DomainID, &node.Title,
		&node.Description, &node.CreatedAt, &node.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrNodeNotFound
	}
	
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	
	return node, nil
}

// GetByDomainAndContent 는 도메인 ID와 콘텐츠로 노드를 조회합니다.
func (r *sqliteNodeRepository) GetByDomainAndContent(domainID int, content string) (*models.Node, error) {
	query := `
		SELECT id, content, domain_id, title, description, created_at, updated_at
		FROM nodes
		WHERE domain_id = ? AND content = ?
	`
	
	node := &models.Node{}
	err := r.QueryRow(query, domainID, content).Scan(
		&node.ID, &node.Content, &node.DomainID, &node.Title,
		&node.Description, &node.CreatedAt, &node.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrNodeNotFound
	}
	
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	
	return node, nil
}

// ListByDomain 은 도메인별 노드 목록을 페이지네이션과 함께 조회합니다.
func (r *sqliteNodeRepository) ListByDomain(domainID int, offset, limit int) ([]models.Node, int, error) {
	// 데이터 조회
	query := `
		SELECT id, content, domain_id, title, description, created_at, updated_at
		FROM nodes
		WHERE domain_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.Query(query, domainID, limit, offset)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	defer rows.Close()
	
	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(&node.ID, &node.Content, &node.DomainID, &node.Title,
			&node.Description, &node.CreatedAt, &node.UpdatedAt)
		if err != nil {
			return nil, 0, MapSQLiteError(err)
		}
		nodes = append(nodes, node)
	}
	
	if err := rows.Err(); err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	// 총 개수 조회
	countQuery := `SELECT COUNT(*) FROM nodes WHERE domain_id = ?`
	var total int
	err = r.QueryRow(countQuery, domainID).Scan(&total)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	return nodes, total, nil
}

// Search 는 도메인 내에서 노드를 검색합니다.
func (r *sqliteNodeRepository) Search(domainID int, query string, offset, limit int) ([]models.Node, int, error) {
	// 데이터 검색
	searchQuery := `
		SELECT id, content, domain_id, title, description, created_at, updated_at
		FROM nodes
		WHERE domain_id = ? AND (title LIKE ? OR content LIKE ? OR description LIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	searchPattern := "%" + query + "%"
	rows, err := r.Query(searchQuery, domainID, searchPattern, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	defer rows.Close()
	
	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(&node.ID, &node.Content, &node.DomainID, &node.Title,
			&node.Description, &node.CreatedAt, &node.UpdatedAt)
		if err != nil {
			return nil, 0, MapSQLiteError(err)
		}
		nodes = append(nodes, node)
	}
	
	if err := rows.Err(); err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	// 총 개수 조회
	countQuery := `
		SELECT COUNT(*)
		FROM nodes
		WHERE domain_id = ? AND (title LIKE ? OR content LIKE ? OR description LIKE ?)
	`
	
	var total int
	err = r.QueryRow(countQuery, domainID, searchPattern, searchPattern, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	return nodes, total, nil
}

// Update 는 노드 정보를 업데이트합니다.
func (r *sqliteNodeRepository) Update(node *models.Node) error {
	query := `
		UPDATE nodes
		SET content = ?, title = ?, description = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING updated_at
	`
	
	err := r.QueryRow(query, node.Content, node.Title, node.Description, node.ID).Scan(
		&node.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return ErrNodeNotFound
	}
	
	if err != nil {
		return MapSQLiteError(err)
	}
	
	return nil
}

// Delete 는 노드를 삭제합니다.
func (r *sqliteNodeRepository) Delete(id int) error {
	query := `DELETE FROM nodes WHERE id = ?`
	
	result, err := r.Execute(query, id)
	if err != nil {
		return MapSQLiteError(err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}
	
	if rowsAffected == 0 {
		return ErrNodeNotFound
	}
	
	return nil
}

// ExistsByDomainAndContent 는 도메인 ID와 콘텐츠로 노드 존재 여부를 확인합니다.
func (r *sqliteNodeRepository) ExistsByDomainAndContent(domainID int, content string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM nodes WHERE domain_id = ? AND content = ?)`
	
	var exists bool
	err := r.QueryRow(query, domainID, content).Scan(&exists)
	if err != nil {
		return false, MapSQLiteError(err)
	}
	
	return exists, nil
}

// CountNodesByDomain 는 도메인별 노드 수를 반환합니다.
func (r *sqliteNodeRepository) CountNodesByDomain(ctx context.Context, domainID int) (int, error) {
	query := `SELECT COUNT(*) FROM nodes WHERE domain_id = ?`
	
	var count int
	err := r.QueryRow(query, domainID).Scan(&count)
	if err != nil {
		return 0, MapSQLiteError(err)
	}
	
	return count, nil
}

// BatchCreate 는 여러 노드를 배치로 생성합니다.
func (r *sqliteNodeRepository) BatchCreate(nodes []models.Node) error {
	if len(nodes) == 0 {
		return nil
	}
	
	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := r.PrepareStatementInTransaction(tx, `
			INSERT INTO nodes (content, domain_id, title, description, created_at, updated_at)
			VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		
		for _, node := range nodes {
			_, err := stmt.Exec(node.Content, node.DomainID, node.Title, node.Description)
			if err != nil {
				return MapSQLiteError(err)
			}
		}
		
		return nil
	})
}

// BatchUpdate 는 여러 노드를 배치로 업데이트합니다.
func (r *sqliteNodeRepository) BatchUpdate(nodes []models.Node) error {
	if len(nodes) == 0 {
		return nil
	}
	
	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := r.PrepareStatementInTransaction(tx, `
			UPDATE nodes
			SET content = ?, title = ?, description = ?, updated_at = datetime('now')
			WHERE id = ?
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		
		for _, node := range nodes {
			_, err := stmt.Exec(node.Content, node.Title, node.Description, node.ID)
			if err != nil {
				return MapSQLiteError(err)
			}
		}
		
		return nil
	})
}

// BatchDelete 는 여러 노드를 배치로 삭제합니다.
func (r *sqliteNodeRepository) BatchDelete(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	
	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := r.PrepareStatementInTransaction(tx, `DELETE FROM nodes WHERE id = ?`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		
		for _, id := range ids {
			_, err := stmt.Exec(id)
			if err != nil {
				return MapSQLiteError(err)
			}
		}
		
		return nil
	})
}

// 트랜잭션 지원 메서드들

// CreateTx 는 트랜잭션 내에서 노드를 생성합니다.
func (r *sqliteNodeRepository) CreateTx(tx *sql.Tx, node *models.Node) error {
	query := `
		INSERT INTO nodes (content, domain_id, title, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
		RETURNING id, created_at, updated_at
	`
	
	err := r.QueryRowInTransaction(tx, query, node.Content, node.DomainID, node.Title, node.Description).Scan(
		&node.ID, &node.CreatedAt, &node.UpdatedAt,
	)
	
	if err != nil {
		return MapSQLiteError(err)
	}
	
	return nil
}

// UpdateTx 는 트랜잭션 내에서 노드를 업데이트합니다.
func (r *sqliteNodeRepository) UpdateTx(tx *sql.Tx, node *models.Node) error {
	query := `
		UPDATE nodes
		SET content = ?, title = ?, description = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING updated_at
	`
	
	err := r.QueryRowInTransaction(tx, query, node.Content, node.Title, node.Description, node.ID).Scan(
		&node.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return ErrNodeNotFound
	}
	
	if err != nil {
		return MapSQLiteError(err)
	}
	
	return nil
}

// DeleteTx 는 트랜잭션 내에서 노드를 삭제합니다.
func (r *sqliteNodeRepository) DeleteTx(tx *sql.Tx, id int) error {
	query := `DELETE FROM nodes WHERE id = ?`
	
	result, err := r.ExecuteInTransaction(tx, query, id)
	if err != nil {
		return MapSQLiteError(err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}
	
	if rowsAffected == 0 {
		return ErrNodeNotFound
	}
	
	return nil
}