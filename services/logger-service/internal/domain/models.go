package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string                 `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string                 `bson:"name" json:"name"`
	Data      string                 `bson:"data" json:"data"`
	Level     string                 `bson:"level,omitempty" json:"level,omitempty"`
	Service   string                 `bson:"service,omitempty" json:"service,omitempty"`
	Action    string                 `bson:"action,omitempty" json:"action,omitempty"`
	UserID    string                 `bson:"user_id,omitempty" json:"user_id,omitempty"`
	RequestID string                 `bson:"request_id,omitempty" json:"request_id,omitempty"`
	Message   string                 `bson:"message,omitempty" json:"message,omitempty"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at" json:"updated_at"`
}

// ExtractFromData extracts structured data from the Data string if it's JSON
func (l *LogEntry) ExtractFromData() {
	if l.Data == "" {
		return
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(l.Data), &data)
	if err != nil {
		// Not valid JSON, just keep the data as is
		return
	}

	// Extract common fields if they exist
	if level, ok := data["level"].(string); ok && l.Level == "" {
		l.Level = level
	}
	if service, ok := data["service"].(string); ok && l.Service == "" {
		l.Service = service
	}
	if action, ok := data["action"].(string); ok && l.Action == "" {
		l.Action = action
	}
	if userID, ok := data["user_id"].(string); ok && l.UserID == "" {
		l.UserID = userID
	}
	if requestID, ok := data["request_id"].(string); ok && l.RequestID == "" {
		l.RequestID = requestID
	}
	if message, ok := data["message"].(string); ok && l.Message == "" {
		l.Message = message
	}
	if metadata, ok := data["metadata"].(map[string]interface{}); ok && l.Metadata == nil {
		l.Metadata = metadata
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	// Extract structured data if available
	entry.ExtractFromData()

	// Log the entry being inserted
	log.Printf("Inserting log entry: Name=%s, Level=%s, Service=%s, Action=%s, UserID=%s",
		entry.Name, entry.Level, entry.Service, entry.Action, entry.UserID)

	result, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		Level:     entry.Level,
		Service:   entry.Service,
		Action:    entry.Action,
		UserID:    entry.UserID,
		RequestID: entry.RequestID,
		Message:   entry.Message,
		Metadata:  entry.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	// Add detailed log with MongoDB ID
	log.Printf("📝 Log entry saved to MongoDB with ID: %v", result.InsertedID)
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// trans string to ObjId of mongo
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	// Extract structured data if available
	l.ExtractFromData()

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"level", l.Level},
				{"service", l.Service},
				{"action", l.Action},
				{"user_id", l.UserID},
				{"request_id", l.RequestID},
				{"message", l.Message},
				{"metadata", l.Metadata},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindByService returns all logs for a specific service
func (l *LogEntry) FindByService(service string) ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(
		context.TODO(),
		bson.D{{"service", service}},
		opts,
	)
	if err != nil {
		log.Println("Finding logs by service error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

// FindByUserAction returns all logs for a specific user action
func (l *LogEntry) FindByUserAction(userID, action string) ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	filter := bson.D{}
	if userID != "" {
		filter = append(filter, bson.E{"user_id", userID})
	}
	if action != "" {
		filter = append(filter, bson.E{"action", action})
	}

	cursor, err := collection.Find(
		context.TODO(),
		filter,
		opts,
	)
	if err != nil {
		log.Println("Finding logs by user action error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

// FindUserActivityLogs retrieves logs related to user activities with optional action type filter
func (l *LogEntry) FindUserActivityLogs(userID string, actionType string) ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	filter := bson.M{
		"user_id": userID,
		"name": bson.M{
			"$regex": "^log\\.INFO\\.user\\.",
		},
	}

	// If actionType is provided, add it to the filter
	if actionType != "" {
		filter["name"] = actionType
	}

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Println("Finding user activity logs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

// FindOrderLogs retrieves logs related to order activities with optional filters
func (l *LogEntry) FindOrderLogs(orderID, orderNumber, actionType string) ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Base filter for order logs
	filter := bson.M{
		"name": bson.M{
			"$regex": "^log\\.INFO\\.order\\.",
		},
	}

	// Apply additional filters if provided
	if orderID != "" {
		// Parse data field as JSON and check for order_id
		filter["$or"] = []bson.M{
			{"data": bson.M{"$regex": fmt.Sprintf("\"order_id\":\\s*\"%s\"", orderID)}},
		}
	}

	if orderNumber != "" {
		// Add orderNumber to the $or query if it doesn't exist
		orConditions, orExists := filter["$or"].([]bson.M)
		if orExists {
			filter["$or"] = append(orConditions, bson.M{
				"data": bson.M{"$regex": fmt.Sprintf("\"order_number\":\\s*\"%s\"", orderNumber)},
			})
		} else {
			filter["$or"] = []bson.M{
				{"data": bson.M{"$regex": fmt.Sprintf("\"order_number\":\\s*\"%s\"", orderNumber)}},
			}
		}
	}

	// If actionType is provided, add it to the filter
	if actionType != "" {
		// Replace the generic regex filter with specific action type
		filter["name"] = fmt.Sprintf("log.INFO.order.%s", actionType)
	}

	// Log the query for debugging
	queryJSON, _ := json.Marshal(filter)
	log.Printf("Order logs query: %s", string(queryJSON))

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Println("Finding order logs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

// LogUserActivity is a helper to log user activity events
func (l *LogEntry) LogUserActivity(action, userID, message string, metadata map[string]interface{}) error {
	// Parse the action to create the proper log name
	var logName string
	switch action {
	case "login_success":
		logName = "log.INFO.user.login_success"
	case "login_failed":
		logName = "log.INFO.user.login_failed"
	case "registered":
		logName = "log.INFO.user.registered"
	case "logout":
		logName = "log.INFO.user.logout"
	case "profile_updated":
		logName = "log.INFO.user.profile_updated"
	case "password_changed":
		logName = "log.INFO.user.password_changed"
	case "password_reset_requested":
		logName = "log.INFO.user.password_reset_requested"
	default:
		logName = "log.INFO.user." + action
	}

	// Create data object
	data := map[string]interface{}{
		"user_id":  userID,
		"action":   action,
		"message":  message,
		"level":    "INFO",
		"service":  "user-service",
		"metadata": metadata,
	}

	// Convert data to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create and insert log entry
	return l.Insert(LogEntry{
		Name:     logName,
		Data:     string(dataBytes),
		Level:    "INFO",
		Service:  "user-service",
		Action:   action,
		UserID:   userID,
		Message:  message,
		Metadata: metadata,
	})
}
