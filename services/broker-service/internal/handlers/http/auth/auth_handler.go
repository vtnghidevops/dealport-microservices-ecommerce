package auth

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"broker-service/internal/util"
	authpb "broker-service/proto/auth"

	"google.golang.org/grpc/metadata"
)

type Config struct {
	AuthClient authpb.AuthServiceClient
}

// Register handles user registration
func (c *Config) Register(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call auth service via gRPC
	res, err := c.AuthClient.Register(r.Context(), &authpb.RegisterRequest{
		Email:     requestPayload.Email,
		Password:  requestPayload.Password,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Username:  requestPayload.Username,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Return response
	responseData := util.JsonResponse{
		Error:   false,
		Message: res.Message,
		Data: map[string]interface{}{
			"user_id": res.UserId,
			"success": res.Success,
		},
	}

	util.WriteJSON(w, http.StatusCreated, responseData)
}

// VerifyRegistration handles OTP verification for new user registration
func (c *Config) VerifyRegistration(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call auth service via gRPC
	fmt.Printf("DEBUG BROKER VerifyRegistration: Verifying OTP for email: %s\n", requestPayload.Email)
	res, err := c.AuthClient.VerifyRegistration(r.Context(), &authpb.VerifyRegistrationRequest{
		Email: requestPayload.Email,
		Otp:   requestPayload.OTP,
	})

	if err != nil {
		fmt.Printf("DEBUG BROKER VerifyRegistration: Error from auth service: %v\n", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if !res.Success {
		fmt.Printf("DEBUG BROKER VerifyRegistration: Auth service returned unsuccessful response: %s\n", res.Message)
		util.ErrorJSON(w, errors.New(res.Message), http.StatusBadRequest)
		return
	}

	// Return response
	responseData := util.JsonResponse{
		Error:   false,
		Message: res.Message,
		Data: map[string]interface{}{
			"user_id":       res.UserId,
			"success":       res.Success,
			"access_token":  res.AccessToken,
			"refresh_token": res.RefreshToken,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// Login handles user authentication
func (c *Config) Login(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call auth service via gRPC
	fmt.Printf("DEBUG BROKER Login: Attempting login for email: %s\n", requestPayload.Email)
	res, err := c.AuthClient.Login(r.Context(), &authpb.LoginRequest{
		Email:    requestPayload.Email,
		Password: requestPayload.Password,
	})

	if err != nil {
		fmt.Printf("DEBUG BROKER Login: Error from auth service: %v\n", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if !res.Success {
		fmt.Printf("DEBUG BROKER Login: Auth service returned unsuccessful response: %s\n", res.Message)
		util.ErrorJSON(w, errors.New(res.Message), http.StatusUnauthorized)
		return
	}

	// Log token details
	fmt.Printf("DEBUG BROKER Login: Successful login for user: %s\n", res.UserId)
	//fmt.Printf("DEBUG BROKER Login: FULL ACCESS TOKEN: %s\n", res.AccessToken)
	//fmt.Printf("DEBUG BROKER Login: FULL REFRESH TOKEN: %s\n", res.RefreshToken)
	//fmt.Printf("DEBUG BROKER Login: Access token length: %d\n", len(res.AccessToken))
	//fmt.Printf("DEBUG BROKER Login: Refresh token length: %d\n", len(res.RefreshToken))

	// Check token structure
	// parts := strings.Split(res.RefreshToken, ".")
	// if len(parts) == 3 {
	// 	fmt.Printf("DEBUG BROKER Login: Refresh token header: %s\n", parts[0])
	// 	// fmt.Printf("DEBUG BROKER Login: Refresh token payload: %s\n", parts[1])
	// 	// fmt.Printf("DEBUG BROKER Login: Refresh token signature: %s\n", parts[2])
	// }

	// Prepare response
	responseData := util.JsonResponse{
		Error:   false,
		Message: "Logged in successfully",
		Data: map[string]interface{}{
			"access_token":  res.AccessToken,
			"refresh_token": res.RefreshToken,
			"user_id":       res.UserId,
			"user_info":     res.UserInfo,
			"success":       res.Success,
			"message":       res.Message,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// ValidateToken validates an authentication token
func (c *Config) ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		util.ErrorJSON(w, errors.New("authorization header missing"), http.StatusUnauthorized)
		return
	}

	// The header format should be "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		util.ErrorJSON(w, errors.New("invalid authorization header format"), http.StatusUnauthorized)
		return
	}

	token := parts[1]

	// Call auth service to validate token
	res, err := c.AuthClient.Validate(r.Context(), &authpb.ValidateRequest{
		Token: token,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if !res.Valid {
		util.ErrorJSON(w, errors.New("invalid token"), http.StatusUnauthorized)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Token is valid",
		Data: map[string]interface{}{
			"user_id": res.UserId,
			"claims":  res.Claims,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// RefreshToken handles token refresh
func (c *Config) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Debug logs - print full token
	token := requestPayload.RefreshToken
	fmt.Printf("DEBUG BROKER RefreshToken: FULL TOKEN: %s\n", token)
	fmt.Printf("DEBUG BROKER RefreshToken: Token length: %d\n", len(token))

	// Check for common issues with tokens
	if strings.Contains(token, " ") {
		fmt.Printf("DEBUG BROKER RefreshToken: Warning - token contains spaces\n")
	}
	if strings.Contains(token, "\n") || strings.Contains(token, "\r") {
		fmt.Printf("DEBUG BROKER RefreshToken: Warning - token contains newlines\n")
	}

	// Check token structure
	parts := strings.Split(token, ".")
	fmt.Printf("DEBUG BROKER RefreshToken: Token has %d parts\n", len(parts))
	if len(parts) == 3 {
		fmt.Printf("DEBUG BROKER RefreshToken: Header: %s\n", parts[0])
		fmt.Printf("DEBUG BROKER RefreshToken: Payload: %s\n", parts[1])
		fmt.Printf("DEBUG BROKER RefreshToken: Signature: %s\n", parts[2])
	} else {
		fmt.Printf("DEBUG BROKER RefreshToken: Invalid token format - should have 3 parts (header.payload.signature)\n")
	}

	// Call auth service to refresh token
	fmt.Printf("DEBUG BROKER RefreshToken: Calling auth service to refresh token\n")
	res, err := c.AuthClient.RefreshToken(r.Context(), &authpb.RefreshTokenRequest{
		RefreshToken: token,
	})

	if err != nil {
		fmt.Printf("DEBUG BROKER RefreshToken: Error from auth service: %v\n", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if !res.Success {
		fmt.Printf("DEBUG BROKER RefreshToken: Auth service returned unsuccessful response: %s\n", res.Message)
		util.ErrorJSON(w, errors.New(res.Message), http.StatusUnauthorized)
		return
	}

	fmt.Printf("DEBUG BROKER RefreshToken: Successfully refreshed token\n")
	responseData := util.JsonResponse{
		Error:   false,
		Message: "Token refreshed successfully",
		Data: map[string]interface{}{
			"access_token":  res.AccessToken,
			"refresh_token": res.RefreshToken,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// RequestPasswordReset handles password reset requests
func (c *Config) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call auth service via gRPC
	fmt.Printf("DEBUG BROKER RequestPasswordReset: Attempting password reset for email: %s\n", requestPayload.Email)
	res, err := c.AuthClient.RequestPasswordReset(r.Context(), &authpb.PasswordResetRequest{
		Email: requestPayload.Email,
	})

	if err != nil {
		fmt.Printf("DEBUG BROKER RequestPasswordReset: Error from auth service: %v\n", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Return response with additional data expected by frontend
	responseData := util.JsonResponse{
		Error:   !res.Success,
		Message: res.Message,
		Data: map[string]interface{}{
			"success":   res.Success,
			"message":   res.Message,
			"expiresIn": 15, // OTP expires in 15 minutes
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// VerifyPasswordReset handles OTP verification for password reset
func (c *Config) VerifyPasswordReset(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email   string `json:"email"`
		OTP     string `json:"otp"`
		Purpose string `json:"purpose"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call auth service via gRPC
	fmt.Printf("DEBUG BROKER VerifyPasswordReset: Verifying OTP for email: %s\n", requestPayload.Email)
	res, err := c.AuthClient.VerifyPasswordReset(r.Context(), &authpb.VerifyPasswordResetRequest{
		Email: requestPayload.Email,
		Otp:   requestPayload.OTP,
	})

	if err != nil {
		fmt.Printf("DEBUG BROKER VerifyPasswordReset: Error from auth service: %v\n", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if !res.Success {
		fmt.Printf("DEBUG BROKER VerifyPasswordReset: Auth service returned unsuccessful response: %s\n", res.Message)
		util.ErrorJSON(w, errors.New(res.Message), http.StatusBadRequest)
		return
	}

	// Return response with token field for compatibility with frontend
	responseData := util.JsonResponse{
		Error:   false,
		Message: res.Message,
		Data: map[string]interface{}{
			"success":     res.Success,
			"token":       res.ResetToken, // Use token field name for frontend compatibility
			"reset_token": res.ResetToken, // Keep original field for backward compatibility
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// UpdatePassword handles password update requests
func (c *Config) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		CurrentPassword string `json:"current_password,omitempty"`
		Email           string `json:"email,omitempty"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Log request details (but mask the password)
	var tokenPreview string
	if len(requestPayload.Token) > 10 {
		tokenPreview = requestPayload.Token[:10] + "..."
	} else {
		tokenPreview = requestPayload.Token + "..."
	}
	fmt.Printf("DEBUG BROKER UpdatePassword: Request payload - Token: %s Email: %s\n",
		tokenPreview,
		requestPayload.Email)

	// For password reset, email is required
	if requestPayload.CurrentPassword == "" && requestPayload.Email == "" {
		util.ErrorJSON(w, errors.New("email is required for password reset"), http.StatusBadRequest)
		return
	}

	// Implement a direct solution to ensure email is included in the request to auth service
	// Since we can't directly use the Email field due to proto regeneration issues,
	// we'll use a context-based solution to pass additional metadata
	ctx := r.Context()

	// Create a context with metadata containing the email
	md := metadata.New(map[string]string{
		"email": requestPayload.Email,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Print the metadata for debugging
	fmt.Printf("DEBUG BROKER UpdatePassword: Added email metadata: %s\n", requestPayload.Email)

	// Create basic request
	req := &authpb.UpdatePasswordRequest{
		Token:           requestPayload.Token,
		Password:        requestPayload.Password,
		CurrentPassword: requestPayload.CurrentPassword,
	}

	// Also try reflection as a backup
	reqVal := reflect.ValueOf(req).Elem()
	emailField := reqVal.FieldByName("Email")
	if emailField.IsValid() && emailField.CanSet() {
		emailField.SetString(requestPayload.Email)
		fmt.Printf("DEBUG BROKER UpdatePassword: Email field found in proto and set to: %s\n", requestPayload.Email)
	} else {
		fmt.Printf("DEBUG BROKER UpdatePassword: Email field not found in proto struct, using metadata instead\n")
	}

	// Call auth service with the context containing email metadata
	fmt.Printf("DEBUG BROKER UpdatePassword: Calling auth service with token and email metadata\n")
	res, err := c.AuthClient.UpdatePassword(ctx, req)

	if err != nil {
		fmt.Printf("DEBUG BROKER UpdatePassword: Error from auth service: %v\n", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if !res.Success {
		fmt.Printf("DEBUG BROKER UpdatePassword: Auth service returned unsuccessful response: %s\n", res.Message)
		util.ErrorJSON(w, errors.New(res.Message), http.StatusBadRequest)
		return
	}

	// Return response
	responseData := util.JsonResponse{
		Error:   false,
		Message: res.Message,
		Data: map[string]interface{}{
			"success": res.Success,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// RequestOTP handles OTP generation requests
func (c *Config) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email   string `json:"email"`
		Purpose string `json:"purpose"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate purpose - chỉ cho phép mục đích cụ thể
	validPurpose := false
	switch requestPayload.Purpose {
	case "registration", "password_reset":
		validPurpose = true
	}

	if !validPurpose {
		util.ErrorJSON(w, errors.New("invalid purpose, only 'registration' and 'password_reset' are allowed"), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if requestPayload.Email == "" {
		util.ErrorJSON(w, errors.New("email is required"), http.StatusBadRequest)
		return
	}

	// Call auth service via gRPC
	fmt.Printf("DEBUG BROKER RequestOTP: Requesting OTP for email: %s, purpose: %s\n",
		requestPayload.Email, requestPayload.Purpose)

	// Gửi yêu cầu đến auth service
	// Lưu ý: Phải tạo lại file auth_grpc.pb.go bằng protoc trước khi sử dụng
	res, err := c.AuthClient.RequestOTP(r.Context(), &authpb.OTPRequest{
		Email:   requestPayload.Email,
		Purpose: requestPayload.Purpose,
	})

	if err != nil {
		fmt.Printf("DEBUG BROKER RequestOTP: Error from auth service: %v\n", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Return response
	responseData := util.JsonResponse{
		Error:   !res.Success,
		Message: res.Message,
		Data: map[string]interface{}{
			"success":   res.Success,
			"message":   res.Message,
			"expiresIn": 15, // OTP expires in 15 minutes
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// VerifyOTP handles OTP verification
// func (c *Config) VerifyOTP(w http.ResponseWriter, r *http.Request) {
// 	var requestPayload struct {
// 		Email   string `json:"email"`
// 		OTP     string `json:"otp"`
// 		Purpose string `json:"purpose"`
// 	}

// 	err := util.ReadJSON(w, r, &requestPayload)
// 	if err != nil {
// 		util.ErrorJSON(w, err, http.StatusBadRequest)
// 		return
// 	}

// 	// Validate purpose - chỉ cho phép mục đích cụ thể
// 	validPurpose := false
// 	switch requestPayload.Purpose {
// 	case "registration", "password_reset":
// 		validPurpose = true
// 	}

// 	if !validPurpose {
// 		util.ErrorJSON(w, errors.New("invalid purpose, only 'registration' and 'password_reset' are allowed"), http.StatusBadRequest)
// 		return
// 	}

// 	// Validate required fields
// 	if requestPayload.Email == "" {
// 		util.ErrorJSON(w, errors.New("email is required"), http.StatusBadRequest)
// 		return
// 	}
// 	if requestPayload.OTP == "" {
// 		util.ErrorJSON(w, errors.New("OTP is required"), http.StatusBadRequest)
// 		return
// 	}

// 	// Use different handlers based on purpose
// 	if requestPayload.Purpose == "registration" {
// 		// For registration verification, use VerifyRegistration endpoint
// 		fmt.Printf("DEBUG BROKER VerifyOTP: Registration purpose detected, redirecting to VerifyRegistration for email: %s\n",
// 			requestPayload.Email)

// 		// Call auth service via gRPC for registration verification
// 		res, err := c.AuthClient.VerifyRegistration(r.Context(), &authpb.VerifyRegistrationRequest{
// 			Email: requestPayload.Email,
// 			Otp:   requestPayload.OTP,
// 		})

// 		if err != nil {
// 			fmt.Printf("DEBUG BROKER VerifyOTP->VerifyRegistration: Error from auth service: %v\n", err)
// 			util.ErrorJSON(w, err, http.StatusInternalServerError)
// 			return
// 		}

// 		if !res.Success {
// 			fmt.Printf("DEBUG BROKER VerifyOTP->VerifyRegistration: Auth service returned unsuccessful response: %s\n", res.Message)
// 			util.ErrorJSON(w, errors.New(res.Message), http.StatusBadRequest)
// 			return
// 		}

// 		// Return response with tokens for successful registration verification
// 		responseData := util.JsonResponse{
// 			Error:   false,
// 			Message: res.Message,
// 			Data: map[string]interface{}{
// 				"success":       res.Success,
// 				"user_id":       res.UserId,
// 				"access_token":  res.AccessToken,
// 				"refresh_token": res.RefreshToken,
// 			},
// 		}

// 		util.WriteJSON(w, http.StatusOK, responseData)
// 		return
// 	}

// 	// For other purposes (password_reset), use regular VerifyOTP
// 	fmt.Printf("DEBUG BROKER VerifyOTP: Verifying OTP for email: %s, purpose: %s\n",
// 		requestPayload.Email, requestPayload.Purpose)

// 	// Gửi yêu cầu đến auth service
// 	res, err := c.AuthClient.VerifyOTP(r.Context(), &authpb.VerifyOTPRequest{
// 		Email:   requestPayload.Email,
// 		Otp:     requestPayload.OTP,
// 		Purpose: requestPayload.Purpose,
// 	})

// 	if err != nil {
// 		fmt.Printf("DEBUG BROKER VerifyOTP: Error from auth service: %v\n", err)
// 		util.ErrorJSON(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	if !res.Success {
// 		fmt.Printf("DEBUG BROKER VerifyOTP: Auth service returned unsuccessful response: %s\n", res.Message)
// 		util.ErrorJSON(w, errors.New(res.Message), http.StatusBadRequest)
// 		return
// 	}

// 	// Return response
// 	responseData := util.JsonResponse{
// 		Error:   false,
// 		Message: res.Message,
// 		Data: map[string]interface{}{
// 			"success": res.Success,
// 			"token":   res.Token,
// 		},
// 	}

// 	util.WriteJSON(w, http.StatusOK, responseData)
// }

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CheckAccountExists checks if an account exists with the given email
func (c *Config) CheckAccountExists(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if requestPayload.Email == "" {
		util.ErrorJSON(w, errors.New("email is required"), http.StatusBadRequest)
		return
	}

	// Call auth service via gRPC
	fmt.Printf("DEBUG BROKER CheckAccountExists: Checking if account exists for email: %s\n", requestPayload.Email)

	// Gửi yêu cầu đến auth service
	res, err := c.AuthClient.CheckAccountExists(r.Context(), &authpb.CheckAccountExistsRequest{
		Email: requestPayload.Email,
	})

	if err != nil {
		fmt.Printf("DEBUG BROKER CheckAccountExists: Error from auth service: %v\n", err)
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Return response
	responseData := util.JsonResponse{
		Error:   false,
		Message: res.Message,
		Data: map[string]interface{}{
			"exists": res.Exists,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
