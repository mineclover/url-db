package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"url-db/internal/models"
)

// sqliteNodeConnectionRepository 는 SQLite 기반 노드 연결 리포지토리 구현체입니다.
type sqliteNodeConnectionRepository struct {
	*BaseRepository
}

// NewSQLiteNodeConnectionRepository 는 새로운 SQLite 노드 연결 리포지토리를 생성합니다.
func NewSQLiteNodeConnectionRepository(db *sql.DB) NodeConnectionRepository {
	return &sqliteNodeConnectionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 는 새로운 노드 연결을 생성합니다.
func (r *sqliteNodeConnectionRepository) Create(ctx context.Context, connection *models.NodeConnection) error {
	query := `
		INSERT INTO node_connections (source_node_id, target_node_id, relationship_type, description, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		RETURNING id, created_at
	`
	
	err := r.QueryRow(query, connection.SourceNodeID, connection.TargetNodeID, connection.RelationshipType, connection.Description).Scan(
		&connection.ID, &connection.CreatedAt,
	)
	
	if err != nil {
		return MapSQLiteError(err)
	}
	
	return nil
}

// GetByID 는 ID로 노드 연결을 조회합니다.
func (r *sqliteNodeConnectionRepository) GetByID(ctx context.Context, id int) (*models.NodeConnection, error) {
	query := `
		SELECT id, source_node_id, target_node_id, relationship_type, description, created_at
		FROM node_connections
		WHERE id = ?
	`
	
	connection := &models.NodeConnection{}
	err := r.QueryRow(query, id).Scan(
		&connection.ID, &connection.SourceNodeID, &connection.TargetNodeID,
		&connection.RelationshipType, &connection.Description, &connection.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNodeConnectionNotFound
		}
		return nil, MapSQLiteError(err)
	}
	
	return connection, nil
}

// GetBySourceAndTarget 는 소스 노드와 타겟 노드 및 관계 유형으로 연결을 조회합니다.
func (r *sqliteNodeConnectionRepository) GetBySourceAndTarget(ctx context.Context, sourceNodeID, targetNodeID int, relationshipType string) (*models.NodeConnection, error) {
	query := `
		SELECT id, source_node_id, target_node_id, relationship_type, description, created_at
		FROM node_connections
		WHERE source_node_id = ? AND target_node_id = ? AND relationship_type = ?
	`
	
	connection := &models.NodeConnection{}
	err := r.QueryRow(query, sourceNodeID, targetNodeID, relationshipType).Scan(
		&connection.ID, &connection.SourceNodeID, &connection.TargetNodeID,
		&connection.RelationshipType, &connection.Description, &connection.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNodeConnectionNotFound
		}
		return nil, MapSQLiteError(err)
	}
	
	return connection, nil
}

