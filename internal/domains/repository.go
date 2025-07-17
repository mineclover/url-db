package domains

import (
	"context"
	"database/sql"

	"github.com/url-db/internal/models"
)

type DomainRepository interface {
	Create(ctx context.Context, domain *models.Domain) error
	GetByID(ctx context.Context, id int) (*models.Domain, error)
	GetByName(ctx context.Context, name string) (*models.Domain, error)
	List(ctx context.Context, page, size int) ([]*models.Domain, int, error)
	Update(ctx context.Context, domain *models.Domain) error
	Delete(ctx context.Context, id int) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type domainRepository struct {
	db *sql.DB
}

func NewDomainRepository(db *sql.DB) DomainRepository {
	return &domainRepository{db: db}
}

func (r *domainRepository) Create(ctx context.Context, domain *models.Domain) error {
	query := `
		INSERT INTO domains (name, description, created_at, updated_at)
		VALUES (?, ?, datetime('now'), datetime('now'))
		RETURNING id, created_at, updated_at
	`
	
	row := r.db.QueryRowContext(ctx, query, domain.Name, domain.Description)
	return row.Scan(&domain.ID, &domain.CreatedAt, &domain.UpdatedAt)
}

func (r *domainRepository) GetByID(ctx context.Context, id int) (*models.Domain, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM domains
		WHERE id = ?
	`
	
	domain := &models.Domain{}
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&domain.ID, &domain.Name, &domain.Description, &domain.CreatedAt, &domain.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (r *domainRepository) GetByName(ctx context.Context, name string) (*models.Domain, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM domains
		WHERE name = ?
	`
	
	domain := &models.Domain{}
	row := r.db.QueryRowContext(ctx, query, name)
	err := row.Scan(&domain.ID, &domain.Name, &domain.Description, &domain.CreatedAt, &domain.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (r *domainRepository) List(ctx context.Context, page, size int) ([]*models.Domain, int, error) {
	offset := (page - 1) * size
	
	// Get total count
	countQuery := `SELECT COUNT(*) FROM domains`
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}
	
	// Get domains with pagination
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM domains
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var domains []*models.Domain
	for rows.Next() {
		domain := &models.Domain{}
		err := rows.Scan(&domain.ID, &domain.Name, &domain.Description, &domain.CreatedAt, &domain.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		domains = append(domains, domain)
	}
	
	return domains, totalCount, nil
}

func (r *domainRepository) Update(ctx context.Context, domain *models.Domain) error {
	query := `
		UPDATE domains 
		SET description = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING updated_at
	`
	
	row := r.db.QueryRowContext(ctx, query, domain.Description, domain.ID)
	return row.Scan(&domain.UpdatedAt)
}

func (r *domainRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM domains WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *domainRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM domains WHERE name = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	return exists, err
}