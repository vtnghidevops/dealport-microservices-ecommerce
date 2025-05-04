package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"logger-service/internal/domain"
	pb "logger-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	gRpcPort = "50001"
)

type Config struct {
	Models domain.Models
}

type LogServer struct {
	pb.UnimplementedLogServiceServer
	Models domain.Models
}

// NewGrpcServer tạo và cấu hình máy chủ gRPC mới
func (app *Config) NewGrpcServer() *grpc.Server {
	// Tạo máy chủ gRPC mới
	server := grpc.NewServer()

	// Đăng ký dịch vụ
	pb.RegisterLogServiceServer(server, &LogServer{
		Models: app.Models,
	})

	// Đăng ký reflection service
	reflection.Register(server)

	return server
}

// WriteLog writes a log entry to MongoDB
func (l *LogServer) WriteLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	// Extract data from the request
	input := req.GetLogEntry()

	// Create a log entry for MongoDB
	logEntry := domain.LogEntry{
		Name:      input.Name,
		Data:      input.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Level:     input.Level,
		Service:   input.Service,
		Action:    input.Action,
		UserID:    input.UserId,
		RequestID: input.RequestId,
		Message:   input.Message,
	}

	// Insert the log entry
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		log.Printf("Error inserting log entry: %v", err)
		return &pb.LogResponse{
			Result: "error",
		}, err
	}

	// Return success response
	res := &pb.LogResponse{Result: "logged"}
	return res, nil
}

// GetUserActivityLogs retrieves logs for user activities
func (l *LogServer) GetUserActivityLogs(ctx context.Context, req *pb.UserActivityLogRequest) (*pb.UserActivityLogResponse, error) {
	userID := req.GetUserId()
	actionType := req.GetActionType()

	// Get logs from database
	logs, err := l.Models.LogEntry.FindUserActivityLogs(userID, actionType)
	if err != nil {
		log.Printf("Error retrieving user activity logs: %v", err)
		return &pb.UserActivityLogResponse{}, err
	}

	// Convert logs to protobuf format
	var responseItems []*pb.Log
	for _, logEntry := range logs {
		pbLog := &pb.Log{
			Name:      logEntry.Name,
			Data:      logEntry.Data,
			Level:     logEntry.Level,
			Service:   logEntry.Service,
			Action:    logEntry.Action,
			UserId:    logEntry.UserID,
			RequestId: logEntry.RequestID,
			Message:   logEntry.Message,
		}
		responseItems = append(responseItems, pbLog)
	}

	return &pb.UserActivityLogResponse{
		Logs: responseItems,
	}, nil
}

// GRPCListen starts the gRPC server
func (app *Config) GRPCListen() {
	log.Println("Starting gRPC server on port", gRpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Create a new gRPC server
	server := app.NewGrpcServer()

	log.Printf("gRPC Server starting on port %s", gRpcPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
