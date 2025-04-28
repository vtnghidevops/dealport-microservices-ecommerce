package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"authentication-service/internal/domain"
	"authentication-service/internal/service"
	pb "authentication-service/proto/auth"
)

// AuthHandler handles gRPC requests for the authentication service
type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register creates a new user account
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	fmt.Printf("DEBUG gRPC Register: Received request for email: %s\n", req.Email)

	// Convert proto request to domain request
	domainReq := &domain.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
	}

	// Call service
	userID, err := h.authService.Register(ctx, domainReq)
	if err != nil {
		fmt.Printf("DEBUG gRPC Register: Registration failed: %v\n", err)
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	fmt.Printf("DEBUG gRPC Register: Registration successful for ID: %s\n", userID)
	return &pb.RegisterResponse{
		Success: true,
		UserId:  userID,
		Message: "User registered successfully",
	}, nil
}

// Login authenticates a user and returns tokens
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// fmt.Printf("DEBUG gRPC Login: Received request for email: %s, password length: %d\n",
	// 	req.Email, len(req.Password))
	// fmt.Printf("DEBUG gRPC Login: Password provided (FULL): %s\n", req.Password)

	// Convert proto request to domain request
	domainReq := &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call service
	tokens, user, err := h.authService.Login(ctx, domainReq)
	if err != nil {
		// fmt.Printf("DEBUG gRPC Login: Authentication failed: %v\n", err)
		return &pb.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	// fmt.Printf("DEBUG gRPC Login: Authentication successful for user ID: %s\n", user.ID)
	// fmt.Printf("DEBUG gRPC Login: Access token: %s\n", tokens.AccessToken)
	// fmt.Printf("DEBUG gRPC Login: Refresh token: %s\n", tokens.RefreshToken)

	// Create response
	response := &pb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		UserId:       user.ID,
		UserInfo:     mapUserToProto(user),
		Success:      true,
		Message:      "Login successful",
	}

	return response, nil
}

// Validate checks if a token is valid
func (h *AuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	// Call service
	metadata, err := h.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	// Create response
	response := &pb.ValidateResponse{
		Valid:  true,
		UserId: metadata.UserID,
		Claims: map[string]string{
			"email": metadata.Email,
			"role":  metadata.Role,
			"uuid":  metadata.UUID,
			"exp":   string(metadata.Exp),
		},
	}

	return response, nil
}

// RefreshToken refreshes an existing token
func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// Call service
	tokens, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return &pb.RefreshTokenResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	// Create response
	response := &pb.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Success:      true,
		Message:      "Token refreshed successfully",
	}

	return response, nil
}

// Helper function to map domain user to proto user
func mapUserToProto(user *domain.User) *pb.UserInfo {
	return &pb.UserInfo{
		UserId:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Role:      user.Role,
		Status:    user.Status,
		Active:    user.Active,
	}
}
