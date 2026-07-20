package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	type jsonResponse struct {
		Body string `json:"body"`
	}
	response := jsonResponse{
		Body: "OK",
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		writeJsonError(w, err.Error(), 500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(bytes)
}

func writeJsonError(w http.ResponseWriter, message string, statusCode int) {
	type jsonError struct {
		Body string `json:"body"`
		Code int    `json:"code"`
	}
	errBody := jsonError{
		Body: message,
		Code: statusCode,
	}
	bytes, _ := json.Marshal(errBody)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bytes)
}
