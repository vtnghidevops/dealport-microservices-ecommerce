package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"product-service/internal/domain"

	"github.com/google/uuid"
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

// GetByID returns a product by ID
func (r *ProductRepository) GetByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			p.id, p.name, p.slug, p.description, p.type, p.price, p.original_price, 
			p.discount, p.image_url, p.category_id, p.category_slug, p.stock_quantity, 
			p.brand, p.features, p.shipping_info, p.created_at, p.updated_at,
			COALESCE(AVG(pr.rating), 0) as avg_rating,
			COUNT(DISTINCT pr.id) as review_count
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
	var avgRating float64
	var reviewCount int
	var originalPrice, discount sql.NullFloat64
	var features, shippingInfo sql.NullString
	var brand sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Type,
		&product.Price,
		&originalPrice,
		&discount,
		&product.ImageURL,
		&product.CategoryID,
		&product.CategorySlug,
		&product.StockQuantity,
		&brand,
		&features,
		&shippingInfo,
		&product.CreatedAt,
		&product.UpdatedAt,
		&avgRating,
		&reviewCount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Set nullable fields
	if originalPrice.Valid {
		product.OriginalPrice = originalPrice.Float64
	}
	if discount.Valid {
		product.Discount = discount.Float64
	}
	if brand.Valid {
		product.Brand = brand.String
	}
	if features.Valid {
		product.Features = json.RawMessage(features.String)
	}
	if shippingInfo.Valid {
		product.ShippingInfo = json.RawMessage(shippingInfo.String)
	}

	// Set reviews summary
	product.Reviews = &domain.ProductReviews{
		AverageRating: avgRating,
		Count:         reviewCount,
	}

	// Get product images
	product.Images, err = r.GetProductImages(product.ID)
	if err != nil {
		return nil, err
	}

	// Get product tags
	product.Tags, err = r.GetProductTags(product.ID)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetBySlug returns a product by slug
