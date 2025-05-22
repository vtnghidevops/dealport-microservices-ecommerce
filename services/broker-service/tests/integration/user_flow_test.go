package integration

import (
	"broker-service/internal/handlers/http/user"
	userpb "broker-service/proto/user"
	"broker-service/tests/helpers"
	"broker-service/tests/mocks"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUserFlow kiểm tra luồng quản lý thông tin người dùng
func TestUserFlow(t *testing.T) {
	// Tạo mock cho UserServiceClient và CheckoutServiceClient
	mockUserClient := mocks.SetupMockUserClient()
	mockCheckoutClient := new(mocks.MockCheckoutServiceClient)

	// Tạo config với mock client
	handler := &user.Config{
		UserClient:     mocks.NewUserClientAdapter(mockUserClient),
		CheckoutClient: mockCheckoutClient,
	}

	// Dữ liệu kiểm thử
	testUser := struct {
		ID        string
		Email     string
		FirstName string
		LastName  string
		Phone     string
		Avatar    string
		Address   string
		City      string
		Country   string
		ProductID string
	}{
		ID:        "user123",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "1234567890",
		Avatar:    "https://example.com/avatar.jpg",
		Address:   "123 Test St",
		City:      "Test City",
		Country:   "Test Country",
		ProductID: "123", // Sử dụng ID dạng số để dễ chuyển đổi sang int32
	}

	// Access token cho request
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMTIzIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMn0.fake_signature"

	// Bước 1: Lấy thông tin profile người dùng
	t.Run("1. Get user profile", func(t *testing.T) {
		// Chuẩn bị mock response
		mockUserResponse := &userpb.UserResponse{
			User: &userpb.User{
				Id:           testUser.ID,
				Email:        testUser.Email,
				FirstName:    testUser.FirstName,
				LastName:     testUser.LastName,
				Phone:        testUser.Phone,
				ProfileImage: testUser.Avatar,
				Addresses: []*userpb.Address{
					{
						Id:          "addr1",
						Line1:       testUser.Address,
						City:        testUser.City,
						Country:     testUser.Country,
						IsDefault:   true,
						AddressType: "shipping",
					},
				},
			},
		}

		// Cấu hình mock
		mockUserClient.On("GetUser", mock.Anything, &userpb.GetUserRequest{
			Id: testUser.ID,
		}, mock.Anything).Return(mockUserResponse, nil).Once()

		// Tạo HTTP request
		req, err := http.NewRequest("GET", "/api/v1/users/me", nil)
		assert.NoError(t, err)

		// Thêm Authorization header và user_id trong context
		helpers.SetAuthHeader(req, accessToken)
		// Mock context - trong test thực tế, user_id sẽ được middleware xác thực thêm vào
		req = req.WithContext(context.WithValue(req.Context(), "user_id", testUser.ID))

		rr := helpers.CreateMockResponseRecorder()

		// Thực hiện request
		handler.GetUserProfile(rr, req)

		// Kiểm tra response
		statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)

		// Kiểm tra dữ liệu người dùng
		assert.Equal(t, false, respBody["error"])
		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, testUser.ID, data["id"])
		assert.Equal(t, testUser.Email, data["email"])
		assert.Equal(t, testUser.FirstName, data["first_name"])
		assert.Equal(t, testUser.LastName, data["last_name"])

		// Bước 2: Cập nhật thông tin người dùng
		t.Run("2. Update user profile", func(t *testing.T) {
			// Chuẩn bị request body
			reqBody := map[string]interface{}{
				"firstName": "Updated",
				"lastName":  "User",
				"phone":     "9876543210",
				"avatar":    testUser.Avatar,
				"addresses": []map[string]interface{}{
					{
						"id":        "addr1",
						"street":    "456 New St",
						"city":      "New City",
						"state":     "New State",
						"zipCode":   "12345",
						"country":   "New Country",
						"isDefault": true,
						"type":      "shipping",
					},
				},
			}

			// Chuẩn bị mock response
			mockUpdateResponse := &userpb.UserResponse{
				User: &userpb.User{
					Id:           testUser.ID,
					Email:        testUser.Email,
					FirstName:    "Updated",
					LastName:     "User",
					Phone:        "9876543210",
					ProfileImage: testUser.Avatar,
					Addresses: []*userpb.Address{
						{
							Id:          "addr1",
							Line1:       "456 New St",
							City:        "New City",
							State:       "New State",
							PostalCode:  "12345",
							Country:     "New Country",
							IsDefault:   true,
							AddressType: "shipping",
						},
					},
				},
			}

			// Cấu hình mock - Sử dụng mock.Anything cho updateReq vì khó match chính xác
			mockUserClient.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(mockUpdateResponse, nil).Once()

			// Tạo HTTP request
			req, err := helpers.CreateTestRequest("PUT", "/api/v1/users/me", reqBody)
			assert.NoError(t, err)

			// Thêm Authorization header và user_id trong context
			helpers.SetAuthHeader(req, accessToken)
			// Mock context
			req = req.WithContext(context.WithValue(req.Context(), "user_id", testUser.ID))

			rr := helpers.CreateMockResponseRecorder()

			// Thực hiện request
			handler.UpdateUserProfile(rr, req)

			// Kiểm tra response
			statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, statusCode)

			// Kiểm tra response data
			assert.Equal(t, false, respBody["error"])
			assert.Contains(t, respBody["message"], "User profile updated successfully")

			// Bước 3: Thêm sản phẩm vào wishlist
			t.Run("3. Add to wishlist", func(t *testing.T) {
				// Chuyển đổi ID sản phẩm từ string sang int32
				productID, _ := strconv.Atoi(testUser.ProductID)

				// Chuẩn bị request body
				reqBody := map[string]interface{}{
					"product_id": productID,
				}

				// Chuẩn bị mock response
				mockWishlistResponse := &userpb.WishlistResponse{
					Success: true,
					Message: "Product added to wishlist",
				}

				// Cấu hình mock
				mockUserClient.On("AddToWishlist", mock.Anything, &userpb.AddToWishlistRequest{
					UserId:    testUser.ID,
					ProductId: int32(productID),
				}, mock.Anything).Return(mockWishlistResponse, nil).Once()

				// Tạo HTTP request
				req, err := helpers.CreateTestRequest("POST", "/api/v1/users/me/wishlist", reqBody)
				assert.NoError(t, err)

				// Thêm Authorization header và user_id trong context
				helpers.SetAuthHeader(req, accessToken)
				// Mock context
				req = req.WithContext(context.WithValue(req.Context(), "user_id", testUser.ID))

				rr := helpers.CreateMockResponseRecorder()

				// Thực hiện request
				handler.AddToWishlist(rr, req)

				// Kiểm tra response
				statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, statusCode)

				// Kiểm tra response data
				assert.Equal(t, false, respBody["error"])
				assert.Contains(t, respBody["message"], "Product added to wishlist")

				// Bước 4: Xem wishlist
				t.Run("4. Get wishlist", func(t *testing.T) {
					// Chuẩn bị mock response
					mockGetWishlistResponse := &userpb.GetWishlistResponse{
						Items: []*userpb.WishlistItem{
							{
								Id:        "wishlist_1",
								UserId:    testUser.ID,
								ProductId: int32(productID),
								AddedAt:   "2023-06-01T10:00:00Z",
							},
						},
						Count:   1,
						Message: "Wishlist retrieved successfully",
					}

					// Cấu hình mock
					mockUserClient.On("GetWishlist", mock.Anything, &userpb.GetWishlistRequest{
						UserId: testUser.ID,
					}, mock.Anything).Return(mockGetWishlistResponse, nil).Once()

					// Tạo HTTP request
					req, err := http.NewRequest("GET", "/api/v1/users/me/wishlist", nil)
					assert.NoError(t, err)

					// Thêm Authorization header và user_id trong context
					helpers.SetAuthHeader(req, accessToken)
					// Mock context
					req = req.WithContext(context.WithValue(req.Context(), "user_id", testUser.ID))

					rr := helpers.CreateMockResponseRecorder()

					// Thực hiện request
					handler.GetWishlist(rr, req)

					// Kiểm tra response
					statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusOK, statusCode)

					// Kiểm tra response data
					assert.Equal(t, false, respBody["error"])

					// Bước 5: Xóa sản phẩm khỏi wishlist
					t.Run("5. Remove from wishlist", func(t *testing.T) {
						// Sử dụng productID đã chuyển đổi trước đó
						productID, _ := strconv.Atoi(testUser.ProductID)

						// Chuẩn bị mock response
						mockRemoveResponse := &userpb.WishlistResponse{
							Success: true,
							Message: "Product removed from wishlist",
						}

						// Cấu hình mock
						mockUserClient.On("RemoveFromWishlist", mock.Anything, &userpb.RemoveFromWishlistRequest{
							UserId:    testUser.ID,
							ProductId: int32(productID),
						}, mock.Anything).Return(mockRemoveResponse, nil).Once()

						// Tạo HTTP request - sử dụng URL parameter thay vì truyền vào body
						req, err := http.NewRequest("DELETE", "/api/v1/users/me/wishlist/"+fmt.Sprint(productID), nil)
						assert.NoError(t, err)

						// Thêm Authorization header và user_id trong context
						helpers.SetAuthHeader(req, accessToken)
						// Mock context
						req = req.WithContext(context.WithValue(req.Context(), "user_id", testUser.ID))

						rr := helpers.CreateMockResponseRecorder()

						// Thực hiện request
						handler.RemoveFromWishlist(rr, req)

						// Kiểm tra response
						statusCode, respBody, err := helpers.GetStatusCodeAndBody(rr)
						assert.NoError(t, err)
						assert.Equal(t, http.StatusOK, statusCode)

						// Kiểm tra response data
						assert.Equal(t, false, respBody["error"])
						assert.Contains(t, respBody["message"], "Product removed from wishlist")
					})
				})
			})
		})
	})

	// Verify tất cả mock calls
	mockUserClient.AssertExpectations(t)
	mockCheckoutClient.AssertExpectations(t)
}
