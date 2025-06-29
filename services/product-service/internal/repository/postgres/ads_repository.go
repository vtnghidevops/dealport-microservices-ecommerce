package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"product-service/internal/domain"
	"time"
)

// AdsRepository implements the domain.AdsRepository interface
type AdsRepository struct {
	db *sql.DB
}

// NewAdsRepository creates a new instance of AdsRepository
func NewAdsRepository(db *sql.DB) *AdsRepository {
	return &AdsRepository{
		db: db,
	}
}

// GetAdsByLocation returns ads for a specific location (banner, display, gaming, new_fashion)
func (r *AdsRepository) GetAdsByLocation(location string) ([]*domain.AdsData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT ap.id, ap.reference_type, ap.reference_id, ap.custom_title, ap.custom_image_url, ap.ui_settings
		FROM ads_placement ap
		WHERE ap.location = $1 AND ap.is_active = true
		ORDER BY ap.display_order ASC
	`

	rows, err := r.db.QueryContext(ctx, query, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var placements []struct {
		ID             int
		ReferenceType  string
		ReferenceID    int
		CustomTitle    sql.NullString
		CustomImageURL sql.NullString
		UISettings     []byte
	}

	for rows.Next() {
		var p struct {
			ID             int
			ReferenceType  string
			ReferenceID    int
			CustomTitle    sql.NullString
			CustomImageURL sql.NullString
			UISettings     []byte
		}
		if err := rows.Scan(&p.ID, &p.ReferenceType, &p.ReferenceID, &p.CustomTitle, &p.CustomImageURL, &p.UISettings); err != nil {
			return nil, err
		}
		placements = append(placements, p)
	}

	var results []*domain.AdsData

	for _, p := range placements {
		// Based on ReferenceType, get data from the corresponding table
		if p.ReferenceType == "product" {
			// Get product information
			productQuery := `
				SELECT p.id, p.name, p.slug, COALESCE(pi.url, p.image_url) as image_url, p.price, p.category_id, c.slug as category_slug
				FROM products p
				LEFT JOIN (
					SELECT DISTINCT ON (product_id) product_id, url 
					FROM product_images 
					WHERE is_primary = true 
					ORDER BY product_id, display_order
				) pi ON p.id = pi.product_id
				LEFT JOIN categories c ON p.category_id = c.id
				WHERE p.id = $1
			`
			var product domain.AdsData
			err := r.db.QueryRowContext(ctx, productQuery, p.ReferenceID).Scan(
				&product.ID, &product.Name, &product.Slug, &product.ImageURL,
				&product.Price, &product.CategoryID, &product.CategorySlug,
			)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return nil, err
				}
				// Skip if product not found
				continue
			}

			product.Type = "product"
			if p.CustomTitle.Valid {
				product.CustomTitle = p.CustomTitle.String
			}
			if p.CustomImageURL.Valid {
				product.ImageURL = p.CustomImageURL.String
			}

			// Set default UI settings based on location if none provided
			if len(p.UISettings) == 0 || string(p.UISettings) == "null" || string(p.UISettings) == "{}" {
				defaultSettings := getDefaultUISettings(location)
				product.UISettings = defaultSettings
			} else {
				product.UISettings = p.UISettings
			}

			// For gaming ads, ensure we use the frontend expected field names
			if location == "gaming" {
				// Frontend expects 'subtitle' instead of 'name'
				if product.CustomTitle != "" {
					product.Name = product.CustomTitle
				}
			}

			results = append(results, &product)
		} else if p.ReferenceType == "category" {
			// Get category information
			categoryQuery := `
				SELECT id, name, slug, image_url
				FROM categories
				WHERE id = $1
			`
			var category domain.AdsData
			err := r.db.QueryRowContext(ctx, categoryQuery, p.ReferenceID).Scan(
				&category.ID, &category.Name, &category.Slug, &category.ImageURL,
			)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return nil, err
				}
				// Skip if category not found
				continue
			}

			category.Type = "category"
			if p.CustomTitle.Valid {
				category.CustomTitle = p.CustomTitle.String
			}
			if p.CustomImageURL.Valid {
				category.ImageURL = p.CustomImageURL.String
			}

			// Set default UI settings based on location if none provided
			if len(p.UISettings) == 0 || string(p.UISettings) == "null" || string(p.UISettings) == "{}" {
				defaultSettings := getDefaultUISettings(location)
				category.UISettings = defaultSettings
			} else {
				category.UISettings = p.UISettings
			}

			// For gaming ads, ensure we use the frontend expected field names
			if location == "gaming" {
				// Frontend expects 'subtitle' instead of 'name'
				if category.CustomTitle != "" {
					category.Name = category.CustomTitle
				}
			}

			results = append(results, &category)
		}
	}

	return results, nil
}

// getDefaultUISettings returns default UI settings based on ad location
func getDefaultUISettings(location string) json.RawMessage {
	switch location {
	case "banner":
		return json.RawMessage(`{"button_type":"primary","button_text":"Shop Now","has_more":false}`)
	case "display":
		return json.RawMessage(`{"button_type":"gradient","discount_img":""}`)
	case "gaming":
		// Gaming ads don't need any special UI settings
		return json.RawMessage(`{}`)
	case "new_fashion":
		return json.RawMessage(`{"button_type":"secondary"}`)
	default:
		return json.RawMessage(`{}`)
	}
}

// CreateAdsPlacement creates a new ads placement
func (r *AdsRepository) CreateAdsPlacement(ad *domain.AdsPlacement) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO ads_placement (
			location, reference_type, reference_id, display_order, 
			custom_title, custom_image_url, ui_settings, is_active, 
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	now := time.Now()
	if ad.CreatedAt.IsZero() {
		ad.CreatedAt = now
	}
	if ad.UpdatedAt.IsZero() {
		ad.UpdatedAt = now
	}

	// Handle empty strings as NULL values
	var customTitle, customImageURL sql.NullString
	if ad.CustomTitle != "" {
		customTitle = sql.NullString{String: ad.CustomTitle, Valid: true}
	}
	if ad.CustomImageURL != "" {
		customImageURL = sql.NullString{String: ad.CustomImageURL, Valid: true}
	}

	var id int
	err := r.db.QueryRowContext(ctx, query,
		ad.Location, ad.ReferenceType, ad.ReferenceID, ad.DisplayOrder,
		customTitle, customImageURL, ad.UISettings, ad.IsActive,
		ad.CreatedAt, ad.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetAdsPlacementByID returns an ads placement by ID
func (r *AdsRepository) GetAdsPlacementByID(id int) (*domain.AdsPlacement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			id, location, reference_type, reference_id, display_order,
			custom_title, custom_image_url, ui_settings, is_active,
			created_at, updated_at
		FROM ads_placement
		WHERE id = $1
	`

	var ad domain.AdsPlacement
	var customTitle, customImageURL sql.NullString
	var uiSettings []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ad.ID, &ad.Location, &ad.ReferenceType, &ad.ReferenceID, &ad.DisplayOrder,
		&customTitle, &customImageURL, &uiSettings, &ad.IsActive,
		&ad.CreatedAt, &ad.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("ads placement not found")
		}
		return nil, err
	}

	// Handle NULL values
	if customTitle.Valid {
		ad.CustomTitle = customTitle.String
	}
	if customImageURL.Valid {
		ad.CustomImageURL = customImageURL.String
	}
	ad.UISettings = uiSettings

	return &ad, nil
}

