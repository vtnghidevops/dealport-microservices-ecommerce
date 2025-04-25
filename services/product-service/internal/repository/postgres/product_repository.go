package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"product-service/internal/domain"

	"github.com/jackc/pgx/v4"
)

const dbTimeout = time.Second * 3

// ProductRepository implements the domain.ProductRepository interface
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// GetProductByID returns a product by ID
func (r *ProductRepository) GetProductByID(id int) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			p.id, p.type, p.name, p.slug, p.description, p.price, COALESCE(p.original_price, 0),
			COALESCE(p.discount, 0), p.category_id, p.category_slug, p.stock_quantity, 
			p.brand, p.features, p.shipping_info, COALESCE(p.orders, 0), p.created_at, p.updated_at,
			ROUND(COALESCE(AVG(pr.rating), 0), 1) as avg_rating,
			COUNT(DISTINCT pr.id) as review_count, p.ui_metadata
		FROM 
			products p
		LEFT JOIN 
			product_reviews pr ON p.id = pr.product_id
		WHERE 
			p.id = $1
		GROUP BY 
			p.id
	`

	var product domain.Product
	var features, shippingInfo, uiMetadata json.RawMessage
	var avgRating float64
	var reviewCount int

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Type,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.OriginalPrice,
		&product.Discount,
		&product.CategoryID,
		&product.CategorySlug,
		&product.StockQuantity,
		&product.Brand,
		&features,
		&shippingInfo,
		&product.Orders,
		&product.CreatedAt,
		&product.UpdatedAt,
		&avgRating,
		&reviewCount,
		&uiMetadata,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Store the ui_metadata
	if len(uiMetadata) > 0 {
		product.UIMetadata = uiMetadata
	}

	// Parse features JSON to the new structure
	if len(features) > 0 {
		// Try to parse as array of strings first
		err = json.Unmarshal(features, &product.Features)
		if err != nil {
			// If that fails, try to parse as map and convert
			var featuresMap map[string]string
			err = json.Unmarshal(features, &featuresMap)
			if err == nil {
				// Convert map to array of strings
				for _, value := range featuresMap {
					product.Features = append(product.Features, value)
				}
			}
		}
	}

	// Parse shipping info JSON
	if len(shippingInfo) > 0 {
		// Try to parse directly to ShippingInfo
		err = json.Unmarshal(shippingInfo, &product.ShippingInfo)
		if err != nil {
			// If that fails, try to parse as map and convert
			var shippingMap map[string]interface{}
			err = json.Unmarshal(shippingInfo, &shippingMap)
			if err == nil {
				// Convert map to ShippingInfo
				if courier, ok := shippingMap["express"].(bool); ok && courier {
					product.ShippingInfo.Courier = "Available"
				}
				if local, ok := shippingMap["standard"].(bool); ok && local {
					product.ShippingInfo.Local = "Free"
				}
				if ups, ok := shippingMap["ups"].(bool); ok && ups {
					product.ShippingInfo.Ups = "Available"
				}
				if global, ok := shippingMap["international"].(bool); ok && global {
					product.ShippingInfo.Global = "Available"
				} else {
					product.ShippingInfo.Global = "Not Available"
				}
			}
		}
	}

	// Set reviews summary
	product.ReviewsAvg = domain.ProductRating{
		AverageRating: avgRating,
		Count:         reviewCount,
	}

	// Get product images
	images, err := r.GetProductImages(product.ID)
	if err == nil && len(images) > 0 {
		product.Images = images

		// Create imgSlider from images if it wasn't loaded from the database
		if len(product.ImgSlider) == 0 {
			product.ImgSlider = []string{}
			for _, img := range images {
				product.ImgSlider = append(product.ImgSlider, img.URL)

				// Use primary image (is_primary=true) for ImageURL
				if img.IsPrimary {
					product.ImageURL = img.URL
				}
			}
		}

		// If no image is marked as primary, use the first one
		if product.ImageURL == "" && len(images) > 0 {
			product.ImageURL = images[0].URL
		}
	}

	// Get product tags
	tags, err := r.GetProductTags(product.ID)
	if err == nil {
		product.Tags = tags
	}

	return &product, nil
}

// GetProductBySlug returns a product by slug
func (r *ProductRepository) GetProductBySlug(slug string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			p.id, p.type, p.name, p.slug, p.description, p.price, COALESCE(p.original_price, 0), 
			COALESCE(p.discount, 0), p.category_id, p.category_slug, p.stock_quantity, 
			p.brand, p.features, p.shipping_info, COALESCE(p.orders, 0), p.created_at, p.updated_at,
			ROUND(COALESCE(AVG(pr.rating), 0), 1) as avg_rating,
			COUNT(DISTINCT pr.id) as review_count, p.ui_metadata
		FROM 
			products p
		LEFT JOIN 
			product_reviews pr ON p.id = pr.product_id
		WHERE 
			p.slug = $1
		GROUP BY 
			p.id
	`

	var product domain.Product
	var features, shippingInfo, uiMetadata json.RawMessage
	var avgRating float64
	var reviewCount int

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&product.ID,
		&product.Type,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.OriginalPrice,
		&product.Discount,
		&product.CategoryID,
		&product.CategorySlug,
		&product.StockQuantity,
		&product.Brand,
		&features,
		&shippingInfo,
		&product.Orders,
		&product.CreatedAt,
		&product.UpdatedAt,
		&avgRating,
		&reviewCount,
		&uiMetadata,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Store the ui_metadata
	if len(uiMetadata) > 0 {
		product.UIMetadata = uiMetadata
	}

	// Parse features JSON to the new structure
	if len(features) > 0 {
		// Try to parse as array of strings first
		err = json.Unmarshal(features, &product.Features)
		if err != nil {
			// If that fails, try to parse as map and convert
			var featuresMap map[string]string
			err = json.Unmarshal(features, &featuresMap)
			if err == nil {
				// Convert map to array of strings
				for _, value := range featuresMap {
					product.Features = append(product.Features, value)
				}
			}
		}
	}

	// Parse shipping info JSON
	if len(shippingInfo) > 0 {
		// Try to parse directly to ShippingInfo
		err = json.Unmarshal(shippingInfo, &product.ShippingInfo)
		if err != nil {
			// If that fails, try to parse as map and convert
			var shippingMap map[string]interface{}
			err = json.Unmarshal(shippingInfo, &shippingMap)
			if err == nil {
				// Convert map to ShippingInfo
				if courier, ok := shippingMap["express"].(bool); ok && courier {
					product.ShippingInfo.Courier = "Available"
				}
				if local, ok := shippingMap["standard"].(bool); ok && local {
					product.ShippingInfo.Local = "Free"
				}
				if ups, ok := shippingMap["ups"].(bool); ok && ups {
					product.ShippingInfo.Ups = "Available"
				}
				if global, ok := shippingMap["international"].(bool); ok && global {
					product.ShippingInfo.Global = "Available"
				} else {
					product.ShippingInfo.Global = "Not Available"
				}
			}
		}
	}

	// Set reviews summary
	product.ReviewsAvg = domain.ProductRating{
		AverageRating: avgRating,
		Count:         reviewCount,
	}

	// Log the ReviewsAvg values for debugging
	// log.Printf("Product %d (%s) has ReviewsAvg: count=%d, avg_rating=%.2f",
	// 	product.ID, product.Name, reviewCount, avgRating)

	// Get product images
	images, err := r.GetProductImages(product.ID)
	if err == nil && len(images) > 0 {
		product.Images = images

		// Create imgSlider from images
		product.ImgSlider = []string{}
		for _, img := range images {
			product.ImgSlider = append(product.ImgSlider, img.URL)

			// Use primary image (is_primary=true) for ImageURL
			if img.IsPrimary {
				product.ImageURL = img.URL
			}
		}

		// If no image is marked as primary, use the first one
		if product.ImageURL == "" && len(images) > 0 {
			product.ImageURL = images[0].URL
		}
	}

	// Get product tags
	product.Tags, err = r.GetProductTags(product.ID)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// List returns products with pagination and optional filtering
