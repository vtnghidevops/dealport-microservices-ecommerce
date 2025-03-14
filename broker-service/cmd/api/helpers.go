package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Broker is a handler function that returns a JSON response
// with a message and data
type jsonResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := int64(1 << 20) // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes) // wrap the request body with a MaxBytesReader
	dec := json.NewDecoder(r.Body) 
	err := dec.Decode(data)

	// if there is an error decoding the request body, return an error
	if ( err != nil ) {
		return err
	}

	err = dec.Decode(&struct{}{}) // if the request body is bigger than 1MB, return an error
	if ( err != io.EOF ) {
		return errors.New("request body is too large")
	}
 
	return nil
}

func (app *Config) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t") // Marshal data into JSON
	if ( err != nil ) {
		return err
	}

	// add headers to response 
	if ( len(headers) > 0 ) {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if (err != nil) {
	 return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if ( len(status) > 0 ) {
		statusCode = status[0]
	}
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()
	app.writeJson(w, statusCode, payload)
}

