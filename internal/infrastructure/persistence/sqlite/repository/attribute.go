package repository

import (
	"context"
	"database/sql"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/infrastructure/persistence/sqlite/mapper"
)

type attributeRepository struct {
	db *sql.DB
}

// NewAttributeRepository creates a new attribute repository
func NewAttributeRepository(db *sql.DB) repository.AttributeRepository {
	return &attributeRepository{db: db}
}

func (r *attributeRepository) Create(ctx context.Context, attribute *entity.Attribute) error {
	query := `
		INSERT INTO attributes (name, type, description, domain_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		attribute.Name(),
		attribute.Type(),
		attribute.Description(),
		attribute.DomainID(),
		attribute.CreatedAt(),
		attribute.UpdatedAt(),
	)
	
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	// Set the generated ID back to the entity
	attribute.SetID(int(id))
	
	return nil
}

func (r *attributeRepository) GetByID(ctx context.Context, id int) (*entity.Attribute, error) {
	query := `
		SELECT id, name, type, description, domain_id, created_at, updated_at 
		FROM attributes 
		WHERE id = ?
	`
	
	dbModel := &mapper.AttributeDBModel{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dbModel.ID,
		&dbModel.Name,
		&dbModel.Type,
		&dbModel.Description,
		&dbModel.DomainID,
		&dbModel.CreatedAt,
		&dbModel.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return mapper.ToAttributeEntity(dbModel), nil
}

func (r *attributeRepository) GetByName(ctx context.Context, domainID int, name string) (*entity.Attribute, error) {
	query := `
		SELECT id, name, type, description, domain_id, created_at, updated_at 
		FROM attributes 
		WHERE domain_id = ? AND name = ?
	`
	
	dbModel := &mapper.AttributeDBModel{}
	err := r.db.QueryRowContext(ctx, query, domainID, name).Scan(
		&dbModel.ID,
		&dbModel.Name,
		&dbModel.Type,
		&dbModel.Description,
		&dbModel.DomainID,
		&dbModel.CreatedAt,
		&dbModel.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return mapper.ToAttributeEntity(dbModel), nil
}

func (r *attributeRepository) ListByDomainID(ctx context.Context, domainID int) ([]*entity.Attribute, error) {
	query := `
		SELECT id, name, type, description, domain_id, created_at, updated_at 
		FROM attributes 
		WHERE domain_id = ?
		ORDER BY name
	`
	
	rows, err := r.db.QueryContext(ctx, query, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var attributes []*entity.Attribute
	for rows.Next() {
		dbModel := &mapper.AttributeDBModel{}
		err := rows.Scan(
			&dbModel.ID,
			&dbModel.Name,
			&dbModel.Type,
			&dbModel.Description,
			&dbModel.DomainID,
			&dbModel.CreatedAt,
			&dbModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		attributes = append(attributes, mapper.ToAttributeEntity(dbModel))
	}
	
	return attributes, rows.Err()
}

func (r *attributeRepository) Update(ctx context.Context, attribute *entity.Attribute) error {
	query := `
		UPDATE attributes 
		SET name = ?, type = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	
	_, err := r.db.ExecContext(ctx, query,
		attribute.Name(),
		attribute.Type(),
		attribute.Description(),
		attribute.UpdatedAt(),
		attribute.ID(),
	)
	
	return err
}

func (r *attributeRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM attributes WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}