# Product Service

This service manages products, categories, and related operations for the e-commerce platform.

## Features

- Product CRUD operations
- Category management
- Product reviews
- Product images
- Banners and Ads management

## Storage Options

The product service supports two storage options for product images:

1. **Local Storage** (default): Images are stored in the local filesystem under `./uploads/products/`
2. **MinIO Storage**: Images are stored in a MinIO object storage bucket

### MinIO Storage Configuration

To use MinIO for image storage, set the following environment variables:

```
STORAGE_PROVIDER=minio
MINIO_ENDPOINT=minio.deploy.io.vn
MINIO_ACCESS_KEY_ID=be-images
MINIO_SECRET_ACCESS_KEY=your_secret_key
MINIO_USE_SSL=true
MINIO_BUCKET_NAME=images
MINIO_LOCATION=us-east-1
MINIO_BASE_URL=https://minio.deploy.io.vn
MINIO_PRESIGNED_TTL=3600
```

When MinIO storage is enabled:

1. Images are uploaded to the specified MinIO bucket
2. The database stores relative URLs to the images
3. When retrieving images, the service generates presigned URLs with a configurable expiry time
4. URLs are cached to improve performance and reduce load on the MinIO server

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

### Product Images API

#### Upload a Product Image

```
POST /api/v1/products/{id}/images
```

**Parameters:**

- `id` (path parameter): The ID of the product to upload the image for

**Form Data:**

- `image`: The image file (multipart/form-data)
- `isPrimary`: Whether this should be the primary product image (true/false)

**Response:**

```json
{
  "status": 201,
  "data": {
    "url": "/api/products/images/1234_1624567890.jpg"
  }
}
```

#### Delete a Product Image

```
DELETE /api/v1/products/{id}/images/{imageId}
```

**Parameters:**

- `id` (path parameter): The ID of the product
- `imageId` (path parameter): The ID of the image to delete

**Response:**

```json
{
  "status": 200,
  "message": "Image deleted successfully"
}
```

#### Set Primary Product Image

```
PUT /api/v1/products/{id}/images/{imageId}/primary
```

**Parameters:**

- `id` (path parameter): The ID of the product
- `imageId` (path parameter): The ID of the image to set as primary

**Response:**

```json
{
  "status": 200,
  "message": "Primary image set successfully"
}
```

The image URLs in the API responses can be used directly in your frontend application.
Images are stored in the backend filesystem and served through the `/api/products/images/` endpoint.

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

### PATCH Product Endpoint

The PATCH endpoint provides a flexible way to update only specific fields of a product without having to send the entire product object.

```
PATCH /api/v1/products/{id}
```

**Parameters:**

- `id` (path parameter): The ID of the product to update

**Request Body:**
A JSON object containing only the fields you want to update. The endpoint supports updating the following fields:

```json
{
  "name": "Updated Product Name",
  "description": "New description text",
  "slug": "updated-product-slug",
  "price": 99.99,
  "original_price": 129.99,
  "discount": 30.0,
  "category_id": 5,
  "category_slug": "electronics",
  "stock_quantity": 150,
  "type": "trending",
  "brand": "BrandName",
  "features": ["Feature 1", "Feature 2", "Feature 3"],
  "shipping_info": {
    "courier": "Express",
    "local": "Free",
    "ups": "Available",
    "global": "Available"
  },
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ],
  "tags": ["tag1", "tag2", "tag3"],
  "orders": 25,
  "ui_metadata": {
    "featured": true,
    "position": "top"
  }
}
```

You can include any combination of these fields in your request, and only the included fields will be updated.

**Response:**

```json
{
  "status": 200,
  "message": "Product updated successfully",
  "data": {
    // The full updated product object
  }
}
```

**Example: Update only the price and stock quantity**

```bash
curl -X PATCH http://localhost:8082/api/v1/products/123 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 89.99,
    "stock_quantity": 75
  }'
```

**Example: Update product images**

```bash
curl -X PATCH http://localhost:8082/api/v1/products/123 \
  -H "Content-Type: application/json" \
  -d '{
    "images": [
      "/api/products/images/123_image1.jpg",
      "/api/products/images/123_image2.jpg",
      "/api/products/images/123_image3.jpg"
    ]
  }'
```

**Example: Update product type and add tags**

```bash
curl -X PATCH http://localhost:8082/api/v1/products/123 \
  -H "Content-Type: application/json" \
  -d '{
    "type": "top-sale",
    "tags": ["featured", "sale", "bestseller"]
  }'
```

The PATCH endpoint is particularly useful for operations like:

- Updating price or stock quantity without changing other fields
- Adding or removing product images
- Changing the product type or category
- Adding or removing tags
- Updating metadata for UI display

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

```
