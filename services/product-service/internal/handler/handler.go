package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"product-service/internal/domain"
)

// Config contains application dependencies and configuration
type Config struct {
	ProductService  domain.ProductService
	CategoryService domain.CategoryService
	BannerService   domain.BannerService
	AdsService      domain.AdsService
}

// Response is the standard response format
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

// ReadJSON reads JSON from request body into data
func (app *Config) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON sends a JSON response
func (app *Config) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ErrorResponse sends an error response
func (app *Config) ErrorResponse(w http.ResponseWriter, err error, status int) {
	payload := Response{
		Status:  status,
		Message: err.Error(),
	}

	app.WriteJSON(w, status, payload)
}

// SuccessResponse sends a success response
func (app *Config) SuccessResponse(w http.ResponseWriter, status int, payload Response) {
	app.WriteJSON(w, status, payload)
}

// ParseJSON parses JSON from request body
func (app *Config) ParseJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New("malformed JSON at position " + string(rune(syntaxError.Offset)))
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed JSON")
		case errors.As(err, &unmarshalTypeError):
			return errors.New("invalid value for field " + unmarshalTypeError.Field)
		case errors.As(err, &invalidUnmarshalError):
			return errors.New("invalid unmarshal error")
		case errors.Is(err, io.EOF):
			return errors.New("request body is empty")
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

// HealthCheck handler
func (app *Config) HealthCheck(w http.ResponseWriter, r *http.Request) {
	payload := Response{
		Status:  http.StatusOK,
		Message: "Product service is healthy",
	}

	app.WriteJSON(w, http.StatusOK, payload)
}
