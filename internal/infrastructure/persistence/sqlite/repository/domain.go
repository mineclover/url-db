package repository

import (
	"context"
	"database/sql"
	"errors"
	"url-db/internal/constants"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/infrastructure/persistence/sqlite/mapper"
)

type domainRepository struct {
	db *sql.DB
}

// NewDomainRepository creates a new SQLite-based domain repository
func NewDomainRepository(db *sql.DB) repository.DomainRepository {
	return &domainRepository{db: db}
}

func (r *domainRepository) Create(ctx context.Context, domain *entity.Domain) error {
	dbModel := mapper.FromDomainEntity(domain)

	query := `INSERT INTO domains (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		dbModel.Name,
		dbModel.Description,
		dbModel.CreatedAt,
		dbModel.UpdatedAt,
	)

	return err
}

func (r *domainRepository) GetByID(ctx context.Context, id int) (*entity.Domain, error) {
	var dbRow mapper.DatabaseDomain

	query := `SELECT id, name, description, created_at, updated_at FROM domains WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dbRow.ID,
		&dbRow.Name,
		&dbRow.Description,
		&dbRow.CreatedAt,
		&dbRow.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainEntity(&dbRow), nil
}

func (r *domainRepository) GetByName(ctx context.Context, name string) (*entity.Domain, error) {
	var dbRow mapper.DatabaseDomain

	query := `SELECT id, name, description, created_at, updated_at FROM domains WHERE name = ?`
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&dbRow.ID,
		&dbRow.Name,
		&dbRow.Description,
		&dbRow.CreatedAt,
		&dbRow.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainEntity(&dbRow), nil
}

func (r *domainRepository) List(ctx context.Context, page, size int) ([]*entity.Domain, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM domains`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * size

	// Get domains with pagination
	query := `SELECT id, name, description, created_at, updated_at FROM domains ORDER BY name LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var domains []*entity.Domain
	for rows.Next() {
		var dbRow mapper.DatabaseDomain
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.Description,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		domain := mapper.ToDomainEntity(&dbRow)
		if domain != nil {
			domains = append(domains, domain)
		}
	}

	return domains, totalCount, nil
}

func (r *domainRepository) Update(ctx context.Context, domain *entity.Domain) error {
	dbModel := mapper.FromDomainEntity(domain)

	query := `UPDATE domains SET description = ?, updated_at = ? WHERE name = ?`
	result, err := r.db.ExecContext(ctx, query,
		dbModel.Description,
		dbModel.UpdatedAt,
		dbModel.Name,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(constants.ErrDomainNotFound)
	}

	return nil
}

func (r *domainRepository) Delete(ctx context.Context, name string) error {
	query := `DELETE FROM domains WHERE name = ?`
	result, err := r.db.ExecContext(ctx, query, name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(constants.ErrDomainNotFound)
	}

	return nil
}

func (r *domainRepository) Exists(ctx context.Context, name string) (bool, error) {
	var exists int
	query := `SELECT 1 FROM domains WHERE name = ? LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
