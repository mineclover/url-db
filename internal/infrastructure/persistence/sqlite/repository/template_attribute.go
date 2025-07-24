package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
)

// SQLiteTemplateAttributeRepository implements TemplateAttributeRepository using SQLite
type SQLiteTemplateAttributeRepository struct {
	db *sql.DB
}

// NewSQLiteTemplateAttributeRepository creates a new SQLite template attribute repository
func NewSQLiteTemplateAttributeRepository(db *sql.DB) repository.TemplateAttributeRepository {
	return &SQLiteTemplateAttributeRepository{db: db}
}

// CreateTemplateAttribute creates a new template attribute
func (r *SQLiteTemplateAttributeRepository) CreateTemplateAttribute(ctx context.Context, templateAttribute *entity.TemplateAttribute) error {
	query := `
		INSERT INTO template_attributes (template_id, attribute_id, value, order_index, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		templateAttribute.TemplateID(),
		templateAttribute.AttributeID(),
		templateAttribute.Value(),
		templateAttribute.OrderIndex(),
		templateAttribute.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to create template attribute: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	templateAttribute.SetID(int(id))
	return nil
}

// GetTemplateAttributeByID retrieves a template attribute by ID
func (r *SQLiteTemplateAttributeRepository) GetTemplateAttributeByID(ctx context.Context, id int) (*entity.TemplateAttribute, error) {
	query := `
		SELECT id, template_id, attribute_id, value, order_index, created_at
		FROM template_attributes
		WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanTemplateAttribute(row)
}

