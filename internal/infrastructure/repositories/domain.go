package repositories

import (
	"database/sql"
	"url-db/internal/models"
)

// sqliteDomainRepository 는 SQLite 기반 도메인 리포지토리 구현체입니다.
type sqliteDomainRepository struct {
	*BaseRepository
}

// NewSQLiteDomainRepository 는 새로운 SQLite 도메인 리포지토리를 생성합니다.
func NewSQLiteDomainRepository(db *sql.DB) DomainRepository {
	return &sqliteDomainRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 는 새로운 도메인을 생성합니다.
func (r *sqliteDomainRepository) Create(domain *models.Domain) error {
	query := `
		INSERT INTO domains (name, description, created_at, updated_at)
		VALUES (?, ?, datetime('now'), datetime('now'))
		RETURNING id, created_at, updated_at
	`

	err := r.QueryRow(query, domain.Name, domain.Description).Scan(
		&domain.ID, &domain.CreatedAt, &domain.UpdatedAt,
	)

	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// GetByID 는 ID로 도메인을 조회합니다.
func (r *sqliteDomainRepository) GetByID(id int) (*models.Domain, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM domains
		WHERE id = ?
	`

	domain := &models.Domain{}
	err := r.QueryRow(query, id).Scan(
		&domain.ID, &domain.Name, &domain.Description,
		&domain.CreatedAt, &domain.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrDomainNotFound
	}

	if err != nil {
		return nil, MapSQLiteError(err)
	}

	return domain, nil
}

// GetByName 은 이름으로 도메인을 조회합니다.
func (r *sqliteDomainRepository) GetByName(name string) (*models.Domain, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM domains
		WHERE name = ?
	`

	domain := &models.Domain{}
	err := r.QueryRow(query, name).Scan(
		&domain.ID, &domain.Name, &domain.Description,
		&domain.CreatedAt, &domain.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrDomainNotFound
	}

	if err != nil {
		return nil, MapSQLiteError(err)
	}

	return domain, nil
}

// List 는 도메인 목록을 페이지네이션과 함께 조회합니다.
func (r *sqliteDomainRepository) List(offset, limit int) ([]models.Domain, int, error) {
	// 데이터 조회
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM domains
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.Query(query, limit, offset)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var domain models.Domain
		err := rows.Scan(&domain.ID, &domain.Name, &domain.Description,
			&domain.CreatedAt, &domain.UpdatedAt)
		if err != nil {
			return nil, 0, MapSQLiteError(err)
		}
		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, MapSQLiteError(err)
	}

	// 총 개수 조회
	countQuery := `SELECT COUNT(*) FROM domains`
	var total int
	err = r.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, MapSQLiteError(err)
	}

	return domains, total, nil
}

// Update 는 도메인 정보를 업데이트합니다.
func (r *sqliteDomainRepository) Update(domain *models.Domain) error {
	query := `
		UPDATE domains
		SET name = ?, description = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING updated_at
	`

	err := r.QueryRow(query, domain.Name, domain.Description, domain.ID).Scan(
		&domain.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return ErrDomainNotFound
	}

	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// Delete 는 도메인을 삭제합니다.
func (r *sqliteDomainRepository) Delete(id int) error {
	query := `DELETE FROM domains WHERE id = ?`

	result, err := r.Execute(query, id)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrDomainNotFound
	}

	return nil
}

// ExistsByName 은 이름으로 도메인 존재 여부를 확인합니다.
func (r *sqliteDomainRepository) ExistsByName(name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM domains WHERE name = ?)`

	var exists bool
	err := r.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, MapSQLiteError(err)
	}

	return exists, nil
}

// 트랜잭션 지원 메서드들

// CreateTx 는 트랜잭션 내에서 도메인을 생성합니다.
func (r *sqliteDomainRepository) CreateTx(tx *sql.Tx, domain *models.Domain) error {
	query := `
		INSERT INTO domains (name, description, created_at, updated_at)
		VALUES (?, ?, datetime('now'), datetime('now'))
		RETURNING id, created_at, updated_at
	`

	err := r.QueryRowInTransaction(tx, query, domain.Name, domain.Description).Scan(
		&domain.ID, &domain.CreatedAt, &domain.UpdatedAt,
	)

	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// UpdateTx 는 트랜잭션 내에서 도메인을 업데이트합니다.
func (r *sqliteDomainRepository) UpdateTx(tx *sql.Tx, domain *models.Domain) error {
	query := `
		UPDATE domains
		SET name = ?, description = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING updated_at
	`

	err := r.QueryRowInTransaction(tx, query, domain.Name, domain.Description, domain.ID).Scan(
		&domain.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return ErrDomainNotFound
	}

	if err != nil {
		return MapSQLiteError(err)
	}

	return nil
}

// DeleteTx 는 트랜잭션 내에서 도메인을 삭제합니다.
func (r *sqliteDomainRepository) DeleteTx(tx *sql.Tx, id int) error {
	query := `DELETE FROM domains WHERE id = ?`

	result, err := r.ExecuteInTransaction(tx, query, id)
	if err != nil {
		return MapSQLiteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapSQLiteError(err)
	}

	if rowsAffected == 0 {
		return ErrDomainNotFound
	}

	return nil
}
