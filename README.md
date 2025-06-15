# E-commerce Microservices Architecture

This project implements a modern e-commerce platform using a microservices architecture built with Golang.

## System Architecture

```
┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │
│  Frontend App   │────▶│  Broker Service │
│                 │     │                 │
└─────────────────┘     └────────┬────────┘
                                 │
                                 │
         ┌─────────────┬─────────┼─────────┬─────────────┐
         │             │         │         │             │
         ▼             ▼         ▼         ▼             ▼
┌─────────────────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌─────────────────┐
│ Authentication  │ │ User  │ │Product│ │ Cart  │ │    Checkout     │
│    Service      │ │Service│ │Service│ │Service│ │     Service     │
└─────────────────┘ └───────┘ └───────┘ └───────┘ └─────────────────┘
         │             │         │         │             │
         │             │         │         │             │
         ▼             ▼         ▼         ▼             ▼
┌─────────────────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌─────────────────┐
│  Postgres DB    │ │Postgres│ │Postgres│ │Redis │ │    MongoDB     │
│  (Auth Data)    │ │  DB   │ │  DB   │ │  DB  │ │  (Order Data)   │
└─────────────────┘ └───────┘ └───────┘ └───────┘ └─────────────────┘
                                                          │
                                                          │
                                                          ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│ Logger Service  │◀───▶│ Listener Service│◀───▶│  Mail Service   │
│                 │     │                 │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    MongoDB      │     │    RabbitMQ     │     │   SMTP Server   │
│  (Log Data)     │     │  (Message Bus)  │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

## Core Services

### 1. Broker Service (API Gateway)

- Serves as the entry point for all client requests
- Routes requests to appropriate microservices
- Handles service discovery and load balancing

### 2. Authentication Service

- Manages user authentication and authorization
- Handles login, logout, and token validation
- Stores credentials in PostgreSQL database
- Supports JWT token-based authentication

### 3. User Service

- Manages user profiles and account information
- Handles user registration and profile updates
- Stores user data in PostgreSQL database
- Features:
  - User profile management
  - Address management
  - Payment method management
  - Wishlist functionality

### 4. Product Service

- Manages product catalog and inventory
- Supports product categories, search, and filtering
- Stores product data in PostgreSQL database
- Features:
  - Product CRUD operations
  - Category management
  - Product reviews and ratings
  - Image management

### 5. Cart Service

- Manages shopping cart functionality
- Handles add, update, remove operations
- Stores cart data in Redis for fast access
- Features:
  - Cart management
  - Coupon application
  - Price calculation

### 6. Checkout Service

- Handles order processing and payment
- Manages order lifecycle and status updates
- Stores order data in MongoDB
- Features:
  - Order creation and management
  - Payment processing
  - Order validation

### 7. Logger Service

- Centralized logging system
- Collects logs from all services
- Stores logs in MongoDB database

### 8. Listener Service

- Event-driven communication between services
- Listens for events on RabbitMQ
- Triggers appropriate actions based on events

### 9. Mail Service

- Handles email notifications
- Sends order confirmations, password resets, etc.
- Connects to SMTP server for email delivery

## Implementation Details

### Service Communication

- **gRPC Implementation**: All services communicate using gRPC with Protocol Buffers for efficient serialization
- **Event-Driven Architecture**: Asynchronous communication via RabbitMQ
- **API Gateway Pattern**: Broker service handles all client requests and routes them to appropriate services

### Security Features

- JWT-based authentication with refresh token mechanism
- Password hashing using bcrypt
- Role-based access control (RBAC)
- HTTPS for API endpoints

### Database Schema

Each service has its own dedicated database following the database-per-service pattern:

- PostgreSQL for structured data (users, products, authentication)
- MongoDB for unstructured data (logs, orders)
- Redis for high-performance caching (cart data)

### Performance Optimizations

- Connection pooling for database connections
- Redis caching for frequently accessed data
- Efficient data serialization with Protocol Buffers
- Load balancing at the API gateway level

## Installation Guide

### Prerequisites

- Go 1.19+
- Docker and Docker Compose
- Git

### Step 1: Clone the Repository

```bash
git clone https://github.com/yourusername/ecommerce-microservices.git
cd ecommerce-microservices
```

### Step 2: Environment Setup

Create environment files for each service:

```bash
# Copy example environment files
cp auth-service/.env.example auth-service/.env
cp broker-service/.env.example broker-service/.env
cp user-service/.env.example user-service/.env
cp product-service/.env.example product-service/.env
cp cart-service/.env.example cart-service/.env
cp checkout-service/.env.example checkout-service/.env
cp logger-service/.env.example logger-service/.env
cp mail-service/.env.example mail-service/.env
```

Edit each .env file to set appropriate values for your environment.

### Step 3: Build and Run with Docker Compose

```bash
# Build all services
docker-compose build

# Start all services
docker-compose up -d
```

### Step 4: Initialize Databases

```bash
# Run database migrations
docker-compose exec auth-service go run ./cmd/migrations
docker-compose exec user-service go run ./cmd/migrations
docker-compose exec product-service go run ./cmd/migrations

