package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "listener-service/proto/logs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoggerClient struct {
	conn   *grpc.ClientConn
	client pb.LogServiceClient
}

func NewLoggerClient(loggerServiceAddress string) (*LoggerClient, error) {
	// Set default address if not provided
	if loggerServiceAddress == "" {
		// Check environment variable first
		loggerServiceAddress = os.Getenv("LOGGER_SERVICE_HOST")
		if loggerServiceAddress == "" {
			loggerServiceAddress = "logger-service:50056"
		}
	}

	log.Printf("Connecting to logger service at: %s", loggerServiceAddress)

	// Create a connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		loggerServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf("Failed to connect to logger service at %s: %v", loggerServiceAddress, err)
		return nil, fmt.Errorf("failed to connect to logger service at %s: %w", loggerServiceAddress, err)
	}

	log.Printf("Successfully connected to logger service at %s", loggerServiceAddress)
	client := pb.NewLogServiceClient(conn)

	return &LoggerClient{
		conn:   conn,
		client: client,
	}, nil
}

func (lc *LoggerClient) Close() error {
	if lc.conn != nil {
		return lc.conn.Close()
	}
	return nil
}

// WriteLog sends a log entry to the logger service
func (lc *LoggerClient) WriteLog(ctx context.Context, name, data string) error {
	// Create a timeout context for the RPC call
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Printf("Sending log to logger service: %s", name)

	// Create request
	req := &pb.LogRequest{
		LogEntry: &pb.Log{
			Name:    name,
			Data:    data,
			Level:   "INFO",
			Service: "listener-service",
		},
	}

	// Log the request details for debugging
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("DEBUG: Request payload: %s", string(reqJSON))

	// Send the log via gRPC
	response, err := lc.client.WriteLog(timeoutCtx, req)

	if err != nil {
		log.Printf("Error calling WriteLog: %v", err)
		return fmt.Errorf("error calling WriteLog: %w", err)
	}

	log.Printf("Log entry successfully sent, result: %s", response.Result)
	return nil
}

// LogUserActivity logs a standard user activity event
func (lc *LoggerClient) LogUserActivity(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error {
	// Create the data payload
	data := map[string]interface{}{
		"level":     "INFO",
		"service":   "user-service",
		"action":    action,
		"user_id":   userID,
		"message":   message,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Add metadata if provided
	if len(metadata) > 0 {
		data["metadata"] = metadata
	}

	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling log data: %w", err)
	}

	// Determine log name based on action
	logName := fmt.Sprintf("log.INFO.user.%s", action)

	log.Printf("Logging user activity: %s for user ID: %s", action, userID)

	// Send to logger service
	return lc.WriteLog(ctx, logName, string(jsonData))
}