func (r *ProductRepository) GetBySlug(slug string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			p.id, p.name, p.slug, p.description, p.type, p.price, p.original_price, 
			p.discount, p.image_url, p.category_id, p.category_slug, p.stock_quantity, 
			p.brand, p.features, p.shipping_info, p.created_at, p.updated_at,
			COALESCE(AVG(pr.rating), 0) as avg_rating,
			COUNT(DISTINCT pr.id) as review_count
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
	var avgRating float64
	var reviewCount int
	var originalPrice, discount sql.NullFloat64
	var features, shippingInfo sql.NullString
	var brand sql.NullString

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Type,
		&product.Price,
		&originalPrice,
		&discount,
		&product.ImageURL,
		&product.CategoryID,
		&product.CategorySlug,
		&product.StockQuantity,
		&brand,
		&features,
		&shippingInfo,
		&product.CreatedAt,
		&product.UpdatedAt,
		&avgRating,
		&reviewCount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Set nullable fields
	if originalPrice.Valid {
		product.OriginalPrice = originalPrice.Float64
	}
	if discount.Valid {
		product.Discount = discount.Float64
	}
	if brand.Valid {
		product.Brand = brand.String
	}
	if features.Valid {
		product.Features = json.RawMessage(features.String)
	}
	if shippingInfo.Valid {
		product.ShippingInfo = json.RawMessage(shippingInfo.String)
	}

	// Set reviews summary
	product.Reviews = &domain.ProductReviews{
		AverageRating: avgRating,
		Count:         reviewCount,
	}

	// Get product images
	product.Images, err = r.GetProductImages(product.ID)
	if err != nil {
		return nil, err
	}

	// Get product tags
	product.Tags, err = r.GetProductTags(product.ID)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// List returns products with optional filtering and pagination
func (r *ProductRepository) List(page, pageSize int, filters map[string]string) ([]*domain.Product, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Base query
	query := `
		SELECT 
			p.id, p.name, p.slug, p.description, p.type, p.price, p.original_price, 
			p.discount, p.image_url, p.category_id, p.category_slug, p.stock_quantity, 
			p.brand, p.features, p.shipping_info, p.created_at, p.updated_at,
			COALESCE(AVG(pr.rating), 0) as avg_rating,
			COUNT(DISTINCT pr.id) as review_count
		FROM 
			products p
		LEFT JOIN 
			product_reviews pr ON p.id = pr.product_id
	`

	// Add WHERE clauses based on filters
	var whereClause string
	var args []interface{}
	var argCount int = 1

	if len(filters) > 0 {
		whereClause = " WHERE "
		whereParts := []string{}

		if categoryID, ok := filters["category_id"]; ok && categoryID != "" {
			whereParts = append(whereParts, fmt.Sprintf("p.category_id = $%d", argCount))
			args = append(args, categoryID)
			argCount++
		}

		if categorySlug, ok := filters["category_slug"]; ok && categorySlug != "" {
			whereParts = append(whereParts, fmt.Sprintf("p.category_slug = $%d", argCount))
			args = append(args, categorySlug)
			argCount++
		}

		if brand, ok := filters["brand"]; ok && brand != "" {
			whereParts = append(whereParts, fmt.Sprintf("p.brand = $%d", argCount))
			args = append(args, brand)
			argCount++
		}

		if productType, ok := filters["type"]; ok && productType != "" {
			whereParts = append(whereParts, fmt.Sprintf("p.type = $%d", argCount))
			args = append(args, productType)
			argCount++
		}

		if len(whereParts) > 0 {
			whereClause += strings.Join(whereParts, " AND ")
		} else {
			whereClause = ""
		}
	}

	// Count total matching products
	countQuery := "SELECT COUNT(*) FROM products p" + whereClause
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Complete the query with GROUP BY, ORDER BY, pagination
	query += whereClause
	query += " GROUP BY p.id "

	if orderBy, ok := filters["order_by"]; ok && orderBy != "" {
		query += "ORDER BY " + orderBy + " "
		if orderDir, ok := filters["order_dir"]; ok && (orderDir == "ASC" || orderDir == "DESC") {
			query += orderDir + " "
		} else {
			query += "ASC "
		}
	} else {
		query += "ORDER BY p.name ASC "
	}

	// Add pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf("LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, pageSize, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product

	for rows.Next() {
		var product domain.Product
		var avgRating float64
		var reviewCount int
		var originalPrice, discount sql.NullFloat64
		var features, shippingInfo sql.NullString
		var brand sql.NullString

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Type,
			&product.Price,
			&originalPrice,
			&discount,
			&product.ImageURL,
			&product.CategoryID,
			&product.CategorySlug,
			&product.StockQuantity,
			&brand,
			&features,
			&shippingInfo,
			&product.CreatedAt,
			&product.UpdatedAt,
			&avgRating,
			&reviewCount,
		)
		if err != nil {
			return nil, 0, err
		}

		// Set nullable fields
		if originalPrice.Valid {
			product.OriginalPrice = originalPrice.Float64
		}
		if discount.Valid {
			product.Discount = discount.Float64
		}
		if brand.Valid {
			product.Brand = brand.String
		}
		if features.Valid {
			product.Features = json.RawMessage(features.String)
		}
		if shippingInfo.Valid {
			product.ShippingInfo = json.RawMessage(shippingInfo.String)
		}

		// Set reviews summary
		product.Reviews = &domain.ProductReviews{
			AverageRating: avgRating,
			Count:         reviewCount,
		}

		// Get product images
		product.Images, err = r.GetProductImages(product.ID)
		if err != nil {
			return nil, 0, err
		}

		// Get product tags
		product.Tags, err = r.GetProductTags(product.ID)
		if err != nil {
			return nil, 0, err
		}

		products = append(products, &product)
	}

	return products, totalCount, nil
}

