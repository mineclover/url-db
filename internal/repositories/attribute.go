package repositories

import (
	"database/sql"
	"url-db/internal/models"
)

// sqliteAttributeRepository 는 SQLite 기반 속성 리포지토리 구현체입니다.
type sqliteAttributeRepository struct {
	*BaseRepository
}

// NewSQLiteAttributeRepository 는 새로운 SQLite 속성 리포지토리를 생성합니다.
func NewSQLiteAttributeRepository(db *sql.DB) AttributeRepository {
	return &sqliteAttributeRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 는 새로운 속성을 생성합니다.
func (r *sqliteAttributeRepository) Create(attribute *models.Attribute) error {
	query := `
		INSERT INTO attributes (domain_id, name, type, description, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		RETURNING id, created_at
	`

	err := r.QueryRow(query, attribute.DomainID, attribute.Name, attribute.Type, attribute.Description).Scan(
		&attribute.ID, &attribute.CreatedAt,
	)

	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// GetByID 는 ID로 속성을 조회합니다.
func (r *sqliteAttributeRepository) GetByID(id int) (*models.Attribute, error) {
	query := `
		SELECT id, domain_id, name, type, description, created_at
		FROM attributes
		WHERE id = ?
	`

	attribute := &models.Attribute{}
	err := r.QueryRow(query, id).Scan(
		&attribute.ID, &attribute.DomainID, &attribute.Name,
		&attribute.Type, &attribute.Description, &attribute.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAttributeNotFound
	}

	if err != nil {
		return nil, MapSQLiteError(err)
	}

	return attribute, nil
}

// GetByDomainAndName 은 도메인 ID와 이름으로 속성을 조회합니다.
func (r *sqliteAttributeRepository) GetByDomainAndName(domainID int, name string) (*models.Attribute, error) {
	query := `
		SELECT id, domain_id, name, type, description, created_at
		FROM attributes
		WHERE domain_id = ? AND name = ?
	`

	attribute := &models.Attribute{}
	err := r.QueryRow(query, domainID, name).Scan(
		&attribute.ID, &attribute.DomainID, &attribute.Name,
		&attribute.Type, &attribute.Description, &attribute.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAttributeNotFound
	}

	if err != nil {
		return nil, MapSQLiteError(err)
	}

	return attribute, nil
}

// ListByDomain 은 도메인별 속성 목록을 조회합니다.
func (r *sqliteAttributeRepository) ListByDomain(domainID int) ([]models.Attribute, error) {
	query := `
		SELECT id, domain_id, name, type, description, created_at
		FROM attributes
		WHERE domain_id = ?
		ORDER BY name
	`

	rows, err := r.Query(query, domainID)
	if err != nil {
		return nil, MapSQLiteError(err)
	}
	defer rows.Close()

	var attributes []models.Attribute
	for rows.Next() {
		var attribute models.Attribute
		err := rows.Scan(&attribute.ID, &attribute.DomainID, &attribute.Name,
			&attribute.Type, &attribute.Description, &attribute.CreatedAt)
		if err != nil {
			return nil, MapSQLiteError(err)
		}
		attributes = append(attributes, attribute)
	}

	if err := rows.Err(); err != nil {
		return nil, MapSQLiteError(err)
	}

	return attributes, nil
}

// Update 는 속성 정보를 업데이트합니다.
func (r *sqliteAttributeRepository) Update(attribute *models.Attribute) error {
	query := `
		UPDATE attributes
		SET name = ?, type = ?, description = ?
		WHERE id = ?
	`

	result, err := r.Execute(query, attribute.Name, attribute.Type, attribute.Description, attribute.ID)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrAttributeNotFound
	}

	return nil
}

// Delete 는 속성을 삭제합니다.
func (r *sqliteAttributeRepository) Delete(id int) error {
	query := `DELETE FROM attributes WHERE id = ?`

	result, err := r.Execute(query, id)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrAttributeNotFound
	}

	return nil
}

// ExistsByDomainAndName 은 도메인 ID와 이름으로 속성 존재 여부를 확인합니다.
func (r *sqliteAttributeRepository) ExistsByDomainAndName(domainID int, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM attributes WHERE domain_id = ? AND name = ?)`

	var exists bool
	err := r.QueryRow(query, domainID, name).Scan(&exists)
	if err != nil {
		return false, MapSQLiteError(err)
	}

	return exists, nil
}

// 트랜잭션 지원 메서드들

// CreateTx 는 트랜잭션 내에서 속성을 생성합니다.
func (r *sqliteAttributeRepository) CreateTx(tx *sql.Tx, attribute *models.Attribute) error {
	query := `
		INSERT INTO attributes (domain_id, name, type, description, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		RETURNING id, created_at
	`

	err := r.QueryRowInTransaction(tx, query, attribute.DomainID, attribute.Name, attribute.Type, attribute.Description).Scan(
		&attribute.ID, &attribute.CreatedAt,
	)

	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// UpdateTx 는 트랜잭션 내에서 속성을 업데이트합니다.
func (r *sqliteAttributeRepository) UpdateTx(tx *sql.Tx, attribute *models.Attribute) error {
	query := `
		UPDATE attributes
		SET name = ?, type = ?, description = ?
		WHERE id = ?
	`

	result, err := r.ExecuteInTransaction(tx, query, attribute.Name, attribute.Type, attribute.Description, attribute.ID)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrAttributeNotFound
	}

	return nil
}

// DeleteTx 는 트랜잭션 내에서 속성을 삭제합니다.
func (r *sqliteAttributeRepository) DeleteTx(tx *sql.Tx, id int) error {
	query := `DELETE FROM attributes WHERE id = ?`

	result, err := r.ExecuteInTransaction(tx, query, id)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrAttributeNotFound
	}

	return nil
}
