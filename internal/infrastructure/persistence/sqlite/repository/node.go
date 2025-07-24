package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"url-db/internal/constants"
	"url-db/internal/domain/entity"
	"url-db/internal/domain/repository"
	"url-db/internal/infrastructure/persistence/sqlite/mapper"
)

type nodeRepository struct {
	db *sql.DB
}

// NewNodeRepository creates a new SQLite-based node repository
func NewNodeRepository(db *sql.DB) repository.NodeRepository {
	return &nodeRepository{db: db}
}

func (r *nodeRepository) Create(ctx context.Context, node *entity.Node) error {
	dbModel := mapper.FromNodeEntity(node)

	query := `INSERT INTO nodes (content, domain_id, title, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		dbModel.Content,
		dbModel.DomainID,
		dbModel.Title,
		dbModel.Description,
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

	node.SetID(int(id))
	return nil
}

func (r *nodeRepository) GetByID(ctx context.Context, id int) (*entity.Node, error) {
	var dbRow mapper.DatabaseNode

	query := `SELECT id, content, domain_id, title, description, created_at, updated_at FROM nodes WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dbRow.ID,
		&dbRow.Content,
		&dbRow.DomainID,
		&dbRow.Title,
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

	return mapper.ToNodeEntity(&dbRow), nil
}

