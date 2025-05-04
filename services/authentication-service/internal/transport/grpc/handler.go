package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"authentication-service/internal/domain"
	"authentication-service/internal/service"
	pb "authentication-service/proto/auth"

	"google.golang.org/grpc/metadata"
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
	fmt.Printf("CÓ VÔ HÀM REGISTER Ở HANDLER.GO")
	// Convert proto request to domain request
	domainReq := &domain.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Phone:     req.Phone,
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

// VerifyRegistration verifies an OTP and completes registration
func (h *AuthHandler) VerifyRegistration(ctx context.Context, req *pb.VerifyRegistrationRequest) (*pb.VerifyRegistrationResponse, error) {
	fmt.Printf("DEBUG gRPC VerifyRegistration: Received request for email: %s\n", req.Email)
	fmt.Printf("CÓ VÔ HÀM NÀY MÀ VERYFYREgistration ?????")
	// Convert proto request to domain request
	domainReq := &domain.VerifyRegistrationRequest{
		Email: req.Email,
		OTP:   req.Otp,
	}

	// Call service
	tokens, user, err := h.authService.VerifyRegistration(ctx, domainReq)
	if err != nil {
		fmt.Printf("DEBUG gRPC VerifyRegistration: Verification failed: %v\n", err)
		return &pb.VerifyRegistrationResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	fmt.Printf("DEBUG gRPC VerifyRegistration: Verification successful for user ID: %s\n", user.ID)
	return &pb.VerifyRegistrationResponse{
		Success:      true,
		UserId:       user.ID,
		Message:      "Registration verified successfully",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// RequestPasswordReset initiates password reset process
func (h *AuthHandler) RequestPasswordReset(ctx context.Context, req *pb.PasswordResetRequest) (*pb.PasswordResetResponse, error) {
	fmt.Printf("DEBUG gRPC RequestPasswordReset: Received request for email: %s\n", req.Email)

	// Call service
	err := h.authService.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		fmt.Printf("DEBUG gRPC RequestPasswordReset: Password reset request failed: %v\n", err)
		return &pb.PasswordResetResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	// Even if user doesn't exist, we return success for security reasons
	fmt.Printf("DEBUG gRPC RequestPasswordReset: Password reset request sent\n")
	return &pb.PasswordResetResponse{
		Success: true,
		Message: "If your email is registered, you will receive a password reset code",
	}, nil
}

// VerifyPasswordReset verifies password reset OTP
func (h *AuthHandler) VerifyPasswordReset(ctx context.Context, req *pb.VerifyPasswordResetRequest) (*pb.VerifyPasswordResetResponse, error) {
	fmt.Printf("DEBUG gRPC VerifyPasswordReset: Received request for email: %s\n", req.Email)

	// Convert proto request to domain request
	domainReq := &domain.VerifyPasswordResetRequest{
		Email: req.Email,
		OTP:   req.Otp,
	}

	// Call service
	resetToken, err := h.authService.VerifyPasswordReset(ctx, domainReq)
	if err != nil {
		fmt.Printf("DEBUG gRPC VerifyPasswordReset: Verification failed: %v\n", err)
		return &pb.VerifyPasswordResetResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	fmt.Printf("DEBUG gRPC VerifyPasswordReset: OTP verification successful\n")
	return &pb.VerifyPasswordResetResponse{
		Success:    true,
		Message:    "OTP verified successfully",
		ResetToken: resetToken,
	}, nil
}

// UpdatePassword updates user's password
func (h *AuthHandler) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	fmt.Printf("DEBUG gRPC UpdatePassword: Received password update request\n")

	// Extract email from metadata
	var email string

	// Try to get email from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// Get email from metadata
		emails := md.Get("email")
		if len(emails) > 0 {
			email = emails[0]
			fmt.Printf("DEBUG gRPC UpdatePassword: Extracted email from metadata: %s\n", email)
		}
	}

	// Create domain request
	domainReq := &domain.UpdatePasswordRequest{
		Token:           req.Token,
		Password:        req.Password,
		CurrentPassword: req.CurrentPassword,
		Email:           email,
	}

	fmt.Printf("DEBUG gRPC UpdatePassword: Processing request with email: %s\n", email)

	// Call service to update password
	err := h.authService.UpdatePassword(ctx, domainReq)
	if err != nil {
		fmt.Printf("DEBUG gRPC UpdatePassword: Password update failed: %v\n", err)
		return &pb.UpdatePasswordResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	fmt.Printf("DEBUG gRPC UpdatePassword: Password updated successfully\n")
	return &pb.UpdatePasswordResponse{
		Success: true,
		Message: "Password updated successfully",
	}, nil
}

// RequestOTP generates and sends a new OTP for a specific purpose
func (h *AuthHandler) RequestOTP(ctx context.Context, req *pb.OTPRequest) (*pb.OTPResponse, error) {
	fmt.Printf("DEBUG gRPC RequestOTP: Received request for email: %s, purpose: %s\n",
		req.Email, req.Purpose)

	// Validate purpose
	validPurpose := false
	switch req.Purpose {
	case "registration", "password_reset":
		validPurpose = true
	}

	if !validPurpose {
		errMsg := fmt.Sprintf("invalid purpose: %s, only registration and password_reset are allowed", req.Purpose)
		fmt.Printf("DEBUG gRPC RequestOTP: %s\n", errMsg)
		return &pb.OTPResponse{
			Success: false,
			Message: errMsg,
		}, status.Error(codes.InvalidArgument, errMsg)
	}

	// Call service
	err := h.authService.RequestOTP(ctx, req.Email, req.Purpose)
	if err != nil {
		fmt.Printf("DEBUG gRPC RequestOTP: OTP request failed: %v\n", err)
		return &pb.OTPResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	// Even if user doesn't exist, we return success for security reasons
	fmt.Printf("DEBUG gRPC RequestOTP: OTP request processed\n")
	return &pb.OTPResponse{
		Success: true,
		Message: "If your email is registered, you will receive an OTP code",
	}, nil
}

// VerifyOTP verifies an OTP code for a specific purpose
func (h *AuthHandler) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	fmt.Printf("DEBUG gRPC VerifyOTP: Received request for email: %s, purpose: %s\n",
		req.Email, req.Purpose)
	fmt.Println("GỌI VERIFY OTP KHÔNG GỌI VERIFY REGISTER")
	// Validate purpose
	validPurpose := false
	switch req.Purpose {
	case "registration", "password_reset":
		validPurpose = true
	}

	if !validPurpose {
		errMsg := fmt.Sprintf("invalid purpose: %s, only registration and password_reset are allowed", req.Purpose)
		fmt.Printf("DEBUG gRPC VerifyOTP: %s\n", errMsg)
		return &pb.VerifyOTPResponse{
			Success: false,
			Message: errMsg,
		}, status.Error(codes.InvalidArgument, errMsg)
	}

	// Convert proto request to domain request
	domainReq := &domain.VerifyOTPRequest{
		Email:   req.Email,
		OTP:     req.Otp,
		Purpose: req.Purpose,
	}

	// Call service
	token, err := h.authService.VerifyOTP(ctx, domainReq)
	if err != nil {
		fmt.Printf("DEBUG gRPC VerifyOTP: Verification failed: %v\n", err)
		return &pb.VerifyOTPResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	fmt.Printf("DEBUG gRPC VerifyOTP: OTP verification successful\n")
	return &pb.VerifyOTPResponse{
		Success: true,
		Message: "OTP verified successfully",
		Token:   token,
	}, nil
}

// CheckAccountExists verifies if an account exists with the given email
func (h *AuthHandler) CheckAccountExists(ctx context.Context, req *pb.CheckAccountExistsRequest) (*pb.CheckAccountExistsResponse, error) {
	fmt.Printf("DEBUG gRPC CheckAccountExists: Received request for email: %s\n", req.Email)

	// Call service
	exists, err := h.authService.CheckAccountExists(ctx, req.Email)
	if err != nil {
		fmt.Printf("DEBUG gRPC CheckAccountExists: Error checking account: %v\n", err)
		return &pb.CheckAccountExistsResponse{
			Exists:  false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var message string
	if exists {
		message = "Account found"
		fmt.Printf("DEBUG gRPC CheckAccountExists: Account exists for email: %s\n", req.Email)
	} else {
		message = "Account not found"
		fmt.Printf("DEBUG gRPC CheckAccountExists: Account does not exist for email: %s\n", req.Email)
	}

	return &pb.CheckAccountExistsResponse{
		Exists:  exists,
		Message: message,
	}, nil
}

// Logout handles logout requests
func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	var err error
	
	// Check if this is a logout from all devices
	if req.LogoutAllDevices {
		fmt.Printf("DEBUG gRPC Logout: Logging out from all devices for user %s\n", req.UserId)
		err = h.authService.LogoutFromAllDevices(ctx, req.UserId, req.Email)
	} else {
		fmt.Printf("DEBUG gRPC Logout: Logging out from current session for user %s\n", req.UserId)
		err = h.authService.Logout(ctx, req.UserId, req.Email)
	}
	
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Message: fmt.Sprintf("Logout failed: %v", err),
		}, status.Error(codes.Internal, "Failed to logout user")
	}

	// Create response
	response := &pb.LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
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
		Phone:     user.Phone,
	}
}
