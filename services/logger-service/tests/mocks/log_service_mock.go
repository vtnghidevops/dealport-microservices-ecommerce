package mocks

import (
	"context"

	pb "logger-service/proto"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockLogServiceClient mô phỏng LogServiceClient từ proto definition
type MockLogServiceClient struct {
	mock.Mock
}

// WriteLog mô phỏng phương thức WriteLog
func (m *MockLogServiceClient) WriteLog(ctx context.Context, in *pb.LogRequest, opts ...grpc.CallOption) (*pb.LogResponse, error) {
	args := m.Called(ctx, in, opts)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*pb.LogResponse), args.Error(1)
}

// GetUserActivityLogs mô phỏng phương thức GetUserActivityLogs
func (m *MockLogServiceClient) GetUserActivityLogs(ctx context.Context, in *pb.UserActivityLogRequest, opts ...grpc.CallOption) (*pb.UserActivityLogResponse, error) {
	args := m.Called(ctx, in, opts)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*pb.UserActivityLogResponse), args.Error(1)
}

// GetOrderLogs mô phỏng phương thức GetOrderLogs
func (m *MockLogServiceClient) GetOrderLogs(ctx context.Context, in *pb.OrderLogRequest, opts ...grpc.CallOption) (*pb.OrderLogResponse, error) {
	args := m.Called(ctx, in, opts)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*pb.OrderLogResponse), args.Error(1)
}

// MockLogServiceServer mô phỏng LogServiceServer 
type MockLogServiceServer struct {
	mock.Mock
	pb.UnimplementedLogServiceServer
}

// WriteLog mô phỏng phương thức WriteLog của server
func (m *MockLogServiceServer) WriteLog(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	args := m.Called(ctx, in)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*pb.LogResponse), args.Error(1)
}

// GetUserActivityLogs mô phỏng phương thức GetUserActivityLogs của server
func (m *MockLogServiceServer) GetUserActivityLogs(ctx context.Context, in *pb.UserActivityLogRequest) (*pb.UserActivityLogResponse, error) {
	args := m.Called(ctx, in)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*pb.UserActivityLogResponse), args.Error(1)
}

// GetOrderLogs mô phỏng phương thức GetOrderLogs của server
func (m *MockLogServiceServer) GetOrderLogs(ctx context.Context, in *pb.OrderLogRequest) (*pb.OrderLogResponse, error) {
	args := m.Called(ctx, in)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*pb.OrderLogResponse), args.Error(1)
} 