func (r *nodeRepository) GetByURL(ctx context.Context, url, domainName string) (*entity.Node, error) {
	var dbRow mapper.DatabaseNode

	query := `SELECT n.id, n.content, n.domain_id, n.title, n.description, n.created_at, n.updated_at 
			  FROM nodes n 
			  JOIN domains d ON n.domain_id = d.id 
			  WHERE n.content = ? AND d.name = ?`
	err := r.db.QueryRowContext(ctx, query, url, domainName).Scan(
		&dbRow.ID,
		&dbRow.Content,
		&dbRow.DomainID,
		&dbRow.Title,
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

	return mapper.ToNodeEntity(&dbRow), nil
}

func (r *nodeRepository) List(ctx context.Context, domainName string, page, size int) ([]*entity.Node, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM nodes n JOIN domains d ON n.domain_id = d.id WHERE d.name = ?`
	err := r.db.QueryRowContext(ctx, countQuery, domainName).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * size

	// Get nodes with pagination
	query := `SELECT n.id, n.content, n.domain_id, n.title, n.description, n.created_at, n.updated_at 
			  FROM nodes n 
			  JOIN domains d ON n.domain_id = d.id 
			  WHERE d.name = ? 
			  ORDER BY n.created_at DESC 
			  LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, domainName, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []*entity.Node
	for rows.Next() {
		var dbRow mapper.DatabaseNode
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Content,
			&dbRow.DomainID,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		node := mapper.ToNodeEntity(&dbRow)
		if node != nil {
			nodes = append(nodes, node)
		}
	}

	return nodes, totalCount, nil
}

func (r *nodeRepository) Update(ctx context.Context, node *entity.Node) error {
	dbModel := mapper.FromNodeEntity(node)

	query := `UPDATE nodes SET title = ?, description = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query,
		dbModel.Title,
		dbModel.Description,
		dbModel.UpdatedAt,
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
		return errors.New(constants.ErrNodeNotFound)
	}

	return nil
}

func (r *nodeRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM nodes WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(constants.ErrNodeNotFound)
	}

	return nil
}

func (r *nodeRepository) Exists(ctx context.Context, url, domainName string) (bool, error) {
	var exists int
	query := `SELECT 1 FROM nodes n JOIN domains d ON n.domain_id = d.id WHERE n.content = ? AND d.name = ? LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, url, domainName).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *nodeRepository) GetBatch(ctx context.Context, ids []int) ([]*entity.Node, error) {
	if len(ids) == 0 {
		return []*entity.Node{}, nil
	}

	// Build query with placeholders
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}

	query := `SELECT id, content, domain_id, title, description, created_at, updated_at FROM nodes WHERE id IN (` +
		strings.Join(placeholders, ",") + `)`

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

	var nodes []*entity.Node
	for rows.Next() {
		var dbRow mapper.DatabaseNode
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Content,
			&dbRow.DomainID,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		node := mapper.ToNodeEntity(&dbRow)
		if node != nil {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

// GetDomainByNodeID retrieves the domain for a given node ID
func (r *nodeRepository) GetDomainByNodeID(ctx context.Context, nodeID int) (*entity.Domain, error) {
	query := `
		SELECT d.id, d.name, d.description, d.created_at, d.updated_at
		FROM domains d
		JOIN nodes n ON d.id = n.domain_id
		WHERE n.id = ?
	`
	
	var dbRow mapper.DatabaseDomain
	err := r.db.QueryRowContext(ctx, query, nodeID).Scan(
		&dbRow.ID,
		&dbRow.Name,
		&dbRow.Description,
		&dbRow.CreatedAt,
		&dbRow.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	
	return mapper.ToDomainEntity(&dbRow), nil
}

// FilterByAttributes retrieves nodes by domain with attribute filters
func (r *nodeRepository) FilterByAttributes(ctx context.Context, domainName string, filters []repository.AttributeFilter, page, size int) ([]*entity.Node, int, error) {
	if len(filters) == 0 {
		// No filters, return regular list
		return r.List(ctx, domainName, page, size)
	}

	// Build the SQL query with JOIN for each filter
	var joins []string
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Add domain condition
	conditions = append(conditions, "d.name = ?")
	args = append(args, domainName)
	argIndex++

	// Add a JOIN and condition for each filter
	for i, filter := range filters {
		joinAlias := "na" + string(rune('0'+i))
		attrAlias := "a" + string(rune('0'+i))
		
		joins = append(joins, 
			"INNER JOIN node_attributes "+joinAlias+" ON n.id = "+joinAlias+".node_id")
		joins = append(joins, 
			"INNER JOIN attributes "+attrAlias+" ON "+joinAlias+".attribute_id = "+attrAlias+".id")
		
		// Add attribute name condition
		conditions = append(conditions, attrAlias+".name = ?")
		args = append(args, filter.Name)
		argIndex++
		
		// Add value condition based on operator
		switch strings.ToLower(filter.Operator) {
		case "equals", "":
			conditions = append(conditions, joinAlias+".value = ?")
			args = append(args, filter.Value)
		case "contains":
			conditions = append(conditions, joinAlias+".value LIKE ?")
			args = append(args, "%"+filter.Value+"%")
		case "starts_with":
			conditions = append(conditions, joinAlias+".value LIKE ?")
			args = append(args, filter.Value+"%")
		case "ends_with":
			conditions = append(conditions, joinAlias+".value LIKE ?")
			args = append(args, "%"+filter.Value)
		default:
			// Default to equals for invalid operators
			conditions = append(conditions, joinAlias+".value = ?")
			args = append(args, filter.Value)
		}
		argIndex++
	}

	// Build the complete query
	baseQuery := `
		SELECT DISTINCT n.id, n.content, n.domain_id, n.title, n.description, n.created_at, n.updated_at
		FROM nodes n
		INNER JOIN domains d ON n.domain_id = d.id
		` + strings.Join(joins, " ") + `
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY n.created_at DESC
	`

	// Count query for total
	countQuery := `
		SELECT COUNT(DISTINCT n.id)
		FROM nodes n
		INNER JOIN domains d ON n.domain_id = d.id
		` + strings.Join(joins, " ") + `
		WHERE ` + strings.Join(conditions, " AND ")

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	offset := (page - 1) * size
	limitQuery := baseQuery + " LIMIT ? OFFSET ?"
	args = append(args, size, offset)

	// Execute the main query
	rows, err := r.db.QueryContext(ctx, limitQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []*entity.Node
	for rows.Next() {
		var dbRow mapper.DatabaseNode
		err := rows.Scan(
			&dbRow.ID,
			&dbRow.Content,
			&dbRow.DomainID,
			&dbRow.Title,
			&dbRow.Description,
			&dbRow.CreatedAt,
			&dbRow.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		node := mapper.ToNodeEntity(&dbRow)
		if node != nil {
			nodes = append(nodes, node)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return nodes, total, nil
}