// UpdateTemplateAttribute updates a template attribute
func (r *SQLiteTemplateAttributeRepository) UpdateTemplateAttribute(ctx context.Context, templateAttribute *entity.TemplateAttribute) error {
	query := `
		UPDATE template_attributes
		SET value = ?, order_index = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		templateAttribute.Value(),
		templateAttribute.OrderIndex(),
		templateAttribute.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update template attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// DeleteTemplateAttributeByID deletes a template attribute by ID
func (r *SQLiteTemplateAttributeRepository) DeleteTemplateAttributeByID(ctx context.Context, id int) error {
	query := `DELETE FROM template_attributes WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete template attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// GetTemplateAttributes retrieves all template attributes for a template
func (r *SQLiteTemplateAttributeRepository) GetTemplateAttributes(ctx context.Context, templateID int) ([]*entity.TemplateAttribute, error) {
	query := `
		SELECT id, template_id, attribute_id, value, order_index, created_at
		FROM template_attributes
		WHERE template_id = ?
		ORDER BY attribute_id, COALESCE(order_index, 0)
	`
	rows, err := r.db.QueryContext(ctx, query, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template attributes: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributes(rows)
}

// GetTemplateAttributesWithDetails retrieves template attributes with attribute details
func (r *SQLiteTemplateAttributeRepository) GetTemplateAttributesWithDetails(ctx context.Context, templateID int) ([]*entity.TemplateAttributeWithDetails, error) {
	query := `
		SELECT ta.id, ta.template_id, ta.attribute_id, ta.value, ta.order_index, ta.created_at,
		       a.name, a.type, a.description
		FROM template_attributes ta
		JOIN attributes a ON ta.attribute_id = a.id
		WHERE ta.template_id = ?
		ORDER BY a.name, COALESCE(ta.order_index, 0)
	`
	rows, err := r.db.QueryContext(ctx, query, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template attributes with details: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributesWithDetails(rows)
}

// DeleteAllTemplateAttributes deletes all template attributes for a template
func (r *SQLiteTemplateAttributeRepository) DeleteAllTemplateAttributes(ctx context.Context, templateID int) error {
	query := `DELETE FROM template_attributes WHERE template_id = ?`
	_, err := r.db.ExecContext(ctx, query, templateID)
	if err != nil {
		return fmt.Errorf("failed to delete template attributes: %w", err)
	}
	return nil
}

// GetByAttributeID retrieves all template attributes for an attribute
func (r *SQLiteTemplateAttributeRepository) GetByAttributeID(ctx context.Context, attributeID int) ([]*entity.TemplateAttribute, error) {
	query := `
		SELECT id, template_id, attribute_id, value, order_index, created_at
		FROM template_attributes
		WHERE attribute_id = ?
		ORDER BY template_id, COALESCE(order_index, 0)
	`
	rows, err := r.db.QueryContext(ctx, query, attributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template attributes by attribute ID: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributes(rows)
}

// DeleteByAttributeID deletes all template attributes for an attribute
func (r *SQLiteTemplateAttributeRepository) DeleteByAttributeID(ctx context.Context, attributeID int) error {
	query := `DELETE FROM template_attributes WHERE attribute_id = ?`
	_, err := r.db.ExecContext(ctx, query, attributeID)
	if err != nil {
		return fmt.Errorf("failed to delete template attributes by attribute ID: %w", err)
	}
	return nil
}

// GetByTemplateAndAttribute retrieves a specific template attribute
func (r *SQLiteTemplateAttributeRepository) GetByTemplateAndAttribute(ctx context.Context, templateID, attributeID int) (*entity.TemplateAttribute, error) {
	query := `
		SELECT id, template_id, attribute_id, value, order_index, created_at
		FROM template_attributes
		WHERE template_id = ? AND attribute_id = ?
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, templateID, attributeID)
	return r.scanTemplateAttribute(row)
}

// Exists checks if a template attribute exists
func (r *SQLiteTemplateAttributeRepository) Exists(ctx context.Context, templateID, attributeID int) (bool, error) {
	query := `SELECT COUNT(1) FROM template_attributes WHERE template_id = ? AND attribute_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, templateID, attributeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check template attribute existence: %w", err)
	}
	return count > 0, nil
}

// CreateTemplateAttributesBatch creates multiple template attributes in a transaction
func (r *SQLiteTemplateAttributeRepository) CreateTemplateAttributesBatch(ctx context.Context, templateAttributes []*entity.TemplateAttribute) error {
	if len(templateAttributes) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO template_attributes (template_id, attribute_id, value, order_index, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, ta := range templateAttributes {
		result, err := stmt.ExecContext(ctx,
			ta.TemplateID(),
			ta.AttributeID(),
			ta.Value(),
			ta.OrderIndex(),
			ta.CreatedAt(),
		)
		if err != nil {
			return fmt.Errorf("failed to create template attribute: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}

		ta.SetID(int(id))
	}

	return tx.Commit()
}

// UpdateTemplateAttributesBatch updates multiple template attributes in a transaction
func (r *SQLiteTemplateAttributeRepository) UpdateTemplateAttributesBatch(ctx context.Context, templateAttributes []*entity.TemplateAttribute) error {
	if len(templateAttributes) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE template_attributes SET value = ?, order_index = ? WHERE id = ?`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, ta := range templateAttributes {
		_, err := stmt.ExecContext(ctx, ta.Value(), ta.OrderIndex(), ta.ID())
		if err != nil {
			return fmt.Errorf("failed to update template attribute: %w", err)
		}
	}

	return tx.Commit()
}

// DeleteBatch deletes multiple template attributes by IDs
func (r *SQLiteTemplateAttributeRepository) DeleteBatch(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	query := fmt.Sprintf("DELETE FROM template_attributes WHERE id IN (%s)", placeholders)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete template attributes: %w", err)
	}

	return nil
}

// FindTemplateAttributesByValue finds template attributes by value
func (r *SQLiteTemplateAttributeRepository) FindTemplateAttributesByValue(ctx context.Context, value string) ([]*entity.TemplateAttribute, error) {
	query := `
		SELECT id, template_id, attribute_id, value, order_index, created_at
		FROM template_attributes
		WHERE value LIKE ?
		ORDER BY template_id, attribute_id
	`
	rows, err := r.db.QueryContext(ctx, query, "%"+value+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find template attributes by value: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributes(rows)
}

// FindByAttributeNameAndValue finds template attributes by attribute name and value
func (r *SQLiteTemplateAttributeRepository) FindByAttributeNameAndValue(ctx context.Context, templateID int, attributeName, value string) ([]*entity.TemplateAttribute, error) {
	query := `
		SELECT ta.id, ta.template_id, ta.attribute_id, ta.value, ta.order_index, ta.created_at
		FROM template_attributes ta
		JOIN attributes a ON ta.attribute_id = a.id
		WHERE ta.template_id = ? AND a.name = ? AND ta.value LIKE ?
		ORDER BY ta.order_index
	`
	rows, err := r.db.QueryContext(ctx, query, templateID, attributeName, "%"+value+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find template attributes by name and value: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributes(rows)
}

// GetOrderedByTemplateAndAttribute gets ordered attributes for a template and attribute
func (r *SQLiteTemplateAttributeRepository) GetOrderedByTemplateAndAttribute(ctx context.Context, templateID, attributeID int) ([]*entity.TemplateAttribute, error) {
	query := `
		SELECT id, template_id, attribute_id, value, order_index, created_at
		FROM template_attributes
		WHERE template_id = ? AND attribute_id = ?
		ORDER BY COALESCE(order_index, 0)
	`
	rows, err := r.db.QueryContext(ctx, query, templateID, attributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered template attributes: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributes(rows)
}

// ReorderAttributes reorders template attributes
func (r *SQLiteTemplateAttributeRepository) ReorderAttributes(ctx context.Context, templateID, attributeID int, orderIndexes []int) error {
	// Get existing attributes for the template and attribute
	attributes, err := r.GetOrderedByTemplateAndAttribute(ctx, templateID, attributeID)
	if err != nil {
		return fmt.Errorf("failed to get existing attributes: %w", err)
	}

	if len(attributes) != len(orderIndexes) {
		return fmt.Errorf("number of attributes (%d) does not match number of order indexes (%d)", len(attributes), len(orderIndexes))
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE template_attributes SET order_index = ? WHERE id = ?`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i, attr := range attributes {
		_, err := stmt.ExecContext(ctx, orderIndexes[i], attr.ID())
		if err != nil {
			return fmt.Errorf("failed to update order index: %w", err)
		}
	}

	return tx.Commit()
}

// CountByTemplateID counts template attributes for a template
func (r *SQLiteTemplateAttributeRepository) CountByTemplateID(ctx context.Context, templateID int) (int, error) {
	query := `SELECT COUNT(*) FROM template_attributes WHERE template_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, templateID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count template attributes: %w", err)
	}
	return count, nil
}

// CountByAttributeID counts template attributes for an attribute
func (r *SQLiteTemplateAttributeRepository) CountByAttributeID(ctx context.Context, attributeID int) (int, error) {
	query := `SELECT COUNT(*) FROM template_attributes WHERE attribute_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, attributeID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count template attributes by attribute ID: %w", err)
	}
	return count, nil
}

// GetTemplateAttributeUsageStats gets attribute usage statistics for a domain
func (r *SQLiteTemplateAttributeRepository) GetTemplateAttributeUsageStats(ctx context.Context, domainName string) (map[string]int, error) {
	query := `
		SELECT a.name, COUNT(ta.id) as usage_count
		FROM attributes a
		LEFT JOIN template_attributes ta ON a.id = ta.attribute_id
		JOIN domains d ON a.domain_id = d.id
		WHERE d.name = ?
		GROUP BY a.id, a.name
		ORDER BY usage_count DESC
	`
	rows, err := r.db.QueryContext(ctx, query, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute usage stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			return nil, fmt.Errorf("failed to scan attribute usage stats: %w", err)
		}
		stats[name] = count
	}

	return stats, nil
}

// GetTemplateAttributesByDomain gets all template attributes for a domain
func (r *SQLiteTemplateAttributeRepository) GetTemplateAttributesByDomain(ctx context.Context, domainName string) ([]*entity.TemplateAttributeWithDetails, error) {
	query := `
		SELECT ta.id, ta.template_id, ta.attribute_id, ta.value, ta.order_index, ta.created_at,
		       a.name, a.type, a.description
		FROM template_attributes ta
		JOIN attributes a ON ta.attribute_id = a.id
		JOIN templates t ON ta.template_id = t.id
		JOIN domains d ON t.domain_id = d.id
		WHERE d.name = ?
		ORDER BY t.name, a.name, COALESCE(ta.order_index, 0)
	`
	rows, err := r.db.QueryContext(ctx, query, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get template attributes by domain: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributesWithDetails(rows)
}

// GetTemplatesWithAttribute gets templates that have a specific attribute
func (r *SQLiteTemplateAttributeRepository) GetTemplatesWithAttribute(ctx context.Context, attributeName, domainName string) ([]*entity.TemplateAttributeWithDetails, error) {
	query := `
		SELECT ta.id, ta.template_id, ta.attribute_id, ta.value, ta.order_index, ta.created_at,
		       a.name, a.type, a.description
		FROM template_attributes ta
		JOIN attributes a ON ta.attribute_id = a.id
		JOIN templates t ON ta.template_id = t.id
		JOIN domains d ON t.domain_id = d.id
		WHERE a.name = ? AND d.name = ?
		ORDER BY t.name, COALESCE(ta.order_index, 0)
	`
	rows, err := r.db.QueryContext(ctx, query, attributeName, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates with attribute: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributesWithDetails(rows)
}

// SetTemplateAttribute sets a single attribute for a template
func (r *SQLiteTemplateAttributeRepository) SetTemplateAttribute(ctx context.Context, templateID, attributeID int, value string, orderIndex *int) error {
	// Check if attribute already exists
	exists, err := r.Exists(ctx, templateID, attributeID)
	if err != nil {
		return fmt.Errorf("failed to check attribute existence: %w", err)
	}

	if exists {
		// Update existing attribute
		query := `UPDATE template_attributes SET value = ?, order_index = ? WHERE template_id = ? AND attribute_id = ?`
		_, err := r.db.ExecContext(ctx, query, value, orderIndex, templateID, attributeID)
		if err != nil {
			return fmt.Errorf("failed to update template attribute: %w", err)
		}
	} else {
		// Create new attribute
		ta, err := entity.NewTemplateAttribute(templateID, attributeID, value, orderIndex)
		if err != nil {
			return fmt.Errorf("failed to create template attribute entity: %w", err)
		}
		return r.CreateTemplateAttribute(ctx, ta)
	}

	return nil
}

// SetTemplateAttributes sets multiple attributes for a template (replaces all)
func (r *SQLiteTemplateAttributeRepository) SetTemplateAttributes(ctx context.Context, templateID int, attributes []repository.TemplateAttributeValue) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing attributes
	_, err = tx.ExecContext(ctx, "DELETE FROM template_attributes WHERE template_id = ?", templateID)
	if err != nil {
		return fmt.Errorf("failed to delete existing template attributes: %w", err)
	}

	// Insert new attributes
	if len(attributes) > 0 {
		query := `
			INSERT INTO template_attributes (template_id, attribute_id, value, order_index, created_at)
			SELECT ?, a.id, ?, ?, datetime('now')
			FROM attributes a
			JOIN domains d ON a.domain_id = d.id
			JOIN templates t ON t.domain_id = d.id
			WHERE a.name = ? AND t.id = ?
		`
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, attr := range attributes {
			_, err := stmt.ExecContext(ctx, templateID, attr.Value, attr.OrderIndex, attr.AttributeName, templateID)
			if err != nil {
				return fmt.Errorf("failed to create template attribute for %s: %w", attr.AttributeName, err)
			}
		}
	}

	return tx.Commit()
}

// DeleteTemplateAttribute deletes a specific attribute from a template
func (r *SQLiteTemplateAttributeRepository) DeleteTemplateAttribute(ctx context.Context, templateID, attributeID int) error {
	query := `DELETE FROM template_attributes WHERE template_id = ? AND attribute_id = ?`
	result, err := r.db.ExecContext(ctx, query, templateID, attributeID)
	if err != nil {
		return fmt.Errorf("failed to delete template attribute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// GetTemplatesByAttribute retrieves templates that have a specific attribute value
func (r *SQLiteTemplateAttributeRepository) GetTemplatesByAttribute(ctx context.Context, domainName, attributeName, attributeValue string) ([]*entity.Template, error) {
	query := `
		SELECT DISTINCT t.id, t.name, t.domain_id, t.template_data, t.title, t.description, t.is_active, t.created_at, t.updated_at
		FROM templates t
		JOIN template_attributes ta ON t.id = ta.template_id
		JOIN attributes a ON ta.attribute_id = a.id
		JOIN domains d ON t.domain_id = d.id
		WHERE d.name = ? AND a.name = ? AND ta.value = ?
		ORDER BY t.name
	`
	rows, err := r.db.QueryContext(ctx, query, domainName, attributeName, attributeValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates by attribute: %w", err)
	}
	defer rows.Close()

	var templates []*entity.Template
	for rows.Next() {
		var id, domainID int
		var name, templateData, title, description, createdAt, updatedAt string
		var isActive bool

		err := rows.Scan(&id, &name, &domainID, &templateData, &title, &description, &isActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		template, err := entity.NewTemplate(name, templateData, title, description, domainID)
		if err != nil {
			return nil, fmt.Errorf("failed to create template entity: %w", err)
		}

		template.SetID(id)
		template.SetActive(isActive)
		// Parse and set timestamps if needed
		// template.SetCreatedAt(parsedTime)
		// template.SetUpdatedAt(parsedTime)

		templates = append(templates, template)
	}

	return templates, nil
}

// GetTemplateAttributesByName retrieves specific attribute values for a template
func (r *SQLiteTemplateAttributeRepository) GetTemplateAttributesByName(ctx context.Context, templateID int, attributeNames []string) ([]*entity.TemplateAttribute, error) {
	if len(attributeNames) == 0 {
		return []*entity.TemplateAttribute{}, nil
	}

	placeholders := strings.Repeat("?,", len(attributeNames))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	query := fmt.Sprintf(`
		SELECT ta.id, ta.template_id, ta.attribute_id, ta.value, ta.order_index, ta.created_at
		FROM template_attributes ta
		JOIN attributes a ON ta.attribute_id = a.id
		WHERE ta.template_id = ? AND a.name IN (%s)
		ORDER BY a.name, COALESCE(ta.order_index, 0)
	`, placeholders)

	args := make([]interface{}, len(attributeNames)+1)
	args[0] = templateID
	for i, name := range attributeNames {
		args[i+1] = name
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get template attributes by name: %w", err)
	}
	defer rows.Close()

	return r.scanTemplateAttributes(rows)
}

// Transaction support methods
func (r *SQLiteTemplateAttributeRepository) BeginTx(ctx context.Context) (interface{}, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *SQLiteTemplateAttributeRepository) CommitTx(ctx context.Context, tx interface{}) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return sqlTx.Commit()
}

func (r *SQLiteTemplateAttributeRepository) RollbackTx(ctx context.Context, tx interface{}) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return sqlTx.Rollback()
}

// Helper methods for scanning results
func (r *SQLiteTemplateAttributeRepository) scanTemplateAttribute(row *sql.Row) (*entity.TemplateAttribute, error) {
	var id, templateID, attributeID int
	var value string
	var orderIndex sql.NullInt64
	var createdAt string

	err := row.Scan(&id, &templateID, &attributeID, &value, &orderIndex, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan template attribute: %w", err)
	}

	var orderPtr *int
	if orderIndex.Valid {
		order := int(orderIndex.Int64)
		orderPtr = &order
	}

	ta, err := entity.NewTemplateAttribute(templateID, attributeID, value, orderPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to create template attribute entity: %w", err)
	}

	ta.SetID(id)
	// Parse and set created_at if needed
	// ta.SetCreatedAt(parsedTime)

	return ta, nil
}

func (r *SQLiteTemplateAttributeRepository) scanTemplateAttributes(rows *sql.Rows) ([]*entity.TemplateAttribute, error) {
	var attributes []*entity.TemplateAttribute

	for rows.Next() {
		var id, templateID, attributeID int
		var value string
		var orderIndex sql.NullInt64
		var createdAt string

		err := rows.Scan(&id, &templateID, &attributeID, &value, &orderIndex, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template attribute: %w", err)
		}

		var orderPtr *int
		if orderIndex.Valid {
			order := int(orderIndex.Int64)
			orderPtr = &order
		}

		ta, err := entity.NewTemplateAttribute(templateID, attributeID, value, orderPtr)
		if err != nil {
			return nil, fmt.Errorf("failed to create template attribute entity: %w", err)
		}

		ta.SetID(id)
		// Parse and set created_at if needed
		// ta.SetCreatedAt(parsedTime)

		attributes = append(attributes, ta)
	}

	return attributes, nil
}

func (r *SQLiteTemplateAttributeRepository) scanTemplateAttributesWithDetails(rows *sql.Rows) ([]*entity.TemplateAttributeWithDetails, error) {
	var attributes []*entity.TemplateAttributeWithDetails

	for rows.Next() {
		var id, templateID, attributeID int
		var value, name, attrType, description string
		var orderIndex sql.NullInt64
		var createdAt string

		err := rows.Scan(&id, &templateID, &attributeID, &value, &orderIndex, &createdAt, &name, &attrType, &description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template attribute with details: %w", err)
		}

		var orderPtr *int
		if orderIndex.Valid {
			order := int(orderIndex.Int64)
			orderPtr = &order
		}

		ta, err := entity.NewTemplateAttribute(templateID, attributeID, value, orderPtr)
		if err != nil {
			return nil, fmt.Errorf("failed to create template attribute entity: %w", err)
		}

		ta.SetID(id)
		// Parse and set created_at if needed
		// ta.SetCreatedAt(parsedTime)

		taWithDetails := entity.NewTemplateAttributeWithDetails(ta, name, attrType, description)
		attributes = append(attributes, taWithDetails)
	}

	return attributes, nil
}
