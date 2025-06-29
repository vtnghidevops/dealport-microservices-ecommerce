// internal/repository/postgres/banner_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"product-service/internal/domain"
	"strings"
	"time"
)

// BannerRepository implements the domain.BannerRepository interface
type BannerRepository struct {
	db *sql.DB
}

// NewBannerRepository creates a new instance of BannerRepository
func NewBannerRepository(db *sql.DB) *BannerRepository {
	return &BannerRepository{
		db: db,
	}
}

// GetBannerByID returns a banner by ID
func (r *BannerRepository) GetBannerByID(id int) (*domain.Banner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, title, subtitle, description, discount, highlight_text, 
		       image_url, link_url, action_text, background_color, text_color, 
		       animation_type, is_active, priority, type, product_id, category_id,
		       created_at, updated_at
		FROM banners
		WHERE id = $1
	`

	var banner domain.Banner
	var subtitle, description, discount, highlightText, linkURL, actionText sql.NullString
	var backgroundColor, textColor, animationType sql.NullString
	var productID, categoryID sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&banner.ID,
		&banner.Title,
		&subtitle,
		&description,
		&discount,
		&highlightText,
		&banner.ImageURL,
		&linkURL,
		&actionText,
		&backgroundColor,
		&textColor,
		&animationType,
		&banner.IsActive,
		&banner.Priority,
		&banner.Type,
		&productID,
		&categoryID,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBannerNotFound
		}
		return nil, err
	}

	// Handle optional fields
	if subtitle.Valid {
		banner.Subtitle = subtitle.String
	}
	if description.Valid {
		banner.Description = description.String
	}
	if discount.Valid {
		banner.Discount = discount.String
	}
	if highlightText.Valid {
		banner.HighlightText = highlightText.String
	}
	if linkURL.Valid {
		banner.LinkURL = linkURL.String
	}
	if actionText.Valid {
		banner.ActionText = actionText.String
	}
	if backgroundColor.Valid {
		banner.BackgroundColor = backgroundColor.String
	}
	if textColor.Valid {
		banner.TextColor = textColor.String
	}
	if animationType.Valid {
		banner.AnimationType = animationType.String
	}
	if productID.Valid {
		pid := int(productID.Int32)
		banner.ProductID = &pid
	}
	if categoryID.Valid {
		cid := int(categoryID.Int32)
		banner.CategoryID = &cid
	}

	return &banner, nil
}

// GetBannersByType returns banners by type with optional limit
func (r *BannerRepository) GetBannersByType(bannerType string, limit int) ([]*domain.Banner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, title, subtitle, description, discount, highlight_text, 
		       image_url, link_url, action_text, background_color, text_color, 
		       animation_type, is_active, priority, type, product_id, category_id,
		       created_at, updated_at
		FROM banners
		WHERE is_active = true
	`

	args := []interface{}{}
	
	if bannerType != "" {
		query += " AND type = $1"
		args = append(args, bannerType)
	}

	query += " ORDER BY priority ASC, id ASC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*domain.Banner

	for rows.Next() {
		var banner domain.Banner
		var subtitle, description, discount, highlightText, linkURL, actionText sql.NullString
		var backgroundColor, textColor, animationType sql.NullString
		var productID, categoryID sql.NullInt32

		err := rows.Scan(
			&banner.ID,
			&banner.Title,
			&subtitle,
			&description,
			&discount,
			&highlightText,
			&banner.ImageURL,
			&linkURL,
			&actionText,
			&backgroundColor,
			&textColor,
			&animationType,
			&banner.IsActive,
			&banner.Priority,
			&banner.Type,
			&productID,
			&categoryID,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Handle optional fields
		if subtitle.Valid {
			banner.Subtitle = subtitle.String
		}
		if description.Valid {
			banner.Description = description.String
		}
		if discount.Valid {
			banner.Discount = discount.String
		}
		if highlightText.Valid {
			banner.HighlightText = highlightText.String
		}
		if linkURL.Valid {
			banner.LinkURL = linkURL.String
		}
		if actionText.Valid {
			banner.ActionText = actionText.String
		}
		if backgroundColor.Valid {
			banner.BackgroundColor = backgroundColor.String
		}
		if textColor.Valid {
			banner.TextColor = textColor.String
		}
		if animationType.Valid {
			banner.AnimationType = animationType.String
		}
		if productID.Valid {
			pid := int(productID.Int32)
			banner.ProductID = &pid
		}
		if categoryID.Valid {
			cid := int(categoryID.Int32)
			banner.CategoryID = &cid
		}

		banners = append(banners, &banner)
	}

	return banners, nil
}

// GetAllBanners returns all banners with pagination and filtering
func (r *BannerRepository) GetAllBanners(page, pageSize int, filters map[string]string) ([]*domain.Banner, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Base query
	baseQuery := `
		SELECT id, title, subtitle, description, discount, highlight_text, 
		       image_url, link_url, action_text, background_color, text_color, 
		       animation_type, is_active, priority, type, product_id, category_id,
		       created_at, updated_at
		FROM banners
	`

	// Add WHERE clauses based on filters
	var whereClause string
	var args []interface{}
	var argCount int = 1

	if len(filters) > 0 {
		whereClause = " WHERE "
		whereParts := []string{}

		if bannerType, ok := filters["type"]; ok && bannerType != "" {
			whereParts = append(whereParts, fmt.Sprintf("type = $%d", argCount))
			args = append(args, bannerType)
			argCount++
		}

		if isActive, ok := filters["is_active"]; ok && isActive != "" {
			whereParts = append(whereParts, fmt.Sprintf("is_active = $%d", argCount))
			args = append(args, isActive == "true")
			argCount++
		}

		if productID, ok := filters["product_id"]; ok && productID != "" {
			whereParts = append(whereParts, fmt.Sprintf("product_id = $%d", argCount))
			args = append(args, productID)
			argCount++
		}

		if categoryID, ok := filters["category_id"]; ok && categoryID != "" {
			whereParts = append(whereParts, fmt.Sprintf("category_id = $%d", argCount))
			args = append(args, categoryID)
			argCount++
		}

		if len(whereParts) > 0 {
			whereClause += strings.Join(whereParts, " AND ")
		} else {
			whereClause = ""
		}
	}

	// Add ORDER BY clause
	orderClause := " ORDER BY priority ASC, id ASC"

	// Count total banners for pagination
	countQuery := "SELECT COUNT(*) FROM banners" + whereClause
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	limitClause := fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	// Final query
	query := baseQuery + whereClause + orderClause + limitClause

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var banners []*domain.Banner

	for rows.Next() {
		var banner domain.Banner
		var subtitle, description, discount, highlightText, linkURL, actionText sql.NullString
		var backgroundColor, textColor, animationType sql.NullString
		var productID, categoryID sql.NullInt32

		err := rows.Scan(
			&banner.ID,
			&banner.Title,
			&subtitle,
			&description,
			&discount,
			&highlightText,
			&banner.ImageURL,
			&linkURL,
			&actionText,
			&backgroundColor,
			&textColor,
			&animationType,
			&banner.IsActive,
			&banner.Priority,
			&banner.Type,
			&productID,
			&categoryID,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		// Handle optional fields (tương tự như trên)
		if subtitle.Valid {
			banner.Subtitle = subtitle.String
		}
		if description.Valid {
			banner.Description = description.String
		}
		if discount.Valid {
			banner.Discount = discount.String
		}
		if highlightText.Valid {
			banner.HighlightText = highlightText.String
		}
		if linkURL.Valid {
			banner.LinkURL = linkURL.String
		}
		if actionText.Valid {
			banner.ActionText = actionText.String
		}
		if backgroundColor.Valid {
			banner.BackgroundColor = backgroundColor.String
		}
		if textColor.Valid {
			banner.TextColor = textColor.String
		}
		if animationType.Valid {
			banner.AnimationType = animationType.String
		}
		if productID.Valid {
			pid := int(productID.Int32)
			banner.ProductID = &pid
		}
		if categoryID.Valid {
			cid := int(categoryID.Int32)
			banner.CategoryID = &cid
		}

		banners = append(banners, &banner)
	}

	return banners, total, nil
}

// CreateBanner creates a new banner
func (r *BannerRepository) CreateBanner(banner *domain.Banner) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO banners (
			title, subtitle, description, discount, highlight_text, 
			image_url, link_url, action_text, background_color, text_color, 
			animation_type, is_active, priority, type, product_id, category_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id
	`

	now := time.Now()

	var id int
	err := r.db.QueryRowContext(
		ctx, 
		query,
		banner.Title,
		banner.Subtitle,
		banner.Description,
		banner.Discount,
		banner.HighlightText,
		banner.ImageURL,
		banner.LinkURL,
		banner.ActionText,
		banner.BackgroundColor,
		banner.TextColor,
		banner.AnimationType,
		banner.IsActive,
		banner.Priority,
		banner.Type,
		banner.ProductID,
		banner.CategoryID,
		now,
		now,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateBanner updates an existing banner
func (r *BannerRepository) UpdateBanner(banner *domain.Banner) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE banners SET
			title = $1,
			subtitle = $2,
			description = $3,
			discount = $4,
			highlight_text = $5,
			image_url = $6,
			link_url = $7,
			action_text = $8,
			background_color = $9,
			text_color = $10,
			animation_type = $11,
			is_active = $12,
			priority = $13,
			type = $14,
			product_id = $15,
			category_id = $16,
			updated_at = $17
		WHERE id = $18
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		banner.Title,
		banner.Subtitle,
		banner.Description,
		banner.Discount,
		banner.HighlightText,
		banner.ImageURL,
		banner.LinkURL,
		banner.ActionText,
		banner.BackgroundColor,
		banner.TextColor,
		banner.AnimationType,
		banner.IsActive,
		banner.Priority,
		banner.Type,
		banner.ProductID,
		banner.CategoryID,
		time.Now(),
		banner.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteBanner deletes a banner by ID
func (r *BannerRepository) DeleteBanner(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "DELETE FROM banners WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}