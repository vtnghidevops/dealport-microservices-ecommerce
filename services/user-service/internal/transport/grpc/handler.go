package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"user-service/internal/domain"
	"user-service/internal/service"
	pb "user-service/proto/user"
)

// UserHandler handles gRPC requests for the user service
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get user: %v", err)
	}

	return &pb.UserResponse{
		User:    convertDomainUserToProto(user),
		Message: "User retrieved successfully",
	}, nil
}

// GetUsers retrieves a list of users with optional filtering
func (h *UserHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	// Convert proto filter to domain filter
	filter := &domain.UserFilter{
		Page:      int(req.Page),
		Limit:     int(req.Limit),
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Role:      req.Role,
	}

	users, total, err := h.userService.GetUsers(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	// Convert domain users to proto users
	var protoUsers []*pb.User
	for _, user := range users {
		protoUsers = append(protoUsers, convertDomainUserToProto(user))
	}

	return &pb.GetUsersResponse{
		Users:   protoUsers,
		Total:   int32(total),
		Page:    int32(filter.Page),
		Limit:   int32(filter.Limit),
		Message: "Users retrieved successfully",
	}, nil
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	// Convert proto request to domain request
	domainReq := &domain.CreateUserRequest{
		Email:        req.Email,
		Password:     req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		DisplayName:  req.DisplayName,
		Phone:        req.Phone,
		ProfileImage: req.ProfileImage,
		Role:         req.Role,
	}

	// Convert addresses if any
	if len(req.Addresses) > 0 {
		domainReq.Addresses = make([]domain.Address, len(req.Addresses))
		for i, addr := range req.Addresses {
			domainReq.Addresses[i] = convertProtoAddressToDomain(addr)
		}
	}

	// Create user
	user, err := h.userService.CreateUser(ctx, domainReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.UserResponse{
		User:    convertDomainUserToProto(user),
		Message: "User created successfully",
	}, nil
}

// UpdateUser updates user information
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	// Convert proto request to domain request
	domainReq := &domain.UpdateUserRequest{
		ID:           req.Id,
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		DisplayName:  req.DisplayName,
		Phone:        req.Phone,
		ProfileImage: req.ProfileImage,
		Role:         req.Role,
		Active:       req.Active,
	}

	// Convert addresses if any
	if len(req.Addresses) > 0 {
		domainReq.Addresses = make([]domain.Address, len(req.Addresses))
		for i, addr := range req.Addresses {
			domainReq.Addresses[i] = convertProtoAddressToDomain(addr)
		}
	}

	// Update user
	user, err := h.userService.UpdateUser(ctx, domainReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UserResponse{
		User:    convertDomainUserToProto(user),
		Message: "User updated successfully",
	}, nil
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := h.userService.DeleteUser(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

// SearchUsers searches for users based on criteria
func (h *UserHandler) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.GetUsersResponse, error) {
	// Convert proto search params to domain search params
	params := &domain.SearchParams{
		Query: req.Query,
		Field: req.Field,
		Page:  int(req.Page),
		Limit: int(req.Limit),
	}

	users, total, err := h.userService.SearchUsers(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search users: %v", err)
	}

	// Convert domain users to proto users
	var protoUsers []*pb.User
	for _, user := range users {
		protoUsers = append(protoUsers, convertDomainUserToProto(user))
	}

	return &pb.GetUsersResponse{
		Users:   protoUsers,
		Total:   int32(total),
		Page:    int32(params.Page),
		Limit:   int32(params.Limit),
		Message: "Users found successfully",
	}, nil
}

// Helper methods to convert between domain and proto models

// convertDomainUserToProto converts a domain user to a proto user
func convertDomainUserToProto(user *domain.User) *pb.User {
	protoUser := &pb.User{
		Id:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		DisplayName:  user.DisplayName,
		Phone:        getStringValue(user.Phone),
		ProfileImage: getStringValue(user.ProfileImage),
		Role:         user.Role,
		Active:       user.Active,
	}

	// Convert created_at and updated_at to protobuf timestamps
	if !user.CreatedAt.IsZero() {
		protoUser.CreatedAt = timestamppb.New(user.CreatedAt).AsTime().Format(time.RFC3339)
	}

	if !user.UpdatedAt.IsZero() {
		protoUser.UpdatedAt = timestamppb.New(user.UpdatedAt).AsTime().Format(time.RFC3339)
	}

	// Convert addresses if any
	if len(user.Addresses) > 0 {
		protoUser.Addresses = make([]*pb.Address, len(user.Addresses))
		for i, addr := range user.Addresses {
			protoUser.Addresses[i] = convertDomainAddressToProto(&addr)
		}
	}

	return protoUser
}

// convertProtoAddressToDomain converts a proto address to a domain address
func convertProtoAddressToDomain(addr *pb.Address) domain.Address {
	return domain.Address{
		ID:          addr.Id,
		Line1:       addr.Line1,
		City:        addr.City,
		State:       addr.State,
		PostalCode:  addr.PostalCode,
		Country:     addr.Country,
		IsDefault:   addr.IsDefault,
		AddressType: addr.AddressType,
	}
}

// convertDomainAddressToProto converts a domain address to a proto address
func convertDomainAddressToProto(addr *domain.Address) *pb.Address {
	return &pb.Address{
		Id:          addr.ID,
		Line1:       addr.Line1,
		City:        addr.City,
		State:       addr.State,
		PostalCode:  addr.PostalCode,
		Country:     addr.Country,
		IsDefault:   addr.IsDefault,
		AddressType: addr.AddressType,
	}
}

// Add helper function for handling nil pointers at the end of the file
// getStringValue safely gets string value from a pointer that might be nil
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