# Seed initial data (optional)
docker-compose exec product-service go run ./cmd/seed
```

## Usage Guide

### API Endpoints

#### Authentication Service (port 50051)

- `RegisterUser`: Register a new user
- `Login`: Authenticate user and generate JWT tokens
- `ValidateToken`: Validate JWT token
- `RefreshToken`: Get a new access token using refresh token

#### User Service (port 50052)

- `GetUser`: Get user profile by ID
- `UpdateUser`: Update user profile
- `AddAddress`: Add a new address for user
- `UpdateAddress`: Update existing address
- `DeleteAddress`: Delete user address
- `AddPaymentMethod`: Add payment method to user account
- `GetPaymentMethods`: List user payment methods

#### Product Service (port 50053/8082)

- gRPC endpoints (port 50053):

  - `GetProduct`: Get product by ID
  - `ListProducts`: List products with pagination
  - `CreateProduct`: Add new product
  - `UpdateProduct`: Update product details
  - `DeleteProduct`: Remove product

- REST endpoints (port 8082):
  - `GET /api/products`: List products with filtering options
  - `GET /api/products/:id`: Get product details
  - `GET /api/categories`: List product categories
  - `POST /api/products/:id/reviews`: Add product review

#### Cart Service (port 50054)

- `GetCart`: Get user cart contents
- `AddItem`: Add item to cart
- `UpdateItem`: Update cart item quantity
- `RemoveItem`: Remove item from cart
- `ApplyCoupon`: Apply discount coupon to cart
- `ClearCart`: Empty user cart

#### Checkout Service (port 50055)

- `CreateOrder`: Create new order from cart
- `GetOrder`: Get order details
- `UpdateOrderStatus`: Update order status
- `ProcessPayment`: Process order payment

### Example: Using the API with gRPC

#### Go Client Example

```go
package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "github.com/yourusername/ecommerce-microservices/product-service/proto"
)

func main() {
    conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewProductServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    response, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: 1})
    if err != nil {
        log.Fatalf("could not get product: %v", err)
    }
    log.Printf("Product: %s", response.Product.Name)
}
```

### Example: Using the REST API

```bash
# Get product list
curl -X GET http://localhost:8080/api/products

# Create a user account
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword","firstName":"John","lastName":"Doe"}'

# Login and get tokens
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

## Monitoring and Management

### Service Health Checks

- Each service exposes a health check endpoint at `/health`
- The broker service periodically checks the health of all services

### Logging

- Centralized logging through the logger service
- Log levels: DEBUG, INFO, WARNING, ERROR
- All logs can be queried through the logger service API

### Metrics

- Prometheus metrics available at `/metrics` endpoint on each service
- Grafana dashboards for visualizing service performance

## Data Models

### User Model

```go
type User struct {
    ID           string
    Email        string
    FirstName    string
    LastName     string
    DisplayName  string
    Phone        *string
    ProfileImage *string
    Addresses    []Address
    Role         string
    Status       string
    Active       bool
    PaymentMethods []PaymentMethod
    // Additional fields...
}
```

### Product Model

```go
type Product struct {
    ID            int
    Type          string
    Name          string
    Description   string
    Slug          string
    Price         float64
    ImageURL      string
    CategoryID    int
    CategorySlug  string
    StockQuantity int
    Images        []ProductImage
    Brand         string
    Tags          []string
    ReviewsAvg    ProductRating
    // Additional fields...
}
```

### Cart Model

```go
type Cart struct {
    ID             string
    UserID         string
    Items          []CartItem
    Totals         CartTotals
    CouponCode     string
    DiscountAmount float64
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

### Order Model

```go
type Order struct {
    ID           string
    UserID       string
    OrderNumber  string
    Status       string
    Items        []OrderItem
    BillingInfo  BillingInfo
    ShippingInfo ShippingInfo
    PaymentInfo  PaymentInfo
    Totals       OrderTotals
    CouponCode   string
    Notes        string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

## Communication Patterns

1. **Synchronous Communication (gRPC)**

   - Used for direct service-to-service communication
   - Implemented with Protocol Buffers for efficient serialization
   - Used by Authentication, User, Product, Cart, and Checkout services

2. **Asynchronous Communication (RabbitMQ)**
   - Used for event-driven communication
   - Implemented with RabbitMQ as message broker
   - Handled by Listener service for processing events

## Technologies Used

- **Backend**: Go (Golang)
- **API Communication**: gRPC, REST
- **Databases**:
  - PostgreSQL (Auth, User, Product services)
  - MongoDB (Checkout, Logger services)
  - Redis (Cart service)
- **Message Queue**: RabbitMQ
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **API Documentation**: Swagger/OpenAPI
- **Testing**: Go testing package, integration tests

## Development Workflow

### Local Development

1. Start only the required services:
   ```bash
   docker-compose up -d postgres redis mongodb rabbitmq
   ```
2. Run the service you're working on locally:
   ```bash
   cd product-service
   go run ./cmd/main.go
   ```

### Adding New Features

1. Implement the feature in the appropriate service
2. Add unit tests and integration tests
3. Update the proto file if needed
4. Regenerate gRPC code if proto files changed
5. Submit a pull request

### Generating gRPC Code

If you modify any Proto file, regenerate the Go code:

```bash
cd proto
protoc --go_out=. --go-grpc_out=. *.proto
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Service Ports

- Frontend: 3001
- Broker Service: 8080
- Authentication Service: 50051
- User Service: 50052
- Product Service: 50053/8082
- Cart Service: 50054
- Checkout Service: 50055
- Logger Service: 50056
- Mail Service: 50057
