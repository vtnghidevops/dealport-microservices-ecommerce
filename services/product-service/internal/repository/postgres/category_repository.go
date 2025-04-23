package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"product-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// CategoryRepository implements the domain.CategoryRepository interface
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new instance of CategoryRepository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// GetByID returns a category by ID
func (r *CategoryRepository) GetCategoryByID(id int) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			c.id, c.name, c.slug, c.description, c.image_url,
			c.is_active, c.is_visible, c.created_at, c.updated_at,
			COUNT(p.id) AS product_count
		FROM 
			categories c
		LEFT JOIN 
			products p ON c.id = p.category_id
		WHERE 
			c.id = $1
		GROUP BY 
			c.id
	`

	var category domain.Category
	var description, imageURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&description,
		&imageURL,
		&category.IsActive,
		&category.IsVisible,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.ProductCount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	if description.Valid {
		category.Description = description.String
	}
	if imageURL.Valid {
		category.ImageURL = imageURL.String
	}

	return &category, nil
}

// GetBySlug returns a category by slug
func (r *CategoryRepository) GetCategoryBySlug(slug string) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			c.id, c.name, c.slug, c.description, c.image_url,
			c.is_active, c.is_visible, c.created_at, c.updated_at,
			COUNT(p.id) AS product_count
		FROM 
			categories c
		LEFT JOIN 
			products p ON c.id = p.category_id
		WHERE 
			c.slug = $1
		GROUP BY 
			c.id
	`

	var category domain.Category
	var description, imageURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&description,
		&imageURL,
		&category.IsActive,
		&category.IsVisible,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.ProductCount,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	if description.Valid {
		category.Description = description.String
	}
	if imageURL.Valid {
		category.ImageURL = imageURL.String
	}

	return &category, nil
}

// List returns categories with optional filtering
func (r *CategoryRepository) GetAllCategories(filters map[string]string) ([]*domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Base query with product count calculation
	baseQuery := `
		SELECT 
			c.id, c.name, c.slug, c.description, c.image_url,
			c.is_active, c.is_visible, c.created_at, c.updated_at,
			COUNT(p.id) AS product_count
		FROM 
			categories c
		LEFT JOIN 
			products p ON c.id = p.category_id
	`

	// Add WHERE clauses based on filters
	var whereClause string
	var args []interface{}
	var argCount int = 1

	if len(filters) > 0 {
		whereClause = " WHERE "
		whereParts := []string{}

		if isActive, ok := filters["is_active"]; ok && isActive != "" {
			whereParts = append(whereParts, fmt.Sprintf("c.is_active = $%d", argCount))
			args = append(args, isActive == "true")
			argCount++
		}

		if isVisible, ok := filters["is_visible"]; ok && isVisible != "" {
			whereParts = append(whereParts, fmt.Sprintf("c.is_visible = $%d", argCount))
			args = append(args, isVisible == "true")
			argCount++
		}

		if len(whereParts) > 0 {
			whereClause += strings.Join(whereParts, " AND ")
		} else {
			whereClause = ""
		}
	}

	// Add GROUP BY clause
	groupByClause := " GROUP BY c.id"

	// Add ORDER BY clause
	var orderClause string

	if orderBy, ok := filters["order_by"]; ok && orderBy != "" {
		orderClause = " ORDER BY " + orderBy + " "
		if orderDir, ok := filters["order_dir"]; ok && (orderDir == "ASC" || orderDir == "DESC") {
			orderClause += orderDir
		} else {
			orderClause += "ASC"
		}
	} else {
		orderClause = " ORDER BY c.id ASC"
	}

	// Build the final query
	query := baseQuery + whereClause + groupByClause + orderClause

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category

	for rows.Next() {
		var category domain.Category
		var description, imageURL sql.NullString

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&description,
			&imageURL,
			&category.IsActive,
			&category.IsVisible,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.ProductCount,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			category.Description = description.String
		}
		if imageURL.Valid {
			category.ImageURL = imageURL.String
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

// Create adds a new category
func (r *CategoryRepository) CreateCategory(category *domain.Category) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Generate a new ID if not provided
	if category.ID == 0 {
		// Lấy chữ cái đầu tiên của category name (nếu có)
		prefix := "C" // Mặc định là "C" nếu name rỗng
		if category.Name != "" {
			prefix = strings.ToUpper(string(category.Name[0]))
		}

		// Tạo UUID và cắt ngắn để đảm bảo tổng độ dài không quá 36 ký tự
		uuidStr := uuid.New().String()
		// Giữ lại phần đầu của UUID và nối với prefix (tổng không quá 36 ký tự)
		maxLen := 35 - len(prefix)
		if len(uuidStr) > maxLen {
			uuidStr = uuidStr[:maxLen]
		}
		var err error // declare err variable
		category.ID, err = strconv.Atoi(prefix + uuidStr)
		if err != nil {
			return 0, err // handle error from Atoi conversion
		}
	}

	// Insert category - removing product_count as it will be calculated dynamically
	query := `INSERT INTO categories (id, name, slug, description, image_url, is_active, is_visible, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			 RETURNING id`

	now := time.Now()

	var id int
	err := r.db.QueryRowContext(ctx, query,
		category.ID,
		category.Name,
		category.Slug,
		category.Description,
		category.ImageURL,
		category.IsActive,
		category.IsVisible,
		now,
		now,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// Get the actual product count for the new category (should be 0)
	countQuery := `SELECT COUNT(*) FROM products WHERE category_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, id).Scan(&category.ProductCount)
	if err != nil {
		// Non-critical error, just set to 0
		category.ProductCount = 0
	}

	return id, nil
}

// Update updates an existing category
func (r *CategoryRepository) UpdateCategory(category *domain.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Update category without updating product_count
	query := `UPDATE categories
			 SET name = $1, slug = $2, description = $3, image_url = $4, 
			 is_active = $5, is_visible = $6, updated_at = $7
			 WHERE id = $8`

	_, err := r.db.ExecContext(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.ImageURL,
		category.IsActive,
		category.IsVisible,
		time.Now(),
		category.ID,
	)

	if err != nil {
		return err
	}

	// Update the product_count in memory to reflect actual count
	countQuery := `SELECT COUNT(*) FROM products WHERE category_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, category.ID).Scan(&category.ProductCount)
	if err != nil {
		// Non-critical error, just don't update the count
		return nil
	}

	return nil
}

// Delete removes a category
func (r *CategoryRepository) DeleteCategory(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

// SyncProductCounts updates product counts for all categories
// This method is now redundant since counts are calculated dynamically,
// but kept for backward compatibility
func (r *CategoryRepository) SyncProductCounts() error {
	// No-op as counts are calculated dynamically in queries
	return nil
}
