package util

import (
	"broker-service/internal/util"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Cấu trúc data test cho readJSON
type testData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// TestReadJSON kiểm tra chức năng đọc JSON từ request
func TestReadJSON(t *testing.T) {
	// Test case 1: JSON hợp lệ
	t.Run("valid json", func(t *testing.T) {
		// Chuẩn bị request với JSON hợp lệ
		jsonData := `{"name":"Test User","email":"test@example.com","age":30}`
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(jsonData))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		// Chuẩn bị response recorder
		rr := httptest.NewRecorder()

		// Đối tượng cần test
		var data testData

		// Thực hiện đọc JSON
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := util.ReadJSON(w, r, &data)
			assert.NoError(t, err)

			// Kiểm tra dữ liệu đã parse
			assert.Equal(t, "Test User", data.Name)
			assert.Equal(t, "test@example.com", data.Email)
			assert.Equal(t, 30, data.Age)
		})

		handler.ServeHTTP(rr, req)

		// Kiểm tra status code
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test case 2: JSON không hợp lệ
	t.Run("invalid json", func(t *testing.T) {
		// Chuẩn bị request với JSON không hợp lệ
		invalidJSON := `{"name":"Test User","email":"test@example.com",age:30}`
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(invalidJSON))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		// Chuẩn bị response recorder
		rr := httptest.NewRecorder()

		// Đối tượng cần test
		var data testData

		// Thực hiện đọc JSON
		handlerWithError := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := util.ReadJSON(w, r, &data)
			assert.Error(t, err)
		})

		handlerWithError.ServeHTTP(rr, req)

		// Kiểm tra status code - nên là 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// Test case 3: Empty body
	t.Run("empty body", func(t *testing.T) {
		// Chuẩn bị request với body trống
		req, err := http.NewRequest("POST", "/test", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		// Chuẩn bị response recorder
		rr := httptest.NewRecorder()

		// Đối tượng cần test
		var data testData

		// Thực hiện đọc JSON
		handlerWithError := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := util.ReadJSON(w, r, &data)
			assert.Error(t, err)
		})

		handlerWithError.ServeHTTP(rr, req)

		// Kiểm tra status code - nên là 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

// TestWriteJSON kiểm tra chức năng ghi JSON vào response
func TestWriteJSON(t *testing.T) {
	// Test case 1: Ghi JSON thành công
	t.Run("successful write", func(t *testing.T) {
		// Chuẩn bị response recorder
		rr := httptest.NewRecorder()

		// Dữ liệu cần ghi
		data := testData{
			Name:  "Test User",
			Email: "test@example.com",
			Age:   30,
		}

		// Thực hiện ghi JSON
		err := util.WriteJSON(rr, http.StatusOK, data)
		assert.NoError(t, err)

		// Kiểm tra status code
		assert.Equal(t, http.StatusOK, rr.Code)

		// Kiểm tra Content-Type
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		// Kiểm tra response body
		var responseData testData
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		assert.NoError(t, err)
		assert.Equal(t, data.Name, responseData.Name)
		assert.Equal(t, data.Email, responseData.Email)
		assert.Equal(t, data.Age, responseData.Age)
	})

	// Test case 2: Ghi JSON với status code khác
	t.Run("write with different status code", func(t *testing.T) {
		// Chuẩn bị response recorder
		rr := httptest.NewRecorder()

		// Dữ liệu cần ghi
		data := testData{
			Name:  "Test User",
			Email: "test@example.com",
			Age:   30,
		}

		// Thực hiện ghi JSON với status code Created
		err := util.WriteJSON(rr, http.StatusCreated, data)
		assert.NoError(t, err)

		// Kiểm tra status code
		assert.Equal(t, http.StatusCreated, rr.Code)
	})
}

// TestErrorJSON kiểm tra chức năng ghi lỗi JSON
func TestErrorJSON(t *testing.T) {
	// Test case 1: Lỗi với status code mặc định
	t.Run("error with default status code", func(t *testing.T) {
		// Chuẩn bị response recorder
		rr := httptest.NewRecorder()

		// Thực hiện ghi lỗi JSON
		util.ErrorJSON(rr, errors.New("test error"))

		// Kiểm tra status code mặc định - nên là 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Kiểm tra Content-Type
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		// Kiểm tra response body
		var response struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Error)
		assert.Equal(t, "test error", response.Message)
	})

	// Test case 2: Lỗi với status code tùy chỉnh
	t.Run("error with custom status code", func(t *testing.T) {
		// Chuẩn bị response recorder
		rr := httptest.NewRecorder()

		// Thực hiện ghi lỗi JSON với status code Unauthorized
		util.ErrorJSON(rr, errors.New("unauthorized access"), http.StatusUnauthorized)

		// Kiểm tra status code
		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		// Kiểm tra response body
		var response struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Error)
		assert.Equal(t, "unauthorized access", response.Message)
	})
}