// Create adds a new product
func (r *ProductRepository) Create(product *domain.Product) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Generate a new UUID if not provided
	if product.ID == "" {
		product.ID = uuid.New().String()
	}

	// Insert product
	stmt := `
		INSERT INTO products (
			id, name, slug, description, type, price, original_price, discount,
			image_url, category_id, category_slug, stock_quantity, brand,
			features, shipping_info, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id
	`

	now := time.Now()

	// Handle NULL values for optional fields
	var originalPrice, discount interface{}
	var brand, features, shippingInfo interface{}

	// Only set non-zero values, otherwise use NULL
	if product.OriginalPrice > 0 {
		originalPrice = product.OriginalPrice
	} else {
		originalPrice = nil
	}

	if product.Discount > 0 {
		discount = product.Discount
	} else {
		discount = nil
	}

	if product.Brand != "" {
		brand = product.Brand
	} else {
		brand = nil
	}

	if len(product.Features) > 0 {
		features = product.Features
	} else {
		features = nil
	}

	if len(product.ShippingInfo) > 0 {
		shippingInfo = product.ShippingInfo
	} else {
		shippingInfo = nil
	}

	var id string
	err = tx.QueryRowContext(ctx, stmt,
		product.ID,
		product.Name,
		product.Slug,
		product.Description,
		product.Type,
		product.Price,
		originalPrice,
		discount,
		product.ImageURL,
		product.CategoryID,
		product.CategorySlug,
		product.StockQuantity,
		brand,
		features,
		shippingInfo,
		now,
		now,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	// Insert product images if any
	if len(product.Images) > 0 {
		imageStmt := `
			INSERT INTO product_images (
				id, product_id, url, is_primary, display_order, created_at
			) VALUES ($1, $2, $3, $4, $5, $6)
		`

		for i, img := range product.Images {
			// Generate ID if not present
			imgID := img.ID
			if imgID == "" {
				imgID = uuid.New().String()
			}

			_, err = tx.ExecContext(ctx, imageStmt,
				imgID,
				id,
				img.URL,
				img.IsPrimary,
				i, // Use index as display order if not specified
				now,
			)
			if err != nil {
				return "", err
			}
		}
	}

	// Insert product tags if any
	if len(product.Tags) > 0 {
		tagStmt := `
			INSERT INTO product_tags (product_id, tag)
			VALUES ($1, $2)
		`

		for _, tag := range product.Tags {
			_, err = tx.ExecContext(ctx, tagStmt, id, tag)
			if err != nil {
				return "", err
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return id, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Handle NULL values for optional fields
	var originalPrice, discount interface{}
	var brand, features, shippingInfo interface{}

	// Only set non-zero values, otherwise use NULL
	if product.OriginalPrice > 0 {
		originalPrice = product.OriginalPrice
	} else {
		originalPrice = nil
	}

	if product.Discount > 0 {
		discount = product.Discount
	} else {
		discount = nil
	}

	if product.Brand != "" {
		brand = product.Brand
	} else {
		brand = nil
	}

	if len(product.Features) > 0 {
		features = product.Features
	} else {
		features = nil
	}

	if len(product.ShippingInfo) > 0 {
		shippingInfo = product.ShippingInfo
	} else {
		shippingInfo = nil
	}

	// Update product
	stmt := `
		UPDATE products
		SET 
			name = $1,
			slug = $2,
			description = $3,
			type = $4,
			price = $5,
			original_price = $6,
			discount = $7,
			image_url = $8,
			category_id = $9,
			category_slug = $10,
			stock_quantity = $11,
			brand = $12,
			features = $13,
			shipping_info = $14,
			updated_at = $15
		WHERE id = $16
	`

	_, err = tx.ExecContext(ctx, stmt,
		product.Name,
		product.Slug,
		product.Description,
		product.Type,
		product.Price,
		originalPrice,
		discount,
		product.ImageURL,
		product.CategoryID,
		product.CategorySlug,
		product.StockQuantity,
		brand,
		features,
		shippingInfo,
		time.Now(),
		product.ID,
	)
	if err != nil {
		return err
	}

	// Update images: first delete existing ones, then insert new ones
	_, err = tx.ExecContext(ctx, "DELETE FROM product_images WHERE product_id = $1", product.ID)
	if err != nil {
		return err
	}

	if len(product.Images) > 0 {
		imageStmt := `
			INSERT INTO product_images (
				id, product_id, url, is_primary, display_order, created_at
			) VALUES ($1, $2, $3, $4, $5, $6)
		`

		for i, img := range product.Images {
			// Generate ID if not present
			imgID := img.ID
			if imgID == "" {
				imgID = uuid.New().String()
			}

			_, err = tx.ExecContext(ctx, imageStmt,
				imgID,
				product.ID,
				img.URL,
				img.IsPrimary,
				i, // Use index as display order if not specified
				time.Now(),
			)
			if err != nil {
				return err
			}
		}
	}

	// Update tags: first delete existing ones, then insert new ones
	_, err = tx.ExecContext(ctx, "DELETE FROM product_tags WHERE product_id = $1", product.ID)
	if err != nil {
		return err
	}

	if len(product.Tags) > 0 {
		tagStmt := `
			INSERT INTO product_tags (product_id, tag)
			VALUES ($1, $2)
		`

		for _, tag := range product.Tags {
			_, err = tx.ExecContext(ctx, tagStmt, product.ID, tag)
			if err != nil {
				return err
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a product
func (r *ProductRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := "DELETE FROM products WHERE id = $1"

	_, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// GetProductImages returns all images for a product
func (r *ProductRepository) GetProductImages(productID string) ([]domain.ProductImage, error) {
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
func (r *ProductRepository) GetProductTags(productID string) ([]string, error) {
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
func (r *ProductRepository) GetProductReviews(productID string, page, pageSize int) ([]*domain.ProductReview, int, error) {
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
		ORDER BY created_at DESC
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
func (r *ProductRepository) AddProductReview(review *domain.ProductReview) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Generate a new UUID if not provided
	if review.ID == "" {
		review.ID = uuid.New().String()
	}

	stmt := `
		INSERT INTO product_reviews (
			id, product_id, user_id, rating, comment, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id string
	err := r.db.QueryRowContext(ctx, stmt,
		review.ID,
		review.ProductID,
		review.UserID,
		review.Rating,
		review.Comment,
		time.Now(),
	).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

// DeleteProductReview deletes a review
func (r *ProductRepository) DeleteProductReview(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := "DELETE FROM product_reviews WHERE id = $1"

	_, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
