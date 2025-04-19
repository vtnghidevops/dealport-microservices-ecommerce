# Product Service

This service manages products, categories, and reviews for the e-commerce platform.

## Structure

The service follows a layered architecture:

- `cmd/api`: Entry point of the application
- `internal/handler`: HTTP request handlers
- `internal/transport`: HTTP and gRPC transport layers
- `internal/service`: Business logic
- `internal/repository`: Data access layer
- `data`: Data models and database operations

## Running the Service

### Prerequisites

- Go 1.16 or higher
- PostgreSQL database
- Environment variables:
  - `DSN`: Database connection string (optional, defaults to local development settings)

### Start the service

```bash
cd services/product-service
go run ./cmd/api
```

The service will start on port 8082 by default.

## API Endpoints

### Health Check
- `GET /health`: Basic health check endpoint

### Products
- `GET /api/v1/products`: Get all products with filtering and pagination
- `GET /api/v1/products/{id}`: Get product by ID
- `GET /api/v1/products/slug/{slug}`: Get product by slug
- `POST /api/v1/products`: Create a new product
- `PUT /api/v1/products/{id}`: Update an existing product
- `PATCH /api/v1/products/{id}`: Update an existing product
- `DELETE /api/v1/products/{id}`: Delete a product

### Product Reviews
- `GET /api/v1/products/{id}/reviews`: Get reviews for a product
- `POST /api/v1/products/{id}/reviews`: Add a review to a product

### Categories
- `GET /api/v1/categories`: Get all categories
- `GET /api/v1/categories/{id}`: Get category by ID
- `GET /api/v1/categories/slug/{slug}`: Get category by slug
- `GET /api/v1/categories/tree`: Get the full category hierarchy

## Testing

### Unit Tests

Run unit tests for the handlers:

```bash
cd services/product-service
go test ./internal/handler
```

### Integration Tests

To run integration tests, you need to:

1. Start the service in a separate terminal
2. Set the `RUN_INTEGRATION_TESTS` environment variable
3. Run the tests

```bash
# Terminal 1: Start the service
cd services/product-service
go run ./cmd/api

# Terminal 2: Run the integration tests
cd services/product-service
export RUN_INTEGRATION_TESTS=1
go test ./tests
```

## Example API Usage

### Creating a Product

```bash
curl -X POST http://localhost:8082/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smartphone X",
    "slug": "smartphone-x",
    "description": "Latest smartphone with amazing features",
    "price": 799.99,
    "category_id": "electronics-id",
    "category_slug": "electronics",
    "stock_quantity": 100,
    "brand": "TechBrand"
  }'
```

### Getting Products

```bash
# Get all products with filtering
curl "http://localhost:8082/api/v1/products?page=1&limit=10&category_slug=electronics&brand=TechBrand"

# Get product by ID
curl http://localhost:8082/api/v1/products/your-product-id
```

### Partial Update with PATCH

The PATCH endpoint allows updating specific fields of a product without requiring the full product object.

```bash
curl -X PATCH http://localhost:8082/api/v1/products/7d1e7337-2722-4a20-a685-3e53eaa51669 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 899.99,
    "discount": 200.00,
    "stock_quantity": 45
  }'
```

This request will only update the price, discount, and stock quantity of the product, leaving all other fields unchanged.

### Complex Partial Update

You can also update nested objects like features, shipping info, images, and tags:

```bash
curl -X PATCH http://localhost:8082/api/v1/products/7d1e7337-2722-4a20-a685-3e53eaa51669 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 899.99,
    "features": {
      "display": "6.1-inch Super Retina XDR",
      "chip": "A16 Bionic",
      "camera": "48MP Main | Ultra Wide | Telephoto",
      "battery": "Up to 25 hours video playback",
      "color": "Deep Purple"
    },
    "tags": ["phone", "premium", "apple", "ios", "sale"]
  }'
```

This request will update the price, completely replace the features object with the new one, and update the tags array.

## Database Structure

The product service uses PostgreSQL with the following main tables:

- `products`: Main product information
- `product_images`: Images associated with products
- `product_tags`: Tags for products
- `product_reviews`: User reviews for products
- `categories`: Product categories

## Development

### Running the Service

```bash
# Start the PostgreSQL container
docker-compose up -d postgres-products

# Start the service
go run ./cmd/api
```

### Environment Variables

- `DSN`: Database connection string

## Example API Usage

### Creating a Product

```