func (r *ProductRepository) GetAllProducts(page, pageSize int, filters map[string]string) ([]*domain.Product, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Build WHERE clause based on filters
	var whereClause string
	var args []interface{}
	var argCount int = 1

	if len(filters) > 0 {
		whereClause = " WHERE "
		var conditions []string

		if categoryID, ok := filters["category_id"]; ok {
			conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", argCount))
			args = append(args, categoryID)
			argCount++
		}

		if categorySlug, ok := filters["category_slug"]; ok {
			conditions = append(conditions, fmt.Sprintf("p.category_slug = $%d", argCount))
			args = append(args, categorySlug)
			argCount++
		}

		if brand, ok := filters["brand"]; ok {
			conditions = append(conditions, fmt.Sprintf("p.brand = $%d", argCount))
			args = append(args, brand)
			argCount++
		}

		if productType, ok := filters["type"]; ok {
			conditions = append(conditions, fmt.Sprintf("p.type = $%d", argCount))
			args = append(args, productType)
			argCount++
		}

		whereClause += strings.Join(conditions, " AND ")
	}

	// Count total matching products
	countQuery := "SELECT COUNT(*) FROM products p" + whereClause
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Prepare ORDER BY clause
	var orderClause string
	if orderBy, ok := filters["order_by"]; ok {
		orderDir := "ASC"
		if dir, ok := filters["order_dir"]; ok && (dir == "desc" || dir == "DESC") {
			orderDir = "DESC"
		}

		switch orderBy {
		case "price":
			orderClause = fmt.Sprintf(" ORDER BY p.price %s", orderDir)
		case "name":
			orderClause = fmt.Sprintf(" ORDER BY p.name %s", orderDir)
		case "created_at":
			orderClause = fmt.Sprintf(" ORDER BY p.created_at %s", orderDir)
		case "id":
			orderClause = fmt.Sprintf(" ORDER BY p.id %s", orderDir)
		default:
			orderClause = " ORDER BY p.id ASC"
		}
	} else {
		orderClause = " ORDER BY p.id ASC"
	}

	// Pagination clause
	var paginationClause string
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		paginationClause = fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
		args = append(args, pageSize, offset)
	}

	// Main query using subquery to get accurate review statistics
	mainQuery := `
	SELECT 
		p.id, p.type, p.name, p.slug, p.description, p.price, COALESCE(p.original_price, 0), 
		COALESCE(p.discount, 0), p.category_id, p.category_slug, p.stock_quantity, 
		p.brand, p.features, p.shipping_info, COALESCE(p.orders, 0), p.created_at, p.updated_at,
		COALESCE(r.avg_rating, 0) as avg_rating,
		COALESCE(r.review_count, 0) as review_count, p.ui_metadata
	FROM 
		products p
	LEFT JOIN (
		SELECT 
			product_id,
			ROUND(AVG(rating), 1) as avg_rating,
			COUNT(*) as review_count
		FROM 
			product_reviews
		GROUP BY 
			product_id
	) r ON p.id = r.product_id
	` + whereClause + orderClause + paginationClause

	// Execute main query
	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product

	for rows.Next() {
		var product domain.Product
		var features, shippingInfo, uiMetadata json.RawMessage
		var avgRating float64
		var reviewCount int

		err := rows.Scan(
			&product.ID,
			&product.Type,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Price,
			&product.OriginalPrice,
			&product.Discount,
			&product.CategoryID,
			&product.CategorySlug,
			&product.StockQuantity,
			&product.Brand,
			&features,
			&shippingInfo,
			&product.Orders,
			&product.CreatedAt,
			&product.UpdatedAt,
			&avgRating,
			&reviewCount,
			&uiMetadata,
		)
		if err != nil {
			return nil, 0, err
		}

		// Store the ui_metadata
		if len(uiMetadata) > 0 {
			product.UIMetadata = uiMetadata
		}

		// Parse features JSON to the new structure
		if len(features) > 0 {
			// Try to parse as array of strings first
			err = json.Unmarshal(features, &product.Features)
			if err != nil {
				// If that fails, try to parse as map and convert
				var featuresMap map[string]string
				err = json.Unmarshal(features, &featuresMap)
				if err == nil {
					// Convert map to array of strings
					for _, value := range featuresMap {
						product.Features = append(product.Features, value)
					}
				}
			}
		}

		// Parse shipping info JSON
		if len(shippingInfo) > 0 {
			// Try to parse directly to ShippingInfo
			err = json.Unmarshal(shippingInfo, &product.ShippingInfo)
			if err != nil {
				// If that fails, try to parse as map and convert
				var shippingMap map[string]interface{}
				err = json.Unmarshal(shippingInfo, &shippingMap)
				if err == nil {
					// Convert map to ShippingInfo
					if courier, ok := shippingMap["express"].(bool); ok && courier {
						product.ShippingInfo.Courier = "Available"
					}
					if local, ok := shippingMap["standard"].(bool); ok && local {
						product.ShippingInfo.Local = "Free"
					}
					if ups, ok := shippingMap["ups"].(bool); ok && ups {
						product.ShippingInfo.Ups = "Available"
					}
					if global, ok := shippingMap["international"].(bool); ok && global {
						product.ShippingInfo.Global = "Available"
					} else {
						product.ShippingInfo.Global = "Not Available"
					}
				}
			}
		}

		// Set reviews summary
		product.ReviewsAvg = domain.ProductRating{
			AverageRating: avgRating,
			Count:         reviewCount,
		}

		// Log ReviewsAvg for debugging
		// log.Printf("Product %d (%s): ReviewsAvg count=%d, avg_rating=%.2f",
		// 	product.ID, product.Name, reviewCount, avgRating)

		// Get product images
		images, err := r.GetProductImages(product.ID)
		if err == nil && len(images) > 0 {
			product.Images = images

			// Create imgSlider from images
			product.ImgSlider = []string{}
			for _, img := range images {
				product.ImgSlider = append(product.ImgSlider, img.URL)

				// Use primary image (is_primary=true) for ImageURL
				if img.IsPrimary {
					product.ImageURL = img.URL
				}
			}
			// If no image is marked as primary, use the first one
			if product.ImageURL == "" && len(images) > 0 {
				product.ImageURL = images[0].URL
			}
		}

		// Get product tags
		product.Tags, err = r.GetProductTags(product.ID)
		if err != nil {
			return nil, 0, err
		}

		products = append(products, &product)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

// Insert adds a new product
func (r *ProductRepository) Insert(product *domain.Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Reset sequence to fix duplicate key issue
	// This ensures the sequence counter is synced with the highest ID value in the table
	_, err = tx.ExecContext(ctx, `SELECT setval('products_id_seq', (SELECT COALESCE(MAX(id), 0) FROM products), true)`)
	if err != nil {
		return 0, fmt.Errorf("failed to reset sequence: %w", err)
	}

	// Insert product and get the auto-incremented ID
	query := `
		INSERT INTO products (
			type, name, description, slug, price, 
			original_price, discount, category_id, category_slug, 
			stock_quantity, brand, features, shipping_info, orders, ui_metadata,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		) RETURNING id
	`

	// Convert features array to JSON
	featuresJSON, err := json.Marshal(product.Features)
	if err != nil {
		return 0, err
	}

	// Convert shipping info to JSON
	shippingJSON, err := json.Marshal(product.ShippingInfo)
	if err != nil {
		return 0, err
	}

	// Set default values
	now := time.Now()
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}
	if product.UpdatedAt.IsZero() {
		product.UpdatedAt = now
	}

	var id int
	err = tx.QueryRowContext(ctx, query,
		product.Type,
		product.Name,
		product.Description,
		product.Slug,
		product.Price,
		product.OriginalPrice,
		product.Discount,
		product.CategoryID,
		product.CategorySlug,
		product.StockQuantity,
		product.Brand,
		featuresJSON,
		shippingJSON,
		product.Orders,
		product.UIMetadata,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// Save product images if any
	if len(product.Images) > 0 {
		// Reset the product_images sequence to avoid duplicate key violations
		_, err = tx.ExecContext(ctx, `SELECT setval('product_images_id_seq', (SELECT COALESCE(MAX(id), 0) FROM product_images), true)`)
		if err != nil {
			return id, fmt.Errorf("failed to reset product_images sequence: %w", err)
		}

		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO product_images (product_id, url, is_primary, display_order, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`)
		if err != nil {
			return id, fmt.Errorf("error preparing image insert: %w", err)
		}
		defer stmt.Close()

		for i, img := range product.Images {
			_, err = stmt.ExecContext(ctx,
				id,
				img.URL,
				img.IsPrimary,
				i, // Use index as display order
				now,
			)
			if err != nil {
				return id, fmt.Errorf("error inserting image: %w", err)
			}
		}
	}

	// Save product tags if any
	if len(product.Tags) > 0 {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO product_tags (product_id, tag)
			VALUES ($1, $2)
		`)
		if err != nil {
			return id, fmt.Errorf("error preparing tag insert: %w", err)
		}
		defer stmt.Close()

		for _, tag := range product.Tags {
			_, err = stmt.ExecContext(ctx, id, tag)
			if err != nil {
				return id, fmt.Errorf("error inserting tag: %w", err)
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	// Luôn cập nhật thông tin ReviewsAvg từ bảng product_reviews
	// Việc này được thực hiện sau khi commit transaction để không ảnh hưởng đến việc lưu chính sản phẩm
	// Đếm số lượng reviews trong database và tính rating trung bình
	var averageRating float64
	var reviewCount int

	countQuery := `
		SELECT 
			COUNT(*) as count, 
			COALESCE(AVG(rating), 0) as avg_rating 
		FROM product_reviews 
		WHERE product_id = $1
	`

	row := r.db.QueryRowContext(context.Background(), countQuery, id)
	err = row.Scan(&reviewCount, &averageRating)

	// Luôn cập nhật lại trường ReviewsAvg trong model để đảm bảo tính nhất quán
	if err == nil {
		// Cập nhật lại product.ReviewsAvg
		product.ReviewsAvg.Count = reviewCount
		product.ReviewsAvg.AverageRating = averageRating

		// Log thông tin đã cập nhật
		// log.Printf("Updated ReviewsAvg for product %d: count=%d, avg_rating=%.2f",
		// 	id, reviewCount, averageRating)
	} else {
		log.Printf("Error updating ReviewsAvg for product %d: %v", id, err)
	}

	return id, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Reset product ID sequence if needed (for testing)
	if product.ID > 0 {
		_, err = tx.ExecContext(ctx, "SELECT setval('products_id_seq', (SELECT MAX(id) FROM products), true)")
		if err != nil {
			return err
		}
	}

	// Marshal features to JSON
	featuresJSON, err := json.Marshal(product.Features)
	if err != nil {
		return err
	}

	// Marshal shipping info to JSON
	shippingInfoJSON, err := json.Marshal(product.ShippingInfo)
	if err != nil {
		return err
	}

	// Marshal UI metadata to JSON if it exists
	var uiMetadataJSON = []byte("{}")
	if product.UIMetadata != nil && len(product.UIMetadata) > 0 {
		// Use the UIMetadata directly as it's already RawMessage ([]byte)
		uiMetadataJSON = product.UIMetadata
	}

	// Update the product
	query := `
		UPDATE products SET
			type = $1,
			name = $2,
			slug = $3,
			description = $4,
			price = $5,
			original_price = $6,
			discount = $7,
			category_id = $8,
			category_slug = $9,
			stock_quantity = $10,
			brand = $11,
			features = $12,
			shipping_info = $13,
			ui_metadata = $14,
			updated_at = NOW()
		WHERE id = $15
	`

	_, err = tx.ExecContext(ctx, query,
		product.Type,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.OriginalPrice,
		product.Discount,
		product.CategoryID,
		product.CategorySlug,
		product.StockQuantity,
		product.Brand,
		featuresJSON,
		shippingInfoJSON,
		uiMetadataJSON,
		product.ID,
	)

	if err != nil {
		return err
	}

	// Handle product images
	// First, check if we have at least one primary image
	hasPrimary := false
	for _, img := range product.Images {
		if img.IsPrimary {
			hasPrimary = true
			break
		}
	}

	// If no primary image is set but we have images, set the first one as primary
	if !hasPrimary && len(product.Images) > 0 {
		product.Images[0].IsPrimary = true
	}

	// Delete existing images for the product
	_, err = tx.ExecContext(ctx, "DELETE FROM product_images WHERE product_id = $1", product.ID)
	if err != nil {
		return err
	}

	// Reset product_images ID sequence
	_, err = tx.ExecContext(ctx, "SELECT setval('product_images_id_seq', (SELECT COALESCE(MAX(id), 0) FROM product_images), true)")
	if err != nil {
		return err
	}

	// Insert new images if any
	if len(product.Images) > 0 {
		for i, img := range product.Images {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO product_images (product_id, url, is_primary, display_order) VALUES ($1, $2, $3, $4)",
				product.ID, img.URL, img.IsPrimary, i)
			if err != nil {
				return err
			}
		}
	}

	// Delete existing tags for the product
	_, err = tx.ExecContext(ctx, "DELETE FROM product_tags WHERE product_id = $1", product.ID)
	if err != nil {
		return err
	}

	// Insert new tags if any
	if len(product.Tags) > 0 {
		for _, tag := range product.Tags {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO product_tags (product_id, tag) VALUES ($1, $2)",
				product.ID, tag)
			if err != nil {
				return err
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a product
func (r *ProductRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Bắt đầu transaction để đảm bảo tính nhất quán
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Xóa sản phẩm
	stmt := "DELETE FROM products WHERE id = $1"
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetProductImages returns all images for a product
func (r *ProductRepository) GetProductImages(productID int) ([]domain.ProductImage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, product_id, url, is_primary, display_order, created_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []domain.ProductImage

	for rows.Next() {
		var img domain.ProductImage
		err := rows.Scan(
			&img.ID,
			&img.ProductID,
			&img.URL,
			&img.IsPrimary,
			&img.DisplayOrder,
			&img.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}

// GetProductTags returns all tags for a product
func (r *ProductRepository) GetProductTags(productID int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT tag
		FROM product_tags
		WHERE product_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string

	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetProductReviews returns reviews for a product with pagination
func (r *ProductRepository) GetProductReviews(productID int, page, pageSize int) ([]*domain.ProductReview, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Count total reviews
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM product_reviews WHERE product_id = $1"
	err := r.db.QueryRowContext(ctx, countQuery, productID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated reviews
	query := `
		SELECT id, product_id, user_id, rating, comment, created_at
		FROM product_reviews
		WHERE product_id = $1
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, productID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []*domain.ProductReview

	for rows.Next() {
		var review domain.ProductReview
		err := rows.Scan(
			&review.ID,
			&review.ProductID,
			&review.UserID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		reviews = append(reviews, &review)
	}

	return reviews, totalCount, nil
}

// AddProductReview adds a review for a product
func (r *ProductRepository) AddProductReview(review *domain.ProductReview) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Reset sequence to fix duplicate key issue for product_reviews
	_, err = tx.ExecContext(ctx, `SELECT setval('product_reviews_id_seq', (SELECT COALESCE(MAX(id), 0) FROM product_reviews), true)`)
	if err != nil {
		return 0, fmt.Errorf("failed to reset product_reviews sequence: %w", err)
	}

	// Insert review and get the auto-incremented ID
	query := `
		INSERT INTO product_reviews (product_id, user_id, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	if review.CreatedAt.IsZero() {
		review.CreatedAt = now
	}

	var id int
	err = tx.QueryRowContext(ctx, query,
		review.ProductID,
		review.UserID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteProductReview deletes a review
func (r *ProductRepository) DeleteProductReview(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Lấy product_id trước khi xóa để tính toán lại reviewsAvg
	var productID int
	productIDQuery := "SELECT product_id FROM product_reviews WHERE id = $1"
	err := r.db.QueryRowContext(ctx, productIDQuery, id).Scan(&productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Review không tồn tại, không cần làm gì thêm
			return nil
		}
		return err
	}

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Reset sequence to fix duplicate key issue for product_reviews
	_, err = tx.ExecContext(ctx, `SELECT setval('product_reviews_id_seq', (SELECT COALESCE(MAX(id), 0) FROM product_reviews), true)`)
	if err != nil {
		return fmt.Errorf("failed to reset product_reviews sequence: %w", err)
	}

	stmt := "DELETE FROM product_reviews WHERE id = $1"
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Tính toán lại ReviewsAvg sau khi xóa review
	var averageRating float64
	var reviewCount int

	countQuery := `
		SELECT 
			COUNT(*) as count, 
			COALESCE(AVG(rating), 0) as avg_rating 
		FROM product_reviews 
		WHERE product_id = $1
	`

	err = r.db.QueryRowContext(context.Background(), countQuery, productID).Scan(&reviewCount, &averageRating)
	if err == nil {
		log.Printf("Updated ReviewsAvg after deleting review for product %d: count=%d, avg_rating=%.2f",
			productID, reviewCount, averageRating)
	} else {
		log.Printf("Error calculating ReviewsAvg after deleting review for product %d: %v", productID, err)
	}

	return nil
}

// UpdateProductReview updates an existing review
func (r *ProductRepository) UpdateProductReview(review *domain.ProductReview) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Reset sequence to fix duplicate key issue for product_reviews
	_, err = tx.ExecContext(ctx, `SELECT setval('product_reviews_id_seq', (SELECT COALESCE(MAX(id), 0) FROM product_reviews), true)`)
	if err != nil {
		return fmt.Errorf("failed to reset product_reviews sequence: %w", err)
	}

	stmt := `
		UPDATE product_reviews 
		SET rating = $1, comment = $2, updated_at = $3
		WHERE id = $4 AND product_id = $5
	`

	_, err = tx.ExecContext(ctx, stmt,
		review.Rating,
		review.Comment,
		time.Now(),
		review.ID,
		review.ProductID,
	)

	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Tính toán lại ReviewsAvg sau khi cập nhật review
	var averageRating float64
	var reviewCount int

	countQuery := `
		SELECT 
			COUNT(*) as count, 
			COALESCE(AVG(rating), 0) as avg_rating 
		FROM product_reviews 
		WHERE product_id = $1
	`

	err = r.db.QueryRowContext(context.Background(), countQuery, review.ProductID).Scan(&reviewCount, &averageRating)
	if err == nil {
		log.Printf("Updated ReviewsAvg after updating review for product %d: count=%d, avg_rating=%.2f",
			review.ProductID, reviewCount, averageRating)
	} else {
		log.Printf("Error calculating ReviewsAvg after updating review for product %d: %v", review.ProductID, err)
	}

	return nil
}

// CreateProduct creates a new product
func (r *ProductRepository) CreateProduct(product *domain.Product) (int, error) {
	return r.Insert(product)
}

// UpdateProduct updates an existing product
func (r *ProductRepository) UpdateProduct(product *domain.Product) error {
	return r.Update(product)
}

// DeleteProduct deletes a product by ID
func (r *ProductRepository) DeleteProduct(id int) error {
	return r.Delete(id)
}

// GetRandomTopRatedReviews returns random 5-star reviews for the HappyCustomers section
func (r *ProductRepository) GetRandomTopRatedReviews(limit int) ([]*domain.Testimonial, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			pr.id, 
			pr.user_id, 
			COALESCE(pr.user_name, 'Happy Customer') as user_name, 
			pr.comment as review, 
			pr.rating,
			pr.product_id
		FROM 
			product_reviews pr
		WHERE 
			pr.rating = 5 AND pr.comment != ''
		ORDER BY 
			RANDOM()
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testimonials []*domain.Testimonial

	// Default avatar URLs if user service is not implemented yet
	defaultAvatars := []string{
		"/images/common/avatars/emily.png",
		"/images/common/avatars/johnd.png",
		"/images/common/avatars/ahmedM.png",
		"/images/common/avatars/alexT.png",
		"/images/common/avatars/priyaR.png",
		"/images/common/avatars/davidH.png",
		"/images/common/avatars/ahmedN.png",
	}

	avatarIndex := 0
	for rows.Next() {
		var testimonial domain.Testimonial
		var userID string

		err := rows.Scan(
			&testimonial.ID,
			&userID,
			&testimonial.UserName,
			&testimonial.Review,
			&testimonial.Rating,
			&testimonial.ProductID,
		)
		if err != nil {
			return nil, err
		}

		// Use default avatars in rotation until user service is implemented
		testimonial.Avatar = defaultAvatars[avatarIndex%len(defaultAvatars)]
		avatarIndex++

		// Format the review text with quotes
		if !strings.HasPrefix(testimonial.Review, "\"") {
			testimonial.Review = "\"" + testimonial.Review + "\""
		}

		testimonials = append(testimonials, &testimonial)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return testimonials, nil
}

// AddProductImage adds an image for a product and returns the image ID
func (r *ProductRepository) AddProductImage(productID int, imageURL string, isPrimary bool, displayOrder int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// If this is primary, update all other images to not primary
	if isPrimary {
		_, err = tx.ExecContext(ctx, `
			UPDATE product_images
			SET is_primary = false
			WHERE product_id = $1
		`, productID)
		if err != nil {
			return 0, err
		}
	}

	// Insert the new image
	query := `
		INSERT INTO product_images (product_id, url, is_primary, display_order, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err = tx.QueryRowContext(ctx, query,
		productID,
		imageURL,
		isPrimary,
		displayOrder,
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteProductImage deletes a product image by ID
func (r *ProductRepository) DeleteProductImage(imageID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if this is a primary image
	var isPrimary bool
	var productID int
	err = tx.QueryRowContext(ctx, `
		SELECT is_primary, product_id FROM product_images WHERE id = $1
	`, imageID).Scan(&isPrimary, &productID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Image doesn't exist, nothing to do
			return nil
		}
		return err
	}

	// Delete the image
	_, err = tx.ExecContext(ctx, `
		DELETE FROM product_images WHERE id = $1
	`, imageID)
	if err != nil {
		return err
	}

	// If this was a primary image, set a new primary
	if isPrimary {
		// Find another image for this product
		var newPrimaryID int
		err = tx.QueryRowContext(ctx, `
			SELECT id FROM product_images 
			WHERE product_id = $1 
			ORDER BY display_order ASC LIMIT 1
		`, productID).Scan(&newPrimaryID)

		if err == nil {
			// Set as primary
			_, err = tx.ExecContext(ctx, `
				UPDATE product_images 
				SET is_primary = true 
				WHERE id = $1
			`, newPrimaryID)
			if err != nil {
				return err
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			// Only return error if it's not a "no rows" error
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// UpdateProductImageOrder updates the display order of a product image
func (r *ProductRepository) UpdateProductImageOrder(imageID int, displayOrder int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE product_images
		SET display_order = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, displayOrder, imageID)
	return err
}

// SetPrimaryProductImage sets the primary image for a product
func (r *ProductRepository) SetPrimaryProductImage(productID int, imageID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Set all images for this product to non-primary
	_, err = tx.ExecContext(ctx, `
		UPDATE product_images
		SET is_primary = false
		WHERE product_id = $1
	`, productID)
	if err != nil {
		return err
	}

	// Set the selected image as primary
	_, err = tx.ExecContext(ctx, `
		UPDATE product_images
		SET is_primary = true
		WHERE id = $1 AND product_id = $2
	`, imageID, productID)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
