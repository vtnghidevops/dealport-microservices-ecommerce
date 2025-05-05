package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"checkout-service/internal/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OrderRepository handles order data persistence with MongoDB
type OrderRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewOrderRepository creates a new MongoDB order repository
func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		db:         db,
		collection: db.Collection("orders"),
	}
}

// CreateOrder creates a new order in the database
func (r *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// Debug: In ra thông tin đơn hàng trước khi lưu
	log.Printf("OrderRepository.CreateOrder received: %+v", order)
	log.Printf("BillingInfo in repo: %+v", order.BillingInfo)
	log.Printf("ShippingInfo in repo: %+v", order.ShippingInfo)
	log.Printf("PaymentInfo in repo: %+v", order.PaymentInfo)

	// Generate IDs and timestamps if not provided
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Generate order number if not provided
	if order.OrderNumber == "" {
		timestamp := time.Now().Unix()
		order.OrderNumber = fmt.Sprintf("ORD-%s-%d", uuid.New().String()[:8], timestamp)
	}

	// Set status if not provided
	if order.Status == "" {
		order.Status = "pending"
	}

	// Set timestamps
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	// Prepare a map for MongoDB insertion with proper types
	orderDoc := bson.M{
		"_id":           order.ID,
		"id":            order.ID, // Add this field to satisfy schema validation
		"user_id":       order.UserID,
		"order_number":  order.OrderNumber,
		"status":        order.Status,
		"items":         order.Items,
		"billing_info":  order.BillingInfo,
		"shipping_info": order.ShippingInfo,
		"payment_info":  order.PaymentInfo,
		"totals":        order.Totals,
		"coupon_code":   order.CouponCode,
		"notes":         order.Notes,
		"created_at":    order.CreatedAt, // Use time.Time object
		"updated_at":    order.UpdatedAt, // Use time.Time object
	}

	// Debug: BSON document trước khi insert
	log.Printf("BSON document to be inserted: %+v", orderDoc)

	// Insert order
	log.Printf("DEBUG: About to insert order with ID: %s, OrderNumber: %s", order.ID, order.OrderNumber)
	result, err := r.collection.InsertOne(ctx, orderDoc)
	if err != nil {
		log.Printf("ERROR: MongoDB InsertOne error: %v", err)
		return nil, errors.Join(domain.ErrDatabaseOperation, err)
	}

	log.Printf("SUCCESS: Order inserted with MongoDB ID: %v", result.InsertedID)

	// Verify the order was saved correctly by querying it back
	var savedOrder domain.Order
	filter := bson.M{"_id": order.ID}
	err = r.collection.FindOne(ctx, filter).Decode(&savedOrder)
	if err != nil {
		log.Printf("ERROR: Failed to verify order insertion: %v", err)
		if err == mongo.ErrNoDocuments {
			log.Printf("CRITICAL: Order with ID %s was not found after insertion!", order.ID)
		}
	} else {
		log.Printf("VERIFICATION: Successfully retrieved the order after insertion")
		log.Printf("VERIFICATION: Retrieved order ID: %s, OrderNumber: %s, UserID: %s",
			savedOrder.ID, savedOrder.OrderNumber, savedOrder.UserID)
		log.Printf("VERIFICATION: Order has %d items, total amount: %+v",
			len(savedOrder.Items), savedOrder.Totals)
	}

	// List all orders in collection to debug
	log.Printf("DEBUGGING: Listing all orders in collection...")
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("ERROR: Failed to list orders: %v", err)
	} else {
		defer cursor.Close(ctx)
		var orders []bson.M
		if err = cursor.All(ctx, &orders); err != nil {
			log.Printf("ERROR: Failed to decode orders: %v", err)
		} else {
			log.Printf("DEBUGGING: Found %d total orders in collection", len(orders))
			for i, o := range orders {
				if i < 5 { // Only show up to 5 orders to avoid log spam
					log.Printf("DEBUGGING: Order #%d: ID=%v, OrderNumber=%v",
						i+1, o["_id"], o["order_number"])
				}
			}
		}
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	var order domain.Order

	filter := bson.M{"_id": orderID}
	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrOrderNotFound
		}
		return nil, errors.Join(domain.ErrDatabaseOperation, err)
	}

	return &order, nil
}

