package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/infrastructure/persistence/sqlite/mapper"
)

type templateRepository struct {
	db *sql.DB
}

// NewTemplateRepository creates a new SQLite-based template repository
func NewTemplateRepository(db *sql.DB) repository.TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(ctx context.Context, template *entity.Template) error {
	dbModel := mapper.FromTemplateEntity(template)

	query := `INSERT INTO templates (name, domain_id, template_data, title, description, is_active, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		dbModel.Name,
		dbModel.DomainID,
		dbModel.TemplateData,
		dbModel.Title,
		dbModel.Description,
		dbModel.IsActive,
		dbModel.CreatedAt,
		dbModel.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	template.SetID(int(id))
	return nil
}

func (r *templateRepository) GetByID(ctx context.Context, id int) (*entity.Template, error) {
	var dbRow mapper.DatabaseTemplate

	query := `SELECT id, name, domain_id, template_data, title, description, is_active, created_at, updated_at 
			  FROM templates WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dbRow.ID,
		&dbRow.Name,
		&dbRow.DomainID,
		&dbRow.TemplateData,
		&dbRow.Title,
		&dbRow.Description,
		&dbRow.IsActive,
		&dbRow.CreatedAt,
		&dbRow.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return mapper.ToTemplateEntity(&dbRow), nil
}

func (r *templateRepository) GetByName(ctx context.Context, name, domainName string) (*entity.Template, error) {
	var dbRow mapper.DatabaseTemplateWithDomain

	query := `SELECT t.id, t.name, t.domain_id, d.name as domain_name, t.template_data, t.title, t.description, t.is_active, t.created_at, t.updated_at 
			  FROM templates t
			  JOIN domains d ON t.domain_id = d.id
			  WHERE t.name = ? AND d.name = ?`
	err := r.db.QueryRowContext(ctx, query, name, domainName).Scan(
		&dbRow.ID,
		&dbRow.Name,
		&dbRow.DomainID,
		&dbRow.DomainName,
		&dbRow.TemplateData,
		&dbRow.Title,
		&dbRow.Description,
		&dbRow.IsActive,
		&dbRow.CreatedAt,
		&dbRow.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return mapper.ToTemplateEntityWithDomain(&dbRow), nil
}

func (r *templateRepository) List(ctx context.Context, domainName string, page, size int) ([]*entity.Template, int, error) {
	offset := (page - 1) * size

	// Get total count
	countQuery := `SELECT COUNT(*) FROM templates t JOIN domains d ON t.domain_id = d.id WHERE d.name = ?`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, domainName).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get templates
	query := `SELECT t.id, t.name, t.domain_id, t.template_data, t.title, t.description, t.is_active, t.created_at, t.updated_at 
			  FROM templates t
			  JOIN domains d ON t.domain_id = d.id
			  WHERE d.name = ?
			  ORDER BY t.updated_at DESC
			  LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, domainName, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var dbRow mapper.DatabaseTemplate
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.DomainID,
			&dbRow.TemplateData,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.IsActive,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		template := mapper.ToTemplateEntity(&dbRow)
		if template != nil {
			templates = append(templates, template)
		}
	}

	return templates, total, nil
}

func (r *templateRepository) ListActive(ctx context.Context, domainName string, page, size int) ([]*entity.Template, int, error) {
	offset := (page - 1) * size

	// Get total count
	countQuery := `SELECT COUNT(*) FROM templates t JOIN domains d ON t.domain_id = d.id WHERE d.name = ? AND t.is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, domainName).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get active templates
	query := `SELECT t.id, t.name, t.domain_id, t.template_data, t.title, t.description, t.is_active, t.created_at, t.updated_at 
			  FROM templates t
			  JOIN domains d ON t.domain_id = d.id
			  WHERE d.name = ? AND t.is_active = true
			  ORDER BY t.updated_at DESC
			  LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, domainName, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var dbRow mapper.DatabaseTemplate
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.DomainID,
			&dbRow.TemplateData,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.IsActive,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		template := mapper.ToTemplateEntity(&dbRow)
		if template != nil {
			templates = append(templates, template)
		}
	}

	return templates, total, nil
}

func (r *templateRepository) ListByType(ctx context.Context, domainName, templateType string, page, size int) ([]*entity.Template, int, error) {
	offset := (page - 1) * size

	// Get total count
	countQuery := `SELECT COUNT(*) FROM templates t 
				   JOIN domains d ON t.domain_id = d.id 
				   WHERE d.name = ? AND JSON_EXTRACT(t.template_data, '$.type') = ?`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, domainName, templateType).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get templates by type
	query := `SELECT t.id, t.name, t.domain_id, t.template_data, t.title, t.description, t.is_active, t.created_at, t.updated_at 
			  FROM templates t
			  JOIN domains d ON t.domain_id = d.id
			  WHERE d.name = ? AND JSON_EXTRACT(t.template_data, '$.type') = ?
			  ORDER BY t.updated_at DESC
			  LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, domainName, templateType, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var dbRow mapper.DatabaseTemplate
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.DomainID,
			&dbRow.TemplateData,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.IsActive,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		template := mapper.ToTemplateEntity(&dbRow)
		if template != nil {
			templates = append(templates, template)
		}
	}

	return templates, total, nil
}

