// data/models.go
package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New creates a new instance of the data package
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Product:       Product{},
		Category:      Category{},
		ProductImage:  ProductImage{},
		ProductTag:    ProductTag{},
		ProductReview: ProductReview{},
	}
}

// Models wraps the different models used by the application
type Models struct {
	Product       Product
	Category      Category
	ProductImage  ProductImage
	ProductTag    ProductTag
	ProductReview ProductReview
}

// Product model represents a product in the store
type Product struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Slug          string          `json:"slug"`
	Description   string          `json:"description"`
	Type          string          `json:"type"`
	Price         float64         `json:"price"`
	OriginalPrice float64         `json:"original_price,omitempty"`
	Discount      float64         `json:"discount,omitempty"`
	ImageURL      string          `json:"image_url"`
	CategoryID    string          `json:"category_id"`
	CategorySlug  string          `json:"category_slug"`
	StockQuantity int             `json:"stock_quantity"`
	Brand         string          `json:"brand,omitempty"`
	Features      json.RawMessage `json:"features,omitempty"`
	ShippingInfo  json.RawMessage `json:"shipping_info,omitempty"`
	Images        []ProductImage  `json:"images,omitempty"`
	Tags          []string        `json:"tags,omitempty"`
	Reviews       ProductReviews  `json:"reviews,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ProductReviews represents a summary of product reviews
type ProductReviews struct {
	AverageRating float64 `json:"average_rating"`
	Count         int     `json:"count"`
}

// ProductImage represents an image associated with a product
type ProductImage struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	URL          string    `json:"url"`
	IsPrimary    bool      `json:"is_primary"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// ProductTag represents a tag for a product
type ProductTag struct {
	ProductID string `json:"product_id"`
	Tag       string `json:"tag"`
}

// ProductReview represents a review for a product
type ProductReview struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// Category represents a product category
type Category struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Slug            string              `json:"slug"`
	Description     string              `json:"description,omitempty"`
	ImageURL        string              `json:"image_url,omitempty"`
	Icon            string              `json:"icon,omitempty"`
	BannerURL       string              `json:"banner_url,omitempty"`
	ProductCount    int                 `json:"product_count"`
	IsActive        bool                `json:"is_active"`
	IsVisible       bool                `json:"is_visible"`
	DisplayOrder    int                 `json:"display_order"`
	MetaTitle       string              `json:"meta_title,omitempty"`
	MetaDescription string              `json:"meta_description,omitempty"`
	ParentID        string              `json:"parent_id,omitempty"`
	Level           int                 `json:"level"`
	Attributes      []CategoryAttribute `json:"attributes,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Children        []*Category         `json:"children,omitempty"`
}

// CategoryAttribute represents an attribute for a category
type CategoryAttribute struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"category_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Required   bool      `json:"required"`
	Options    []string  `json:"options,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// GetAll returns all products with optional filtering
