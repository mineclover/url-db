package repositories

import (
	"database/sql"
	"url-db/internal/models"
)

// sqliteNodeAttributeRepository 는 SQLite 기반 노드 속성 리포지토리 구현체입니다.
type sqliteNodeAttributeRepository struct {
	*BaseRepository
}

// NewSQLiteNodeAttributeRepository 는 새로운 SQLite 노드 속성 리포지토리를 생성합니다.
func NewSQLiteNodeAttributeRepository(db *sql.DB) NodeAttributeRepository {
	return &sqliteNodeAttributeRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 는 새로운 노드 속성을 생성합니다.
func (r *sqliteNodeAttributeRepository) Create(nodeAttribute *models.NodeAttribute) error {
	query := `
		INSERT INTO node_attributes (node_id, attribute_id, value, order_index, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		RETURNING id, created_at
	`

	err := r.QueryRow(query, nodeAttribute.NodeID, nodeAttribute.AttributeID,
		nodeAttribute.Value, nodeAttribute.OrderIndex).Scan(
		&nodeAttribute.ID, &nodeAttribute.CreatedAt,
	)

	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// GetByID 는 ID로 노드 속성을 조회합니다.
func (r *sqliteNodeAttributeRepository) GetByID(id int) (*models.NodeAttribute, error) {
	query := `
		SELECT id, node_id, attribute_id, value, order_index, created_at
		FROM node_attributes
		WHERE id = ?
	`

	nodeAttribute := &models.NodeAttribute{}
	err := r.QueryRow(query, id).Scan(
		&nodeAttribute.ID, &nodeAttribute.NodeID, &nodeAttribute.AttributeID,
		&nodeAttribute.Value, &nodeAttribute.OrderIndex, &nodeAttribute.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNodeAttributeNotFound
	}

	if err != nil {
		return nil, MapSQLiteError(err)
	}

	return nodeAttribute, nil
}

// GetByNodeAndAttribute 는 노드 ID와 속성 ID로 노드 속성을 조회합니다.
func (r *sqliteNodeAttributeRepository) GetByNodeAndAttribute(nodeID, attributeID int) (*models.NodeAttribute, error) {
	query := `
		SELECT id, node_id, attribute_id, value, order_index, created_at
		FROM node_attributes
		WHERE node_id = ? AND attribute_id = ?
	`

	nodeAttribute := &models.NodeAttribute{}
	err := r.QueryRow(query, nodeID, attributeID).Scan(
		&nodeAttribute.ID, &nodeAttribute.NodeID, &nodeAttribute.AttributeID,
		&nodeAttribute.Value, &nodeAttribute.OrderIndex, &nodeAttribute.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNodeAttributeNotFound
	}

	if err != nil {
		return nil, MapSQLiteError(err)
	}

	return nodeAttribute, nil
}

// ListByNode 은 노드 ID로 노드 속성 목록을 조회합니다 (속성 정보 포함).
func (r *sqliteNodeAttributeRepository) ListByNode(nodeID int) ([]models.NodeAttributeWithInfo, error) {
	query := `
		SELECT na.id, na.node_id, na.attribute_id, na.value, na.order_index, na.created_at,
		       a.name, a.type
		FROM node_attributes na
		JOIN attributes a ON na.attribute_id = a.id
		WHERE na.node_id = ?
		ORDER BY a.name, na.order_index
	`

	rows, err := r.Query(query, nodeID)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	defer rows.Close()

	var attributes []models.NodeAttributeWithInfo
	for rows.Next() {
		var attr models.NodeAttributeWithInfo
		err := rows.Scan(&attr.ID, &attr.NodeID, &attr.AttributeID,
			&attr.Value, &attr.OrderIndex, &attr.CreatedAt,
			&attr.Name, &attr.Type)
		if err != nil {
			return nil, MapSQLiteError(err)
		}
		attributes = append(attributes, attr)
	}

	if err := rows.Err(); err != nil {
		return nil, MapSQLiteError(err)
	}

	return attributes, nil
}

// ListByAttribute 는 속성 ID로 노드 속성 목록을 조회합니다.
func (r *sqliteNodeAttributeRepository) ListByAttribute(attributeID int) ([]models.NodeAttribute, error) {
	query := `
		SELECT id, node_id, attribute_id, value, order_index, created_at
		FROM node_attributes
		WHERE attribute_id = ?
		ORDER BY order_index
	`

	rows, err := r.Query(query, attributeID)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	defer rows.Close()

	var nodeAttributes []models.NodeAttribute
	for rows.Next() {
		var nodeAttribute models.NodeAttribute
		err := rows.Scan(&nodeAttribute.ID, &nodeAttribute.NodeID, &nodeAttribute.AttributeID,
			&nodeAttribute.Value, &nodeAttribute.OrderIndex, &nodeAttribute.CreatedAt)
		if err != nil {
			return nil, MapSQLiteError(err)
		}
		nodeAttributes = append(nodeAttributes, nodeAttribute)
	}

	if err := rows.Err(); err != nil {
		return nil, MapSQLiteError(err)
	}

	return nodeAttributes, nil
}

// Update 는 노드 속성 정보를 업데이트합니다.
func (r *sqliteNodeAttributeRepository) Update(nodeAttribute *models.NodeAttribute) error {
	query := `
		UPDATE node_attributes
		SET value = ?, order_index = ?
		WHERE id = ?
	`

	result, err := r.Execute(query, nodeAttribute.Value, nodeAttribute.OrderIndex, nodeAttribute.ID)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrNodeAttributeNotFound
	}

	return nil
}

// Delete 는 노드 속성을 삭제합니다.
func (r *sqliteNodeAttributeRepository) Delete(id int) error {
	query := `DELETE FROM node_attributes WHERE id = ?`

	result, err := r.Execute(query, id)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrNodeAttributeNotFound
	}

	return nil
}