func (r *templateRepository) Update(ctx context.Context, template *entity.Template) error {
	dbModel := mapper.FromTemplateEntity(template)

	query := `UPDATE templates 
			  SET template_data = ?, title = ?, description = ?, is_active = ?, updated_at = ?
			  WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query,
		dbModel.TemplateData,
		dbModel.Title,
		dbModel.Description,
		dbModel.IsActive,
		time.Now(),
		dbModel.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *templateRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM templates WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *templateRepository) Exists(ctx context.Context, name, domainName string) (bool, error) {
	query := `SELECT 1 FROM templates t
			  JOIN domains d ON t.domain_id = d.id
			  WHERE t.name = ? AND d.name = ?`
	var exists int
	err := r.db.QueryRowContext(ctx, query, name, domainName).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *templateRepository) GetBatch(ctx context.Context, ids []int) ([]*entity.Template, error) {
	if len(ids) == 0 {
		return []*entity.Template{}, nil
	}

	// Create placeholders for the IN clause
	placeholders := strings.Repeat("?,", len(ids)-1) + "?"
	query := fmt.Sprintf(`SELECT id, name, domain_id, template_data, title, description, is_active, created_at, updated_at 
						  FROM templates WHERE id IN (%s)`, placeholders)

	// Convert ids to interface slice
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var dbRow mapper.DatabaseTemplate
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.DomainID,
			&dbRow.TemplateData,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.IsActive,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		template := mapper.ToTemplateEntity(&dbRow)
		if template != nil {
			templates = append(templates, template)
		}
	}

	return templates, nil
}

func (r *templateRepository) GetDomainByTemplateID(ctx context.Context, templateID int) (*entity.Domain, error) {
	var dbRow mapper.DatabaseDomain

	query := `SELECT d.id, d.name, d.description, d.created_at, d.updated_at
			  FROM domains d
			  JOIN templates t ON d.id = t.domain_id
			  WHERE t.id = ?`
	err := r.db.QueryRowContext(ctx, query, templateID).Scan(
		&dbRow.ID,
		&dbRow.Name,
		&dbRow.Description,
		&dbRow.CreatedAt,
		&dbRow.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return mapper.ToDomainEntity(&dbRow), nil
}

func (r *templateRepository) FilterByAttributes(ctx context.Context, domainName string, filters []repository.AttributeFilter, page, size int) ([]*entity.Template, int, error) {
	if len(filters) == 0 {
		return r.List(ctx, domainName, page, size)
	}

	offset := (page - 1) * size

	// Build the WHERE clause for attribute filters
	var whereConditions []string
	var args []interface{}
	baseArgIndex := 1

	for _, filter := range filters {
		switch filter.Operator {
		case "equals":
			whereConditions = append(whereConditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM template_attributes ta
				JOIN attributes a ON ta.attribute_id = a.id
				WHERE ta.template_id = t.id AND a.name = ? AND ta.value = ?
			)`))
			args = append(args, filter.Name, filter.Value)
			baseArgIndex += 2
		case "contains":
			whereConditions = append(whereConditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM template_attributes ta
				JOIN attributes a ON ta.attribute_id = a.id
				WHERE ta.template_id = t.id AND a.name = ? AND ta.value LIKE ?
			)`))
			args = append(args, filter.Name, "%"+filter.Value+"%")
			baseArgIndex += 2
		}
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " AND " + strings.Join(whereConditions, " AND ")
	}

	// Add domain name to args
	finalArgs := append([]interface{}{domainName}, args...)

	// Get total count
	countQuery := fmt.Sprintf(`SELECT COUNT(DISTINCT t.id) FROM templates t
							  JOIN domains d ON t.domain_id = d.id
							  WHERE d.name = ?%s`, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, finalArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get templates
	query := fmt.Sprintf(`SELECT DISTINCT t.id, t.name, t.domain_id, t.template_data, t.title, t.description, t.is_active, t.created_at, t.updated_at 
						  FROM templates t
						  JOIN domains d ON t.domain_id = d.id
						  WHERE d.name = ?%s
						  ORDER BY t.updated_at DESC
						  LIMIT ? OFFSET ?`, whereClause)

	// Add pagination args
	finalArgs = append(finalArgs, size, offset)

	rows, err := r.db.QueryContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var dbRow mapper.DatabaseTemplate
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.DomainID,
			&dbRow.TemplateData,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.IsActive,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		template := mapper.ToTemplateEntity(&dbRow)
		if template != nil {
			templates = append(templates, template)
		}
	}

	return templates, total, nil
}

func (r *templateRepository) Clone(ctx context.Context, sourceID int, newName, newTitle, newDescription string) (*entity.Template, error) {
	// Get source template
	sourceTemplate, err := r.GetByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}

	// Create new template with same data but different name
	newTemplate, err := entity.NewTemplate(
		newName,
		sourceTemplate.TemplateData(),
		newTitle,
		newDescription,
		sourceTemplate.DomainID(),
	)
	if err != nil {
		return nil, err
	}

	// Create the new template
	err = r.Create(ctx, newTemplate)
	if err != nil {
		return nil, err
	}

	return newTemplate, nil
}

func (r *templateRepository) SetActive(ctx context.Context, id int, active bool) error {
	query := `UPDATE templates SET is_active = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, active, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *templateRepository) GetTemplatesByDomainID(ctx context.Context, domainID int, page, size int) ([]*entity.Template, int, error) {
	offset := (page - 1) * size

	// Get total count
	countQuery := `SELECT COUNT(*) FROM templates WHERE domain_id = ?`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, domainID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get templates
	query := `SELECT id, name, domain_id, template_data, title, description, is_active, created_at, updated_at 
			  FROM templates
			  WHERE domain_id = ?
			  ORDER BY updated_at DESC
			  LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, domainID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var dbRow mapper.DatabaseTemplate
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.DomainID,
			&dbRow.TemplateData,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.IsActive,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		template := mapper.ToTemplateEntity(&dbRow)
		if template != nil {
			templates = append(templates, template)
		}
	}

	return templates, total, nil
}

func (r *templateRepository) CountByType(ctx context.Context, domainName string) (map[string]int, error) {
	query := `SELECT JSON_EXTRACT(t.template_data, '$.type') as template_type, COUNT(*) as count
			  FROM templates t
			  JOIN domains d ON t.domain_id = d.id
			  WHERE d.name = ?
			  GROUP BY JSON_EXTRACT(t.template_data, '$.type')`

	rows, err := r.db.QueryContext(ctx, query, domainName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var templateType sql.NullString
		var count int
		err := rows.Scan(&templateType, &count)
		if err != nil {
			return nil, err
		}

		if templateType.Valid {
			counts[templateType.String] = count
		}
	}

	return counts, nil
}

func (r *templateRepository) GetRecentlyModified(ctx context.Context, domainName string, limit int) ([]*entity.Template, error) {
	query := `SELECT t.id, t.name, t.domain_id, t.template_data, t.title, t.description, t.is_active, t.created_at, t.updated_at 
			  FROM templates t
			  JOIN domains d ON t.domain_id = d.id
			  WHERE d.name = ?
			  ORDER BY t.updated_at DESC
			  LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, domainName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var dbRow mapper.DatabaseTemplate
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.DomainID,
			&dbRow.TemplateData,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.IsActive,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		template := mapper.ToTemplateEntity(&dbRow)
		if template != nil {
			templates = append(templates, template)
		}
	}

	return templates, nil
}

