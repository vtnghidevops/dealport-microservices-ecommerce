package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"broker-service/internal/util"
	authpb "broker-service/proto/auth"
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

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
