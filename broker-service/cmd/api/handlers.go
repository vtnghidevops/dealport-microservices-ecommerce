package main

import (
	"encoding/json"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error: false,
		Message: "Broker service is up and running",
		Data: nil,
	}
	_ = app.writeJson(w, http.StatusOK, payload) // Write response to client

	// out, _ := json.MarshalIndent(payload, "", "\t") 
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out) // Write response to client
}