// DeleteByNode 는 노드 ID로 모든 노드 속성을 삭제합니다.
func (r *sqliteNodeAttributeRepository) DeleteByNode(nodeID int) error {
	query := `DELETE FROM node_attributes WHERE node_id = ?`

	_, err := r.Execute(query, nodeID)
	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// DeleteByAttribute 는 속성 ID로 모든 노드 속성을 삭제합니다.
func (r *sqliteNodeAttributeRepository) DeleteByAttribute(attributeID int) error {
	query := `DELETE FROM node_attributes WHERE attribute_id = ?`

	_, err := r.Execute(query, attributeID)
	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// ExistsByNodeAndAttribute 는 노드 ID와 속성 ID로 노드 속성 존재 여부를 확인합니다.
func (r *sqliteNodeAttributeRepository) ExistsByNodeAndAttribute(nodeID, attributeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM node_attributes WHERE node_id = ? AND attribute_id = ?)`

	var exists bool
	err := r.QueryRow(query, nodeID, attributeID).Scan(&exists)
	if err != nil {
		return false, MapSQLiteError(err)
	}

	return exists, nil
}

// BatchCreate 는 여러 노드 속성을 배치로 생성합니다.
func (r *sqliteNodeAttributeRepository) BatchCreate(nodeAttributes []models.NodeAttribute) error {
	if len(nodeAttributes) == 0 {
		return nil
	}

	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := r.PrepareStatementInTransaction(tx, `
			INSERT INTO node_attributes (node_id, attribute_id, value, order_index, created_at)
			VALUES (?, ?, ?, ?, datetime('now'))
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, nodeAttribute := range nodeAttributes {
			_, err := stmt.Exec(nodeAttribute.NodeID, nodeAttribute.AttributeID,
				nodeAttribute.Value, nodeAttribute.OrderIndex)
			if err != nil {
				return MapSQLiteError(err)
			}
		}

		return nil
	})
}

// BatchUpdate 는 여러 노드 속성을 배치로 업데이트합니다.
func (r *sqliteNodeAttributeRepository) BatchUpdate(nodeAttributes []models.NodeAttribute) error {
	if len(nodeAttributes) == 0 {
		return nil
	}

	return r.WithTransaction(func(tx *sql.Tx) error {
		stmt, err := r.PrepareStatementInTransaction(tx, `
			UPDATE node_attributes
			SET value = ?, order_index = ?
			WHERE id = ?
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, nodeAttribute := range nodeAttributes {
			_, err := stmt.Exec(nodeAttribute.Value, nodeAttribute.OrderIndex, nodeAttribute.ID)
			if err != nil {
				return MapSQLiteError(err)
			}
		}

		return nil
	})
}

// BatchDeleteByNode 는 노드 ID로 모든 노드 속성을 삭제합니다.
func (r *sqliteNodeAttributeRepository) BatchDeleteByNode(nodeID int) error {
	query := `DELETE FROM node_attributes WHERE node_id = ?`

	_, err := r.Execute(query, nodeID)
	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// BatchDeleteByAttribute 는 속성 ID로 모든 노드 속성을 삭제합니다.
func (r *sqliteNodeAttributeRepository) BatchDeleteByAttribute(attributeID int) error {
	query := `DELETE FROM node_attributes WHERE attribute_id = ?`

	_, err := r.Execute(query, attributeID)
	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// 트랜잭션 지원 메서드들

// CreateTx 는 트랜잭션 내에서 노드 속성을 생성합니다.
func (r *sqliteNodeAttributeRepository) CreateTx(tx *sql.Tx, nodeAttribute *models.NodeAttribute) error {
	query := `
		INSERT INTO node_attributes (node_id, attribute_id, value, order_index, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		RETURNING id, created_at
	`

	err := r.QueryRowInTransaction(tx, query, nodeAttribute.NodeID, nodeAttribute.AttributeID,
		nodeAttribute.Value, nodeAttribute.OrderIndex).Scan(
		&nodeAttribute.ID, &nodeAttribute.CreatedAt,
	)

	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// UpdateTx 는 트랜잭션 내에서 노드 속성을 업데이트합니다.
func (r *sqliteNodeAttributeRepository) UpdateTx(tx *sql.Tx, nodeAttribute *models.NodeAttribute) error {
	query := `
		UPDATE node_attributes
		SET value = ?, order_index = ?
		WHERE id = ?
	`

	result, err := r.ExecuteInTransaction(tx, query, nodeAttribute.Value, nodeAttribute.OrderIndex, nodeAttribute.ID)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrNodeAttributeNotFound
	}

	return nil
}

// DeleteTx 는 트랜잭션 내에서 노드 속성을 삭제합니다.
func (r *sqliteNodeAttributeRepository) DeleteTx(tx *sql.Tx, id int) error {
	query := `DELETE FROM node_attributes WHERE id = ?`

	result, err := r.ExecuteInTransaction(tx, query, id)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrNodeAttributeNotFound
	}

	return nil
}