// ListOrdersByUserID retrieves orders for a user with pagination
func (r *OrderRepository) ListOrdersByUserID(ctx context.Context, userID string, skip, limit int) ([]*domain.Order, int, error) {
	log.Printf("OrderRepository.ListOrdersByUserID: Starting fetch for userID=%s, skip=%d, limit=%d", userID, skip, limit)

	// Add detailed user ID logging
	log.Printf("DEBUG: UserID hex value: %x", []byte(userID))
	log.Printf("DEBUG: UserID length: %d", len(userID))

	filter := bson.M{"user_id": userID}
	log.Printf("MongoDB filter: %+v", filter)

	// Check for any orders with this user_id
	allOrdersFilter := bson.M{}
	distinctValues, distinctErr := r.collection.Distinct(ctx, "user_id", allOrdersFilter)
	if distinctErr == nil {
		log.Printf("DEBUG: All distinct user_ids in database: %v", distinctValues)
	}

	// Count total orders for this user
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("ERROR: CountDocuments failed: %v", err)
		return nil, 0, errors.Join(domain.ErrDatabaseOperation, err)
	}
	log.Printf("Total documents matching filter: %d", total)

	// Configure options for pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1}) // Most recent first

	// Execute the query
	log.Printf("Executing Find with options: skip=%d, limit=%d, sort=created_at:-1", skip, limit)
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("ERROR: Find operation failed: %v", err)
		return nil, 0, errors.Join(domain.ErrDatabaseOperation, err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var orders []*domain.Order
	log.Printf("Decoding results...")
	if err := cursor.All(ctx, &orders); err != nil {
		log.Printf("ERROR: Cursor.All failed during decoding: %v", err)
		return nil, 0, errors.Join(domain.ErrDatabaseOperation, err)
	}

	log.Printf("Successfully decoded %d orders", len(orders))
	// Debug for the first few orders
	for i, order := range orders {
		if i < 3 { // Limit debug output to at most 3 orders
			log.Printf("Order %d: ID=%s, UserID=%s, Total=%v, CreatedAt=%v",
				i+1, order.ID, order.UserID, order.Totals.Total, order.CreatedAt)
		}
	}

	return orders, int(total), nil
}

