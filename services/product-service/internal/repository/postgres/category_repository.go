package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"product-service/internal/domain"

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
func (r *CategoryRepository) GetByID(id string) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			id, name, slug, description, image_url, icon, banner_url,
			product_count, is_active, is_visible, display_order,
			meta_title, meta_description, parent_id, level,
			created_at, updated_at
		FROM 
			categories
		WHERE 
			id = $1
	`

	var category domain.Category
	var parentID sql.NullString
	var description, imageURL, icon, bannerURL, metaTitle, metaDescription sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&description,
		&imageURL,
		&icon,
		&bannerURL,
		&category.ProductCount,
		&category.IsActive,
		&category.IsVisible,
		&category.DisplayOrder,
		&metaTitle,
		&metaDescription,
		&parentID,
		&category.Level,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Set nullable fields
	if description.Valid {
		category.Description = description.String
	}
	if imageURL.Valid {
		category.ImageURL = imageURL.String
	}
	if icon.Valid {
		category.Icon = icon.String
	}
	if bannerURL.Valid {
		category.BannerURL = bannerURL.String
	}
	if metaTitle.Valid {
		category.MetaTitle = metaTitle.String
	}
	if metaDescription.Valid {
		category.MetaDescription = metaDescription.String
	}
	if parentID.Valid {
		category.ParentID = parentID.String
	}

	// Get category attributes
	attributes, err := r.GetCategoryAttributes(category.ID)
	if err != nil {
		return nil, err
	}
	category.Attributes = attributes

	return &category, nil
}

// GetBySlug returns a category by slug
func (r *CategoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			id, name, slug, description, image_url, icon, banner_url,
			product_count, is_active, is_visible, display_order,
			meta_title, meta_description, parent_id, level,
			created_at, updated_at
		FROM 
			categories
		WHERE 
			slug = $1
	`

	var category domain.Category
	var parentID sql.NullString
	var description, imageURL, icon, bannerURL, metaTitle, metaDescription sql.NullString

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&description,
		&imageURL,
		&icon,
		&bannerURL,
		&category.ProductCount,
		&category.IsActive,
		&category.IsVisible,
		&category.DisplayOrder,
		&metaTitle,
		&metaDescription,
		&parentID,
		&category.Level,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Set nullable fields
	if description.Valid {
		category.Description = description.String
	}
	if imageURL.Valid {
		category.ImageURL = imageURL.String
	}
	if icon.Valid {
		category.Icon = icon.String
	}
	if bannerURL.Valid {
		category.BannerURL = bannerURL.String
	}
	if metaTitle.Valid {
		category.MetaTitle = metaTitle.String
	}
	if metaDescription.Valid {
		category.MetaDescription = metaDescription.String
	}
	if parentID.Valid {
		category.ParentID = parentID.String
	}

	// Get category attributes
	attributes, err := r.GetCategoryAttributes(category.ID)
	if err != nil {
		return nil, err
	}
	category.Attributes = attributes

	return &category, nil
}

// List returns categories with optional filtering
func (r *CategoryRepository) List(filters map[string]string) ([]*domain.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start building the query
	baseQuery := `
		SELECT 
			id, name, slug, description, image_url, icon, banner_url,
			product_count, is_active, is_visible, display_order,
			meta_title, meta_description, parent_id, level,
			created_at, updated_at
		FROM 
			categories
	`

	// Add WHERE clauses based on filters
	var whereClause string
	var args []interface{}
	var argCount int = 1

	if len(filters) > 0 {
		whereClause = " WHERE "
		whereParts := []string{}

		if parentID, ok := filters["parent_id"]; ok && parentID != "" {
			if parentID == "null" {
				whereParts = append(whereParts, "parent_id IS NULL")
			} else {
				whereParts = append(whereParts, fmt.Sprintf("parent_id = $%d", argCount))
				args = append(args, parentID)
				argCount++
			}
		}

		if isActive, ok := filters["is_active"]; ok && isActive != "" {
			isActiveBool := isActive == "true"
			whereParts = append(whereParts, fmt.Sprintf("is_active = $%d", argCount))
			args = append(args, isActiveBool)
			argCount++
		}

		if isVisible, ok := filters["is_visible"]; ok && isVisible != "" {
			isVisibleBool := isVisible == "true"
			whereParts = append(whereParts, fmt.Sprintf("is_visible = $%d", argCount))
			args = append(args, isVisibleBool)
			argCount++
		}

		if level, ok := filters["level"]; ok && level != "" {
			whereParts = append(whereParts, fmt.Sprintf("level = $%d", argCount))
			args = append(args, level)
			argCount++
		}

		if len(whereParts) > 0 {
			whereClause += strings.Join(whereParts, " AND ")
		} else {
			whereClause = ""
		}
	}

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
		orderClause = " ORDER BY level ASC, display_order ASC, name ASC"
	}

	// Build the final query
	query := baseQuery + whereClause + orderClause

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category

	for rows.Next() {
		var category domain.Category
		var parentID sql.NullString
		var description, imageURL, icon, bannerURL, metaTitle, metaDescription sql.NullString

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&description,
			&imageURL,
			&icon,
			&bannerURL,
			&category.ProductCount,
			&category.IsActive,
			&category.IsVisible,
			&category.DisplayOrder,
			&metaTitle,
			&metaDescription,
			&parentID,
			&category.Level,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set nullable fields
		if description.Valid {
			category.Description = description.String
		}
		if imageURL.Valid {
			category.ImageURL = imageURL.String
		}
		if icon.Valid {
			category.Icon = icon.String
		}
		if bannerURL.Valid {
			category.BannerURL = bannerURL.String
		}
		if metaTitle.Valid {
			category.MetaTitle = metaTitle.String
		}
		if metaDescription.Valid {
			category.MetaDescription = metaDescription.String
		}
		if parentID.Valid {
			category.ParentID = parentID.String
		}

		// Get category attributes
		attributes, err := r.GetCategoryAttributes(category.ID)
		if err != nil {
			return nil, err
		}
		category.Attributes = attributes

		categories = append(categories, &category)
	}

	return categories, nil
}

// GetCategoryTree returns the full category tree
func (r *CategoryRepository) GetCategoryTree() ([]*domain.Category, error) {
	// First get all root categories (level 0)
	filters := map[string]string{
		"level": "0",
	}

	rootCategories, err := r.List(filters)
	if err != nil {
		return nil, err
	}

	// For each root category, get its children recursively
	for _, rootCat := range rootCategories {
		err = r.populateChildrenRecursively(rootCat)
		if err != nil {
			return nil, err
		}
	}

	return rootCategories, nil
}

// populateChildrenRecursively populates all child categories recursively
func (r *CategoryRepository) populateChildrenRecursively(parent *domain.Category) error {
	// Get direct children
	filters := map[string]string{
		"parent_id": parent.ID,
	}

	children, err := r.List(filters)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		parent.Children = children

		// Recursively get children of children
		for _, child := range children {
			err = r.populateChildrenRecursively(child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetCategoryAttributes returns attributes for a category
func (r *CategoryRepository) GetCategoryAttributes(categoryID string) ([]domain.CategoryAttribute, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, category_id, name, type, required, options, created_at
		FROM category_attributes
		WHERE category_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []domain.CategoryAttribute

	for rows.Next() {
		var attr domain.CategoryAttribute
		var options []string

		err := rows.Scan(
			&attr.ID,
			&attr.CategoryID,
			&attr.Name,
			&attr.Type,
			&attr.Required,
			&options,
			&attr.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		attr.Options = options
		attributes = append(attributes, attr)
	}

	return attributes, nil
}
