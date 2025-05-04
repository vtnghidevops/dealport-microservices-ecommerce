package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "user-service/proto/logs"
)

// LoggerClient is a client for the logger service
type LoggerClient struct {
	client pb.LogServiceClient
	conn   *grpc.ClientConn
}

// NewLoggerClient creates a new logger client
func NewLoggerClient(loggerHost string) (*LoggerClient, error) {
	// Set up connection to the logger service
	conn, err := grpc.Dial(
		loggerHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to logger service: %w", err)
	}

	// Create client
	client := pb.NewLogServiceClient(conn)

	return &LoggerClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (lc *LoggerClient) Close() error {
	if lc.conn != nil {
		return lc.conn.Close()
	}
	return nil
}

// LogUserActivity logs user activity events
func (lc *LoggerClient) LogUserActivity(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error {
	// Create log name based on the action
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

	// Create the log data
	logData := map[string]interface{}{
		"user_id":   userID,
		"action":    action,
		"message":   message,
		"level":     "INFO",
		"service":   "user-service",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Add metadata if provided
	if metadata != nil {
		logData["metadata"] = metadata
	}

	// Convert to JSON
	jsonData, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("failed to marshal log data: %w", err)
	}

	// Create the log request
	req := &pb.LogRequest{
		LogEntry: &pb.Log{
			Name:    logName,
			Data:    string(jsonData),
			UserId:  userID,
			Service: "user-service",
			Action:  action,
			Message: message,
			Level:   "INFO",
		},
	}

	// Send log to the logger service
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = lc.client.WriteLog(ctx, req)
	if err != nil {
		log.Printf("ERROR: Failed to log user activity: %v", err)
		return fmt.Errorf("failed to log user activity: %w", err)
	}

	return nil
}

// GetUserActivityLogs retrieves logs for user activities
func (lc *LoggerClient) GetUserActivityLogs(ctx context.Context, userID string, actionType string) ([]*pb.Log, error) {
	// Create the request
	req := &pb.UserActivityLogRequest{
		UserId:     userID,
		ActionType: actionType,
	}

	// Send request to the logger service
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := lc.client.GetUserActivityLogs(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user activity logs: %w", err)
	}

	return resp.Logs, nil
}