// UpdateOrderStatus updates the status of an order
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	filter := bson.M{"_id": orderID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(domain.ErrDatabaseOperation, err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

// UpdatePaymentInfo updates the payment information for an order
func (r *OrderRepository) UpdatePaymentInfo(ctx context.Context, orderID string, paymentInfo domain.PaymentInfo) error {
	filter := bson.M{"_id": orderID}
	update := bson.M{
		"$set": bson.M{
			"payment_info": paymentInfo,
			"updated_at":   time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Join(domain.ErrDatabaseOperation, err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

// GetOrderByPaymentTransactionID finds an order by its payment transaction ID
func (r *OrderRepository) GetOrderByPaymentTransactionID(ctx context.Context, transactionID string) (*domain.Order, error) {
	var order domain.Order

	filter := bson.M{"payment_info.transaction_id": transactionID}
	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrOrderNotFound
		}
		return nil, errors.Join(domain.ErrDatabaseOperation, err)
	}

	return &order, nil
}

// ListAllOrders retrieves all orders with pagination, without user filtering
// This is used by admin users only
func (r *OrderRepository) ListAllOrders(ctx context.Context, skip, limit int) ([]*domain.Order, int, error) {
	log.Printf("OrderRepository.ListAllOrders: Starting fetch for ALL orders, skip=%d, limit=%d", skip, limit)

	// Empty filter to get all orders
	filter := bson.M{}
	log.Printf("MongoDB filter for ALL orders: %+v", filter)

	// Count total orders
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("ERROR: CountDocuments failed for all orders: %v", err)
		return nil, 0, errors.Join(domain.ErrDatabaseOperation, err)
	}
	log.Printf("Total documents in orders collection: %d", total)

	// Configure options for pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1}) // Most recent first

	// Execute the query
	log.Printf("Executing Find for ALL orders with options: skip=%d, limit=%d, sort=created_at:-1", skip, limit)
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("ERROR: Find operation failed for all orders: %v", err)
		return nil, 0, errors.Join(domain.ErrDatabaseOperation, err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var orders []*domain.Order
	log.Printf("Decoding results for ALL orders...")
	if err := cursor.All(ctx, &orders); err != nil {
		log.Printf("ERROR: Cursor.All failed during decoding for all orders: %v", err)
		return nil, 0, errors.Join(domain.ErrDatabaseOperation, err)
	}

	log.Printf("Successfully decoded %d orders (from all users)", len(orders))
	// Debug for the first few orders
	for i, order := range orders {
		if i < 5 { // Limit debug output to at most 5 orders for admin view
			log.Printf("Order %d: ID=%s, UserID=%s, OrderNumber=%s, Status=%s, Total=%v, CreatedAt=%v",
				i+1, order.ID, order.UserID, order.OrderNumber, order.Status, order.Totals.Total, order.CreatedAt)
		}
	}

	return orders, int(total), nil
}

// GetUserTotalSpend calculates the total amount spent by a user across all their completed orders
func (r *OrderRepository) GetUserTotalSpend(ctx context.Context, userID string) (float64, error) {
	log.Printf("OrderRepository.GetUserTotalSpend: Calculating total spend for userID=%s", userID)

	// Only count orders with status "paid", "shipped", or "delivered"
	// These are considered completed orders where payment was successful
	filter := bson.M{
		"user_id": userID,
		"status": bson.M{
			"$in": []string{"paid", "shipped", "delivered"},
		},
	}

	// Use MongoDB aggregation pipeline to calculate the total
	pipeline := mongo.Pipeline{
		// Match orders for this user with valid statuses
		bson.D{{"$match", filter}},
		// Group by user_id and sum the totals
		bson.D{
			{"$group", bson.D{
				{"_id", "$user_id"},
				{"totalSpend", bson.D{
					{"$sum", "$totals.total"},
				}},
				{"orderCount", bson.D{
					{"$sum", 1},
				}},
			}},
		},
	}

	// Execute the aggregation
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("ERROR: Aggregation failed: %v", err)
		return 0, errors.Join(domain.ErrDatabaseOperation, err)
	}
	defer cursor.Close(ctx)

	// Process the result
	type result struct {
		ID         string  `bson:"_id"`
		TotalSpend float64 `bson:"totalSpend"`
		OrderCount int     `bson:"orderCount"`
	}

	var results []result
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("ERROR: Failed to decode aggregation results: %v", err)
		return 0, errors.Join(domain.ErrDatabaseOperation, err)
	}

	// If no results, the user has no orders or all orders have $0 total
	if len(results) == 0 {
		log.Printf("No completed orders found for user %s", userID)
		return 0, nil
	}

	// Return the total spend
	totalSpend := results[0].TotalSpend
	orderCount := results[0].OrderCount
	log.Printf("User %s has spent %.2f across %d completed orders", userID, totalSpend, orderCount)

	return totalSpend, nil
}

// GetUserOrderCount counts the number of orders placed by a user
func (r *OrderRepository) GetUserOrderCount(ctx context.Context, userID string) (int, error) {
	log.Printf("OrderRepository.GetUserOrderCount: Counting orders for userID=%s", userID)

	// Count all orders for this user regardless of status
	filter := bson.M{"user_id": userID}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("ERROR: CountDocuments failed: %v", err)
		return 0, errors.Join(domain.ErrDatabaseOperation, err)
	}

	log.Printf("User %s has placed %d orders", userID, count)
	return int(count), nil
}
