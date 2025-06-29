# Checkout Service Testing Documentation

This directory contains tests for the Checkout Service. The tests are organized into several types and categories to ensure comprehensive test coverage.

## Test Structure

```
tests/
├── README.md               # This file
├── mocks/                 # Mock implementations for external dependencies
│   ├── domain/            # Domain models for testing
│   ├── payment_service_mock.go
│   ├── repository_mock.go
│   └── event_emitter_mock.go
├── unit/                  # Unit tests for individual components
│   └── order_service_test.go
├── integration/           # Integration tests that test multiple components together
│   └── api_test.go
├── service/               # Service layer tests
│   ├── mock_adapter.go
│   ├── mock_repository.go
│   ├── mock_event_emitter.go
│   ├── order_service_test.go
│   └── security_test.go
├── security/              # Security-specific tests
│   └── security_test.go
└── helpers/               # Helper utilities for testing
    └── test_utils.go
```

## Test Types

### Unit Tests

Unit tests focus on testing individual functions or methods in isolation. These tests typically mock all dependencies of the unit under test.

Example:
```go
func TestCreateOrder(t *testing.T) {
    mockRepo := mocks.NewMockOrderRepository()
    mockEmitter := mocks.NewMockEventEmitter()
    orderService := service.NewOrderServiceWithMocks(mockRepo, mockEmitter)
    
    // Test the service
    order := &domain.Order{...}
    createdOrder, err := orderService.CreateOrder(context.Background(), order)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, createdOrder)
}
```

### Integration Tests

Integration tests verify that different components of the application work correctly together. These tests may use actual database connections or may mock external services.

Example:
```go
func TestOrderCreationFlow(t *testing.T) {
    // Set up test environment with real or containerized dependencies
    
    // Test the entire flow from API to database
    req := createOrderRequest{...}
    resp, err := http.Post("/api/orders", req)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

### Security Tests

Security tests focus on ensuring the application has proper authentication, authorization, and data validation.

Example:
```go
func TestUnauthorizedOrderAccess(t *testing.T) {
    // Attempt to access an order without proper authorization
    resp, err := http.Get("/api/orders/123", unauthorizedHeader)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 403, resp.StatusCode)
}
```

## Running Tests

To run all tests:

```bash
go test ./...
```

To run specific test types:

```bash
go test ./tests/unit/...  # Run only unit tests
go test ./tests/integration/...  # Run only integration tests
```

To run tests with verbose output:

```bash
go test -v ./...
```

## Writing New Tests

When writing new tests:

1. Choose the appropriate test category (unit, integration, etc.)
2. Follow existing naming conventions
3. Use mocks from the `mocks` directory or create new ones as needed
4. Ensure proper cleanup of test resources
5. Write clear assertions that document expected behavior

## Mocks

The `mocks` directory contains mock implementations used across tests:

- `payment_service_mock.go` - Mocks the payment service interface
- `repository_mock.go` - Mocks database operations
- `event_emitter_mock.go` - Mocks event emission

When creating new mocks, use the testify/mock package for consistency:

```go
type MockExampleService struct {
    mock.Mock
}

func (m *MockExampleService) DoSomething(param string) error {
    args := m.Called(param)
    return args.Error(0)
}
``` 