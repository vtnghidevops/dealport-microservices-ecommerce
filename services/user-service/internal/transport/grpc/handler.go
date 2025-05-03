package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	logger      *log.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
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

// ProcessEvent handles events from other services
func (h *UserHandler) ProcessEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	h.logger.Printf("Received event: %s from %s", req.EventName, req.Source)

	switch req.EventName {
	case "user.registered":
		return h.handleUserRegisteredEvent(ctx, req)
	case "user.password_changed":
		return h.handlePasswordChangedEvent(ctx, req)
	default:
		h.logger.Printf("Unknown event type: %s", req.EventName)
		return &pb.EventResponse{
			Success: false,
			Message: "Unknown event type",
		}, nil
	}
}

// handleUserRegisteredEvent processes user.registered events
func (h *UserHandler) handleUserRegisteredEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	// Log the event data for debugging
	h.logger.Printf("DEBUG USER-SERVICE: Processing user.registered event from %s. Event ID: %s", req.Source, req.EventId)
	h.logger.Printf("DEBUG USER-SERVICE: User registration event data: %s", req.EventData)

	// Parse the complete event data using a struct that matches the UserRegisteredData from auth service
	var userData struct {
		ID           string    `json:"id"`
		Email        string    `json:"email"`
		Username     string    `json:"username"`
		FirstName    string    `json:"first_name"`
		LastName     string    `json:"last_name"`
		DisplayName  string    `json:"display_name,omitempty"`
		Phone        string    `json:"phone,omitempty"`
		ProfileImage string    `json:"profile_image,omitempty"`
		Role         string    `json:"role"`
		Status       string    `json:"status"`
		Active       bool      `json:"active"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Gender       string    `json:"gender,omitempty"`
	}

	if err := json.Unmarshal([]byte(req.EventData), &userData); err != nil {
		h.logger.Printf("ERROR USER-SERVICE: Failed to parse user registration data: %v. Raw data: %s", err, req.EventData)

		// Try alternative approach with a generic map
		var dataMap map[string]interface{}
		if mapErr := json.Unmarshal([]byte(req.EventData), &dataMap); mapErr == nil {
			h.logger.Printf("DEBUG USER-SERVICE: Parsed as generic map: %+v", dataMap)
			// Try to extract essential fields from the map
			if id, ok := dataMap["id"].(string); ok {
				h.logger.Printf("DEBUG USER-SERVICE: Found ID in map: %s", id)
			}
			if email, ok := dataMap["email"].(string); ok {
				h.logger.Printf("DEBUG USER-SERVICE: Found Email in map: %s", email)
			}
		} else {
			h.logger.Printf("ERROR USER-SERVICE: Even generic map parsing failed: %v", mapErr)
		}

		return &pb.EventResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing user data: %v", err),
		}, nil
	}

	// Print parsed user data for debugging
	h.logger.Printf("DEBUG USER-SERVICE: Successfully parsed user data: ID=%s, Email=%s, Name=%s %s, Username=%s, Role=%s, Status=%s, Active=%v",
		userData.ID, userData.Email, userData.FirstName, userData.LastName, userData.Username, userData.Role, userData.Status, userData.Active)

	// Check if user already exists (to handle potential duplicate events)
	existingUser, err := h.userService.GetUserByID(ctx, userData.ID)
	if err == nil {
		h.logger.Printf("DEBUG USER-SERVICE: User %s already exists in database, updating instead of creating", userData.ID)

		// Update the existing user with new data
		updateUserReq := &domain.UpdateUserRequest{
			ID:           existingUser.ID,
			Email:        userData.Email,
			FirstName:    userData.FirstName,
			LastName:     userData.LastName,
			Username:     userData.Username,
			DisplayName:  userData.DisplayName,
			Phone:        userData.Phone,
			ProfileImage: userData.ProfileImage,
			Role:         userData.Role,
			Status:       userData.Status,
			Active:       userData.Active,
		}

		h.logger.Printf("DEBUG USER-SERVICE: Updating existing user with request: %+v", updateUserReq)

		// Update the user
		user, err := h.userService.UpdateUser(ctx, updateUserReq)
		if err != nil {
			h.logger.Printf("ERROR USER-SERVICE: Failed to update existing user: %v", err)
			return &pb.EventResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to update existing user: %v", err),
			}, nil
		}

		h.logger.Printf("DEBUG USER-SERVICE: User %s updated successfully", user.ID)
		return &pb.EventResponse{
			Success: true,
			Message: "User profile updated successfully",
		}, nil
	} else {
		h.logger.Printf("DEBUG USER-SERVICE: User %s not found in database, creating new user profile", userData.ID)
	}

	// Create user profile in user service
	createUserReq := &domain.CreateUserRequest{
		ID:           userData.ID,
		Email:        userData.Email,
		FirstName:    userData.FirstName,
		LastName:     userData.LastName,
		Username:     userData.Username,
		DisplayName:  userData.DisplayName,
		Phone:        userData.Phone,
		ProfileImage: userData.ProfileImage,
		Role:         userData.Role,
		// Note: Password is not included as auth is handled by auth service
		// Status and Active are not part of CreateUserRequest, they'll be set during service creation
	}

	h.logger.Printf("DEBUG USER-SERVICE: Creating new user with request: %+v", createUserReq)

	// Create the user
	user, err := h.userService.CreateUser(ctx, createUserReq)
	if err != nil {
		h.logger.Printf("ERROR USER-SERVICE: Failed to create user profile: %v", err)
		return &pb.EventResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create user profile: %v", err),
		}, nil
	}

	h.logger.Printf("DEBUG USER-SERVICE: User profile created successfully - ID: %s, Email: %s", user.ID, user.Email)
	return &pb.EventResponse{
		Success: true,
		Message: "User profile created successfully",
	}, nil
}

// handlePasswordChangedEvent processes user.password_changed events
func (h *UserHandler) handlePasswordChangedEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	// Parse the password changed data
	var passwordData struct {
		Email     string `json:"email"`
		ChangedAt string `json:"changed_at"`
	}

	if err := json.Unmarshal([]byte(req.EventData), &passwordData); err != nil {
		h.logger.Printf("Error parsing password data: %v", err)
		return &pb.EventResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing password data: %v", err),
		}, nil
	}

	// Get user by email
	user, err := h.userService.GetUserByEmail(ctx, passwordData.Email)
	if err != nil {
		h.logger.Printf("Failed to get user by email: %v", err)
		return &pb.EventResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get user by email: %v", err),
		}, nil
	}

	// Create update request
	updateUserReq := &domain.UpdateUserRequest{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Role:      user.Role,
		Status:    user.Status,
		Active:    user.Active,
	}

	// Update user
	_, err = h.userService.UpdateUser(ctx, updateUserReq)
	if err != nil {
		h.logger.Printf("Failed to update user: %v", err)
		return &pb.EventResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update user: %v", err),
		}, nil
	}

	h.logger.Printf("Password changed timestamp updated for user %s", user.ID)
	return &pb.EventResponse{
		Success: true,
		Message: "Password changed timestamp updated",
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
