package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// CreateTestRequest tạo một HTTP request giả lập cho việc kiểm thử
func CreateTestRequest(method, url string, requestBody interface{}) (*http.Request, error) {
	var reqBody []byte
	var err error

	if requestBody != nil {
		reqBody, err = json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// ParseTestResponse đọc response body và chuyển thành cấu trúc đã định nghĩa
func ParseTestResponse(rec *httptest.ResponseRecorder, response interface{}) error {
	return json.Unmarshal(rec.Body.Bytes(), response)
}

// SetAuthHeader thiết lập Authorization header với token Bearer
func SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

// CreateMockResponseRecorder tạo một HTTP response recorder giả lập
func CreateMockResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// GetStatusCodeAndBody trả về status code và body dạng map từ response
func GetStatusCodeAndBody(rec *httptest.ResponseRecorder) (int, map[string]interface{}, error) {
	var responseBody map[string]interface{}

	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	if err != nil {
		return rec.Code, nil, err
	}

	return rec.Code, responseBody, nil
}

// CreateTestRequestWithJSON tạo một HTTP request với body JSON đã được encode
func CreateTestRequestWithJSON(method, url string, jsonBody string) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBufferString(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
} 