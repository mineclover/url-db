package usecases

import (
	"context"
	"database/sql"
	"fmt"

	"url-db/internal/models"
)

// AttributeRepository defines the interface for attribute data access
type AttributeRepository interface {
	Create(ctx context.Context, attribute *models.Attribute) error
	GetByID(ctx context.Context, id int) (*models.Attribute, error)
	GetByDomainID(ctx context.Context, domainID int) ([]*models.Attribute, error)
	GetByDomainIDAndName(ctx context.Context, domainID int, name string) (*models.Attribute, error)
	Update(ctx context.Context, attribute *models.Attribute) error
	Delete(ctx context.Context, id int) error
	HasValues(ctx context.Context, attributeID int) (bool, error)
}

// SQLiteAttributeRepository implements AttributeRepository for SQLite
type SQLiteAttributeRepository struct {
	db *sql.DB
}

// NewSQLiteAttributeRepository creates a new SQLite attribute repository
func NewSQLiteAttributeRepository(db *sql.DB) *SQLiteAttributeRepository {
	return &SQLiteAttributeRepository{db: db}
}

// Create creates a new attribute
func (r *SQLiteAttributeRepository) Create(ctx context.Context, attribute *models.Attribute) error {
	query := `
		INSERT INTO attributes (domain_id, name, type, description, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		attribute.DomainID,
		attribute.Name,
		attribute.Type,
		attribute.Description,
		attribute.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create attribute: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	attribute.ID = int(id)
	return nil
}

// GetByID retrieves an attribute by ID
func (r *SQLiteAttributeRepository) GetByID(ctx context.Context, id int) (*models.Attribute, error) {
	query := `
		SELECT id, domain_id, name, type, description, created_at
		FROM attributes
		WHERE id = ?
	`

	attribute := &models.Attribute{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attribute.ID,
		&attribute.DomainID,
		&attribute.Name,
		&attribute.Type,
		&attribute.Description,
		&attribute.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAttributeNotFound
		}
		return nil, fmt.Errorf("failed to get attribute by id: %w", err)
	}

	return attribute, nil
}

// GetByDomainID retrieves all attributes for a domain
func (r *SQLiteAttributeRepository) GetByDomainID(ctx context.Context, domainID int) ([]*models.Attribute, error) {
	query := `
		SELECT id, domain_id, name, type, description, created_at
		FROM attributes
		WHERE domain_id = ?
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query, domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attributes by domain id: %w", err)
	}
	defer rows.Close()

	var attributes []*models.Attribute
	for rows.Next() {
		attribute := &models.Attribute{}
		err := rows.Scan(
			&attribute.ID,
			&attribute.DomainID,
			&attribute.Name,
			&attribute.Type,
			&attribute.Description,
			&attribute.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}
		attributes = append(attributes, attribute)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over attributes: %w", err)
	}

	return attributes, nil
}

// GetByDomainIDAndName retrieves an attribute by domain ID and name
func (r *SQLiteAttributeRepository) GetByDomainIDAndName(ctx context.Context, domainID int, name string) (*models.Attribute, error) {
	query := `
		SELECT id, domain_id, name, type, description, created_at
		FROM attributes
		WHERE domain_id = ? AND name = ?
	`

	attribute := &models.Attribute{}
	err := r.db.QueryRowContext(ctx, query, domainID, name).Scan(
		&attribute.ID,
		&attribute.DomainID,
		&attribute.Name,
		&attribute.Type,
		&attribute.Description,
		&attribute.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAttributeNotFound
		}
		return nil, fmt.Errorf("failed to get attribute by domain id and name: %w", err)
	}

	return attribute, nil
}

// Update updates an attribute
func (r *SQLiteAttributeRepository) Update(ctx context.Context, attribute *models.Attribute) error {
	query := `
		UPDATE attributes
		SET description = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, attribute.Description, attribute.ID)
	if err != nil {
		return fmt.Errorf("failed to update attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAttributeNotFound
	}

	return nil
}

// Delete deletes an attribute
func (r *SQLiteAttributeRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM attributes WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAttributeNotFound
	}

	return nil
}

// HasValues checks if an attribute has any values
func (r *SQLiteAttributeRepository) HasValues(ctx context.Context, attributeID int) (bool, error) {
	query := `SELECT COUNT(*) FROM node_attributes WHERE attribute_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, attributeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if attribute has values: %w", err)
	}

	return count > 0, nil
}