func (r *templateRepository) Search(ctx context.Context, domainName, query string, page, size int) ([]*entity.Template, int, error) {
	offset := (page - 1) * size
	searchPattern := "%" + query + "%"

	// Get total count
	countQuery := `SELECT COUNT(*) FROM templates t
				   JOIN domains d ON t.domain_id = d.id
				   WHERE d.name = ? AND (t.name LIKE ? OR t.title LIKE ? OR t.description LIKE ?)`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, domainName, searchPattern, searchPattern, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Search templates
	searchQuery := `SELECT t.id, t.name, t.domain_id, t.template_data, t.title, t.description, t.is_active, t.created_at, t.updated_at 
					FROM templates t
					JOIN domains d ON t.domain_id = d.id
					WHERE d.name = ? AND (t.name LIKE ? OR t.title LIKE ? OR t.description LIKE ?)
					ORDER BY 
						CASE 
							WHEN t.name LIKE ? THEN 1
							WHEN t.title LIKE ? THEN 2
							ELSE 3
						END,
						t.updated_at DESC
					LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, searchQuery, domainName, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var dbRow mapper.DatabaseTemplate
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Name,
			&dbRow.DomainID,
			&dbRow.TemplateData,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.IsActive,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		template := mapper.ToTemplateEntity(&dbRow)
		if template != nil {
			templates = append(templates, template)
		}
	}

	return templates, total, nil
}