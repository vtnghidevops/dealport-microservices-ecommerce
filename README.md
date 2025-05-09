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

## Getting Started

1. Clone the repository
2. Set up environment variables (see `.env.example`)
3. Run `docker-compose up -d` to start all services
4. Access the frontend at http://localhost:3001

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