// UpdateAdsPlacement updates an existing ads placement
func (r *AdsRepository) UpdateAdsPlacement(ad *domain.AdsPlacement) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE ads_placement SET
			location = $1,
			reference_type = $2,
			reference_id = $3,
			display_order = $4,
			custom_title = $5,
			custom_image_url = $6,
			ui_settings = $7,
			is_active = $8,
			updated_at = $9
		WHERE id = $10
	`

	// Handle empty strings as NULL values
	var customTitle, customImageURL sql.NullString
	if ad.CustomTitle != "" {
		customTitle = sql.NullString{String: ad.CustomTitle, Valid: true}
	}
	if ad.CustomImageURL != "" {
		customImageURL = sql.NullString{String: ad.CustomImageURL, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		ad.Location, ad.ReferenceType, ad.ReferenceID, ad.DisplayOrder,
		customTitle, customImageURL, ad.UISettings, ad.IsActive,
		time.Now(), ad.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteAdsPlacement deletes an ads placement by ID
func (r *AdsRepository) DeleteAdsPlacement(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "DELETE FROM ads_placement WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllAdsPlacements returns all ads placements with pagination
func (r *AdsRepository) GetAllAdsPlacements(page, pageSize int) ([]*domain.AdsPlacement, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Count total ads placements
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM ads_placement"
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated ads placements
	query := `
		SELECT 
			id, location, reference_type, reference_id, display_order,
			custom_title, custom_image_url, ui_settings, is_active,
			created_at, updated_at
		FROM ads_placement
		ORDER BY location, display_order
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var ads []*domain.AdsPlacement

	for rows.Next() {
		var ad domain.AdsPlacement
		var customTitle, customImageURL sql.NullString
		var uiSettings []byte

		err := rows.Scan(
			&ad.ID, &ad.Location, &ad.ReferenceType, &ad.ReferenceID, &ad.DisplayOrder,
			&customTitle, &customImageURL, &uiSettings, &ad.IsActive,
			&ad.CreatedAt, &ad.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Handle NULL values
		if customTitle.Valid {
			ad.CustomTitle = customTitle.String
		}
		if customImageURL.Valid {
			ad.CustomImageURL = customImageURL.String
		}
		ad.UISettings = uiSettings

		ads = append(ads, &ad)
	}

	return ads, totalCount, nil
}