// ListBySourceNode 는 소스 노드별로 연결을 조회합니다.
func (r *sqliteNodeConnectionRepository) ListBySourceNode(ctx context.Context, sourceNodeID int, offset, limit int) ([]models.NodeConnectionWithInfo, int, error) {
	query := `
		SELECT nc.id, nc.source_node_id, nc.target_node_id, nc.relationship_type, nc.description, nc.created_at,
			   sn.content as source_node_url, tn.content as target_node_url,
			   sn.title as source_node_title, tn.title as target_node_title
		FROM node_connections nc
		JOIN nodes sn ON nc.source_node_id = sn.id
		JOIN nodes tn ON nc.target_node_id = tn.id
		WHERE nc.source_node_id = ?
		ORDER BY nc.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.Query(query, sourceNodeID, limit, offset)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	defer rows.Close()
	
	var connections []models.NodeConnectionWithInfo
	for rows.Next() {
		var conn models.NodeConnectionWithInfo
		err := rows.Scan(
			&conn.ID, &conn.SourceNodeID, &conn.TargetNodeID, &conn.RelationshipType,
			&conn.Description, &conn.CreatedAt, &conn.SourceNodeURL, &conn.TargetNodeURL,
			&conn.SourceNodeTitle, &conn.TargetNodeTitle,
		)
		if err != nil {
			return nil, 0, MapSQLiteError(err)
		}
		connections = append(connections, conn)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	// 총 개수 조회
	countQuery := `SELECT COUNT(*) FROM node_connections WHERE source_node_id = ?`
	var totalCount int
	err = r.QueryRow(countQuery, sourceNodeID).Scan(&totalCount)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	return connections, totalCount, nil
}

// ListByTargetNode 는 타겟 노드별로 연결을 조회합니다.
func (r *sqliteNodeConnectionRepository) ListByTargetNode(ctx context.Context, targetNodeID int, offset, limit int) ([]models.NodeConnectionWithInfo, int, error) {
	query := `
		SELECT nc.id, nc.source_node_id, nc.target_node_id, nc.relationship_type, nc.description, nc.created_at,
			   sn.content as source_node_url, tn.content as target_node_url,
			   sn.title as source_node_title, tn.title as target_node_title
		FROM node_connections nc
		JOIN nodes sn ON nc.source_node_id = sn.id
		JOIN nodes tn ON nc.target_node_id = tn.id
		WHERE nc.target_node_id = ?
		ORDER BY nc.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.Query(query, targetNodeID, limit, offset)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	defer rows.Close()
	
	var connections []models.NodeConnectionWithInfo
	for rows.Next() {
		var conn models.NodeConnectionWithInfo
		err := rows.Scan(
			&conn.ID, &conn.SourceNodeID, &conn.TargetNodeID, &conn.RelationshipType,
			&conn.Description, &conn.CreatedAt, &conn.SourceNodeURL, &conn.TargetNodeURL,
			&conn.SourceNodeTitle, &conn.TargetNodeTitle,
		)
		if err != nil {
			return nil, 0, MapSQLiteError(err)
		}
		connections = append(connections, conn)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	// 총 개수 조회
	countQuery := `SELECT COUNT(*) FROM node_connections WHERE target_node_id = ?`
	var totalCount int
	err = r.QueryRow(countQuery, targetNodeID).Scan(&totalCount)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	return connections, totalCount, nil
}

// ListByRelationshipType 는 관계 유형별로 연결을 조회합니다.
func (r *sqliteNodeConnectionRepository) ListByRelationshipType(ctx context.Context, relationshipType string, offset, limit int) ([]models.NodeConnectionWithInfo, int, error) {
	query := `
		SELECT nc.id, nc.source_node_id, nc.target_node_id, nc.relationship_type, nc.description, nc.created_at,
			   sn.content as source_node_url, tn.content as target_node_url,
			   sn.title as source_node_title, tn.title as target_node_title
		FROM node_connections nc
		JOIN nodes sn ON nc.source_node_id = sn.id
		JOIN nodes tn ON nc.target_node_id = tn.id
		WHERE nc.relationship_type = ?
		ORDER BY nc.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.Query(query, relationshipType, limit, offset)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	defer rows.Close()
	
	var connections []models.NodeConnectionWithInfo
	for rows.Next() {
		var conn models.NodeConnectionWithInfo
		err := rows.Scan(
			&conn.ID, &conn.SourceNodeID, &conn.TargetNodeID, &conn.RelationshipType,
			&conn.Description, &conn.CreatedAt, &conn.SourceNodeURL, &conn.TargetNodeURL,
			&conn.SourceNodeTitle, &conn.TargetNodeTitle,
		)
		if err != nil {
			return nil, 0, MapSQLiteError(err)
		}
		connections = append(connections, conn)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	// 총 개수 조회
	countQuery := `SELECT COUNT(*) FROM node_connections WHERE relationship_type = ?`
	var totalCount int
	err = r.QueryRow(countQuery, relationshipType).Scan(&totalCount)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	
	return connections, totalCount, nil
}

// Update 는 노드 연결을 업데이트합니다.
func (r *sqliteNodeConnectionRepository) Update(ctx context.Context, connection *models.NodeConnection) error {
	query := `
		UPDATE node_connections
		SET relationship_type = ?, description = ?
		WHERE id = ?
	`
	
	result, err := r.Execute(query, connection.RelationshipType, connection.Description, connection.ID)
	if err != nil {
		return MapSQLiteError(err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}
	
	if rowsAffected == 0 {
		return ErrNodeConnectionNotFound
	}
	
	return nil
}

// Delete 는 노드 연결을 삭제합니다.
func (r *sqliteNodeConnectionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM node_connections WHERE id = ?`
	
	result, err := r.Execute(query, id)
	if err != nil {
		return MapSQLiteError(err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}
	
	if rowsAffected == 0 {
		return ErrNodeConnectionNotFound
	}
	
	return nil
}

// DeleteBySourceNode 는 소스 노드로 연결을 삭제합니다.
func (r *sqliteNodeConnectionRepository) DeleteBySourceNode(ctx context.Context, sourceNodeID int) error {
	query := `DELETE FROM node_connections WHERE source_node_id = ?`
	
	_, err := r.Execute(query, sourceNodeID)
	if err != nil {
		return MapSQLiteError(err)
	}
	
	return nil
}

// DeleteByTargetNode 는 타겟 노드로 연결을 삭제합니다.
func (r *sqliteNodeConnectionRepository) DeleteByTargetNode(ctx context.Context, targetNodeID int) error {
	query := `DELETE FROM node_connections WHERE target_node_id = ?`
	
	_, err := r.Execute(query, targetNodeID)
	if err != nil {
		return MapSQLiteError(err)
	}
	
	return nil
}

// ExistsBySourceAndTarget 는 소스 노드와 타겟 노드 및 관계 유형으로 연결 존재 여부를 확인합니다.
func (r *sqliteNodeConnectionRepository) ExistsBySourceAndTarget(ctx context.Context, sourceNodeID, targetNodeID int, relationshipType string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM node_connections WHERE source_node_id = ? AND target_node_id = ? AND relationship_type = ?)`
	
	var exists bool
	err := r.QueryRow(query, sourceNodeID, targetNodeID, relationshipType).Scan(&exists)
	if err != nil {
		return false, MapSQLiteError(err)
	}
	
	return exists, nil
}

// BatchCreate 는 여러 노드 연결을 배치로 생성합니다.
func (r *sqliteNodeConnectionRepository) BatchCreate(ctx context.Context, connections []models.NodeConnection) error {
	if len(connections) == 0 {
		return nil
	}
	
	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := r.PrepareStatementInTransaction(tx, `
			INSERT INTO node_connections (source_node_id, target_node_id, relationship_type, description, created_at)
			VALUES (?, ?, ?, ?, datetime('now'))
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		
		for _, connection := range connections {
			_, err = stmt.Exec(connection.SourceNodeID, connection.TargetNodeID, connection.RelationshipType, connection.Description)
			if err != nil {
				return MapSQLiteError(err)
			}
		}
		
		return nil
	})
}

