package integration

import (
	"broker-service/internal/handlers/http/auth"
	authpb "broker-service/proto/auth"
	"broker-service/tests/helpers"
	"broker-service/tests/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestAuthenticationFlow kiểm tra toàn bộ luồng xác thực
func TestAuthenticationFlow(t *testing.T) {
	// Tạo mock cho AuthServiceClient bằng adapter
	mockAuthClient := mocks.SetupMockAuthClient()

	// Tạo config với mock client
	handler := &auth.Config{
		AuthClient: mocks.NewAuthClientAdapter(mockAuthClient),
	}

	// Dữ liệu cho luồng test
	testUser := struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Username  string
		UserID    string
		OTP       string
	}{
		Email:     "test@example.com",
		Password:  "Secure123!",
		FirstName: "Test",
		LastName:  "User",
		Username:  "testuser",
		UserID:    "user123",
		OTP:       "123456",
	}

	// Mock tokens
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMTIzIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMn0.fake_signature"
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMTIzIiwicmVmcmVzaCI6dHJ1ZSwiaWF0IjoxNTE2MjM5MDIyfQ.fake_signature"

	// Bước 1: Đăng ký tài khoản mới
	t.Run("1. Register new account", func(t *testing.T) {
		// Chuẩn bị request body
		reqBody := map[string]interface{}{
			"email":      testUser.Email,
			"password":   testUser.Password,
			"first_name": testUser.FirstName,
			"last_name":  testUser.LastName,
			"username":   testUser.Username,
		}

		// Chuẩn bị mock response
		mockRegisterResponse := &authpb.RegisterResponse{
			Success: true,
			Message: "Registration successful. Please check your email for verification code.",
			UserId:  testUser.UserID,
		}

		// Cấu hình mock
		mockAuthClient.On("Register", mock.Anything, mock.MatchedBy(func(req *authpb.RegisterRequest) bool {
			return req.Email == testUser.Email &&
				req.Password == testUser.Password &&
				req.FirstName == testUser.FirstName &&
				req.LastName == testUser.LastName &&
				req.Username == testUser.Username
		}), mock.Anything).Return(mockRegisterResponse, nil).Once()

		// Tạo HTTP request và response recorder
		req, err := helpers.CreateTestRequest("POST", "/auth/register", reqBody)
		assert.NoError(t, err)

		rr := helpers.CreateMockResponseRecorder()

		// Thực hiện request
		handler.Register(rr, req)

		// Kiểm tra response
		statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, statusCode)

		// Kiểm tra response data
		assert.Equal(t, false, respBody["error"])
		assert.Contains(t, respBody["message"], "Registration successful")

		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, testUser.UserID, data["user_id"])
		assert.Equal(t, true, data["success"])

		// Bước 2: Xác thực OTP
		t.Run("2. Verify OTP", func(t *testing.T) {
			// Chuẩn bị request body
			reqBody := map[string]interface{}{
				"email": testUser.Email,
				"otp":   testUser.OTP,
			}

			// Chuẩn bị mock response
			mockVerifyResponse := &authpb.VerifyRegistrationResponse{
				Success:      true,
				Message:      "Account verified successfully",
				UserId:       testUser.UserID,
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}

			// Cấu hình mock
			mockAuthClient.On("VerifyRegistration", mock.Anything, mock.MatchedBy(func(req *authpb.VerifyRegistrationRequest) bool {
				return req.Email == testUser.Email && req.Otp == testUser.OTP
			}), mock.Anything).Return(mockVerifyResponse, nil).Once()

			// Tạo HTTP request và response recorder
			req, err := helpers.CreateTestRequest("POST", "/auth/verify", reqBody)
			assert.NoError(t, err)

			rr := helpers.CreateMockResponseRecorder()

			// Thực hiện request
			handler.VerifyRegistration(rr, req)

			// Kiểm tra response
			statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, statusCode)

			// Kiểm tra response data
			assert.Equal(t, false, respBody["error"])
			assert.Contains(t, respBody["message"], "Account verified successfully")

			data, ok := respBody["data"].(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, testUser.UserID, data["user_id"])
			assert.Equal(t, accessToken, data["access_token"])
			assert.Equal(t, refreshToken, data["refresh_token"])

			// Bước 3: Đăng nhập
			t.Run("3. Login", func(t *testing.T) {
				// Chuẩn bị request body
				reqBody := map[string]interface{}{
					"email":    testUser.Email,
					"password": testUser.Password,
				}

				// Chuẩn bị mock response
				mockLoginResponse := &authpb.LoginResponse{
					Success:      true,
					Message:      "Logged in successfully",
					UserId:       testUser.UserID,
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
					UserInfo: &authpb.UserInfo{
						FirstName: testUser.FirstName,
						LastName:  testUser.LastName,
						Email:     testUser.Email,
					},
				}

				// Cấu hình mock
				mockAuthClient.On("Login", mock.Anything, mock.MatchedBy(func(req *authpb.LoginRequest) bool {
					return req.Email == testUser.Email && req.Password == testUser.Password
				}), mock.Anything).Return(mockLoginResponse, nil).Once()

				// Tạo HTTP request và response recorder
				req, err := helpers.CreateTestRequest("POST", "/auth/login", reqBody)
				assert.NoError(t, err)

				rr := helpers.CreateMockResponseRecorder()

				// Thực hiện request
				handler.Login(rr, req)

				// Kiểm tra response
				statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, statusCode)

				// Kiểm tra response data
				assert.Equal(t, false, respBody["error"])
				assert.Contains(t, respBody["message"], "Logged in successfully")

				data, ok := respBody["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, testUser.UserID, data["user_id"])
				assert.Equal(t, accessToken, data["access_token"])
				assert.Equal(t, refreshToken, data["refresh_token"])

				// Bước 4: Validate token
				t.Run("4. Validate token", func(t *testing.T) {
					// Chuẩn bị mock response
					mockValidateResponse := &authpb.ValidateResponse{
						Valid:  true,
						UserId: testUser.UserID,
						Claims: map[string]string{
							"sub":  testUser.UserID,
							"name": testUser.FirstName + " " + testUser.LastName,
						},
					}

					// Cấu hình mock
					mockAuthClient.On("Validate", mock.Anything, mock.MatchedBy(func(req *authpb.ValidateRequest) bool {
						return req.Token == accessToken
					}), mock.Anything).Return(mockValidateResponse, nil).Once()

					// Tạo HTTP request và response recorder
					req, err := http.NewRequest("POST", "/auth/validate", nil)
					assert.NoError(t, err)

					// Thêm Authorization header
					helpers.SetAuthHeader(req, accessToken)

					rr := helpers.CreateMockResponseRecorder()

					// Thực hiện request
					handler.ValidateToken(rr, req)

					// Kiểm tra response
					statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusOK, statusCode)

					// Kiểm tra response data
					assert.Equal(t, false, respBody["error"])
					assert.Contains(t, respBody["message"], "Token is valid")

					data, ok := respBody["data"].(map[string]interface{})
					assert.True(t, ok)
					assert.Equal(t, testUser.UserID, data["user_id"])

					// Bước 5: Refresh token
					t.Run("5. Refresh token", func(t *testing.T) {
						// Chuẩn bị request body
						reqBody := map[string]interface{}{
							"refresh_token": refreshToken,
						}

						// Chuẩn bị mock response
						newAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMTIzIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyM30.new_fake_signature"

						mockRefreshResponse := &authpb.RefreshTokenResponse{
							Success:      true,
							AccessToken:  newAccessToken,
							RefreshToken: refreshToken,
							Message:      "Token refreshed successfully",
						}

						// Cấu hình mock
						mockAuthClient.On("RefreshToken", mock.Anything, mock.MatchedBy(func(req *authpb.RefreshTokenRequest) bool {
							return req.RefreshToken == refreshToken
						}), mock.Anything).Return(mockRefreshResponse, nil).Once()

						// Tạo HTTP request và response recorder
						req, err := helpers.CreateTestRequest("POST", "/auth/refresh", reqBody)
						assert.NoError(t, err)

						rr := helpers.CreateMockResponseRecorder()

						// Thực hiện request
						handler.RefreshToken(rr, req)

						// Kiểm tra response
						statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
						assert.NoError(t, err)
						assert.Equal(t, http.StatusOK, statusCode)

						// Kiểm tra response data
						assert.Equal(t, false, respBody["error"])

						data, ok := respBody["data"].(map[string]interface{})
						assert.True(t, ok)
						assert.Equal(t, newAccessToken, data["access_token"])
					})
				})
			})
		})
	})

	// Verify tất cả mock calls
	mockAuthClient.AssertExpectations(t)
}