func (p *Product) GetAll(page, pageSize int, filters map[string]string) ([]*Product, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Base query
	query := `
		SELECT 
			p.id, p.name, p.slug, p.description, p.type, p.price, COALESCE(p.original_price, 0) , 
			COALESCE(p.discount, 0), p.image_url, p.category_id, p.category_slug, p.stock_quantity, 
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

		if categoryID, ok := filters["category_id"]; ok {
			whereClause += fmt.Sprintf("p.category_id = $%d ", argCount)
			args = append(args, categoryID)
			argCount++
		}

		if categorySlug, ok := filters["category_slug"]; ok {
			if argCount > 1 {
				whereClause += "AND "
			}
			whereClause += fmt.Sprintf("p.category_slug = $%d ", argCount)
			args = append(args, categorySlug)
			argCount++
		}

		if brand, ok := filters["brand"]; ok {
			if argCount > 1 {
				whereClause += "AND "
			}
			whereClause += fmt.Sprintf("p.brand = $%d ", argCount)
			args = append(args, brand)
			argCount++
		}

		if productType, ok := filters["type"]; ok {
			if argCount > 1 {
				whereClause += "AND "
			}
			whereClause += fmt.Sprintf("p.type = $%d ", argCount)
			args = append(args, productType)
			argCount++
		}
	}

	// Count total matching products
	countQuery := "SELECT COUNT(*) FROM products p" + whereClause
	var totalCount int
	err := db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Complete the query with GROUP BY, pagination
	query += whereClause
	query += " GROUP BY p.id "

	if orderBy, ok := filters["order_by"]; ok {
		query += "ORDER BY " + orderBy + " "
		if orderDir, ok := filters["order_dir"]; ok {
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
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*Product

	for rows.Next() {
		var product Product
		var avgRating float64
		var reviewCount int

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Type,
			&product.Price,
			&product.OriginalPrice,
			&product.Discount,
			&product.ImageURL,
			&product.CategoryID,
			&product.CategorySlug,
			&product.StockQuantity,
			&product.Brand,
			&product.Features,
			&product.ShippingInfo,
			&product.CreatedAt,
			&product.UpdatedAt,
			&avgRating,
			&reviewCount,
		)
		if err != nil {
			return nil, 0, err
		}

		// Set reviews summary
		product.Reviews = ProductReviews{
			AverageRating: avgRating,
			Count:         reviewCount,
		}

		// Get product images
		product.Images, err = p.GetProductImages(product.ID)
		if err != nil {
			log.Println("Error getting product images:", err)
		}

		// Get product tags
		product.Tags, err = p.GetProductTags(product.ID)
		if err != nil {
			log.Println("Error getting product tags:", err)
		}

		products = append(products, &product)
	}

	return products, totalCount, nil
}

// GetProductImages returns all images for a product
func (p *Product) GetProductImages(productID string) ([]ProductImage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, product_id, url, is_primary, display_order, created_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY display_order ASC
	`

	rows, err := db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []ProductImage

	for rows.Next() {
		var img ProductImage
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
func (p *Product) GetProductTags(productID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT tag
		FROM product_tags
		WHERE product_id = $1
	`

	rows, err := db.QueryContext(ctx, query, productID)
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

// GetByID gets a product by ID
func (p *Product) GetByID(id string) (*Product, error) {
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

	var product Product
	var avgRating float64
	var reviewCount int

	err := db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Type,
		&product.Price,
		&product.OriginalPrice,
		&product.Discount,
		&product.ImageURL,
		&product.CategoryID,
		&product.CategorySlug,
		&product.StockQuantity,
		&product.Brand,
		&product.Features,
		&product.ShippingInfo,
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

	// Set reviews summary
	product.Reviews = ProductReviews{
		AverageRating: avgRating,
		Count:         reviewCount,
	}

	// Get product images
	product.Images, err = p.GetProductImages(product.ID)
	if err != nil {
		log.Println("Error getting product images:", err)
	}

	// Get product tags
	product.Tags, err = p.GetProductTags(product.ID)
	if err != nil {
		log.Println("Error getting product tags:", err)
	}

	return &product, nil
}

// GetBySlug gets a product by slug
func (p *Product) GetBySlug(slug string) (*Product, error) {
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

	var product Product
	var avgRating float64
	var reviewCount int

	err := db.QueryRowContext(ctx, query, slug).Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Type,
		&product.Price,
		&product.OriginalPrice,
		&product.Discount,
		&product.ImageURL,
		&product.CategoryID,
		&product.CategorySlug,
		&product.StockQuantity,
		&product.Brand,
		&product.Features,
		&product.ShippingInfo,
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

	// Set reviews summary
	product.Reviews = ProductReviews{
		AverageRating: avgRating,
		Count:         reviewCount,
	}

	// Get product images
	product.Images, err = p.GetProductImages(product.ID)
	if err != nil {
		log.Println("Error getting product images:", err)
	}

	// Get product tags
	product.Tags, err = p.GetProductTags(product.ID)
	if err != nil {
		log.Println("Error getting product tags:", err)
	}

	return &product, nil
}

// Insert adds a new product to the database
func (p *Product) Insert(product Product) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
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

	var id string
	err = tx.QueryRowContext(ctx, stmt,
		product.ID,
		product.Name,
		product.Slug,
		product.Description,
		product.Type,
		product.Price,
		product.OriginalPrice,
		product.Discount,
		product.ImageURL,
		product.CategoryID,
		product.CategorySlug,
		product.StockQuantity,
		product.Brand,
		product.Features,
		product.ShippingInfo,
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
				product_id, url, is_primary, display_order, created_at
			) VALUES ($1, $2, $3, $4, $5)
		`

		for i, img := range product.Images {
			_, err = tx.ExecContext(ctx, imageStmt,
				id,
				img.URL,
				img.IsPrimary,
				i,
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

// Update updates a product in the database
func (p *Product) Update(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
		product.OriginalPrice,
		product.Discount,
		product.ImageURL,
		product.CategoryID,
		product.CategorySlug,
		product.StockQuantity,
		product.Brand,
		product.Features,
		product.ShippingInfo,
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
				product_id, url, is_primary, display_order, created_at
			) VALUES ($1, $2, $3, $4, $5)
		`

		for i, img := range product.Images {
			_, err = tx.ExecContext(ctx, imageStmt,
				product.ID,
				img.URL,
				img.IsPrimary,
				i,
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

// Delete removes a product from the database
func (p *Product) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := "DELETE FROM products WHERE id = $1"

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// GetProductReviews gets all reviews for a product
func (p *Product) GetProductReviews(productID string, page, pageSize int) ([]*ProductReview, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Count total reviews
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM product_reviews WHERE product_id = $1"
	err := db.QueryRowContext(ctx, countQuery, productID).Scan(&totalCount)
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
	rows, err := db.QueryContext(ctx, query, productID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []*ProductReview

	for rows.Next() {
		var review ProductReview
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

		// TODO: Get user name from auth service if needed

		reviews = append(reviews, &review)
	}

	return reviews, totalCount, nil
}

// data/models.go (continued)

// AddProductReview adds a review for a product
func (p *Product) AddProductReview(review ProductReview) (string, error) {
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
	err := db.QueryRowContext(ctx, stmt,
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

// GetAllCategories returns all categories with optional filtering
func (c *Category) GetAllCategories(filters map[string]string) ([]*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Start building the query
	baseQuery := `
		SELECT 
			id, name, slug, description, image_url, icon, banner_url, product_count,
			is_active, is_visible, display_order, meta_title, meta_description, 
			parent_id, level, created_at, updated_at
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
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category

	for rows.Next() {
		var category Category
		var parentID sql.NullString
		var bannerURL, description, imageURL, icon, metaTitle, metaDescription sql.NullString

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
		attributes, err := c.GetCategoryAttributes(category.ID)
		if err != nil {
			log.Println("Error getting category attributes:", err)
		} else {
			category.Attributes = attributes
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

// GetCategoryByID returns a category by ID
func (c *Category) GetCategoryByID(id string) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			id, name, slug, description, image_url, icon, banner_url, product_count,
			is_active, is_visible, display_order, meta_title, meta_description, 
			parent_id, level, created_at, updated_at
		FROM 
			categories
		WHERE 
			id = $1
	`

	var category Category
	var parentID sql.NullString
	var bannerURL, description, imageURL, icon, metaTitle, metaDescription sql.NullString

	err := db.QueryRowContext(ctx, query, id).Scan(
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
	attributes, err := c.GetCategoryAttributes(category.ID)
	if err != nil {
		log.Println("Error getting category attributes:", err)
	} else {
		category.Attributes = attributes
	}

	return &category, nil
}

// GetCategoryBySlug returns a category by slug
func (c *Category) GetCategoryBySlug(slug string) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			id, name, slug, description, image_url, icon, banner_url, product_count,
			is_active, is_visible, display_order, meta_title, meta_description, 
			parent_id, level, created_at, updated_at
		FROM 
			categories
		WHERE 
			slug = $1
	`

	var category Category
	var parentID sql.NullString
	var bannerURL, description, imageURL, icon, metaTitle, metaDescription sql.NullString

	err := db.QueryRowContext(ctx, query, slug).Scan(
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
	attributes, err := c.GetCategoryAttributes(category.ID)
	if err != nil {
		log.Println("Error getting category attributes:", err)
	} else {
		category.Attributes = attributes
	}

	return &category, nil
}

// GetCategoryAttributes returns attributes for a category
func (c *Category) GetCategoryAttributes(categoryID string) ([]CategoryAttribute, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			id, category_id, name, type, required, options, created_at
		FROM 
			category_attributes
		WHERE 
			category_id = $1
		ORDER BY 
			name ASC
	`

	rows, err := db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []CategoryAttribute

	for rows.Next() {
		var attr CategoryAttribute
		var options sql.NullString

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

		if options.Valid {
			// Parse options array from string
			// This is a simple approach - you might want to use proper array types in PostgreSQL
			opts := strings.Trim(options.String, "{}")
			if opts != "" {
				attr.Options = strings.Split(opts, ",")
			}
		}

		attributes = append(attributes, attr)
	}

	return attributes, nil
}

// GetCategoryTree returns the full category hierarchy
func (c *Category) GetCategoryTree() ([]*Category, error) {
	// First get all root categories (level 0)
	filters := map[string]string{
		"level": "0",
	}

	rootCategories, err := c.GetAllCategories(filters)
	if err != nil {
		return nil, err
	}

	// For each root category, get its children recursively
	for _, rootCat := range rootCategories {
		err = c.populateChildrenRecursively(rootCat)
		if err != nil {
			return nil, err
		}
	}

	return rootCategories, nil
}

// populateChildrenRecursively populates all child categories recursively
func (c *Category) populateChildrenRecursively(parent *Category) error {
	// Get direct children
	filters := map[string]string{
		"parent_id": parent.ID,
	}

	children, err := c.GetAllCategories(filters)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		parent.Children = children

		// Recursively get children of children
		for _, child := range children {
			err = c.populateChildrenRecursively(child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