// BatchDelete 는 여러 노드 연결을 배치로 삭제합니다.
func (r *sqliteNodeConnectionRepository) BatchDelete(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	
	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := r.PrepareStatementInTransaction(tx, `DELETE FROM node_connections WHERE id = ?`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		
		for _, id := range ids {
			_, err = stmt.Exec(id)
			if err != nil {
				return MapSQLiteError(err)
			}
		}
		
		return nil
	})
}

// CreateTx 는 트랜잭션 내에서 노드 연결을 생성합니다.
func (r *sqliteNodeConnectionRepository) CreateTx(tx *sql.Tx, connection *models.NodeConnection) error {
	query := `
		INSERT INTO node_connections (source_node_id, target_node_id, relationship_type, description, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		RETURNING id, created_at
	`
	
	err := tx.QueryRow(query, connection.SourceNodeID, connection.TargetNodeID, connection.RelationshipType, connection.Description).Scan(
		&connection.ID, &connection.CreatedAt,
	)
	
	if err != nil {
		return MapSQLiteError(err)
	}
	
	return nil
}

// UpdateTx 는 트랜잭션 내에서 노드 연결을 업데이트합니다.
func (r *sqliteNodeConnectionRepository) UpdateTx(tx *sql.Tx, connection *models.NodeConnection) error {
	query := `
		UPDATE node_connections
		SET relationship_type = ?, description = ?
		WHERE id = ?
	`
	
	result, err := tx.Exec(query, connection.RelationshipType, connection.Description, connection.ID)
	if err != nil {
		return MapSQLiteError(err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}
	
	if rowsAffected == 0 {
		return ErrNodeConnectionNotFound
	}
	
	return nil
}

// DeleteTx 는 트랜잭션 내에서 노드 연결을 삭제합니다.
func (r *sqliteNodeConnectionRepository) DeleteTx(tx *sql.Tx, id int) error {
	query := `DELETE FROM node_connections WHERE id = ?`
	
	result, err := tx.Exec(query, id)
	if err != nil {
		return MapSQLiteError(err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}
	
	if rowsAffected == 0 {
		return ErrNodeConnectionNotFound
	}
	
	return nil
}