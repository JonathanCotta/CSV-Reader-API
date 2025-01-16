package utils

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, m string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	resp := Message{
		Error:   false,
		Message: m,
	}

	json.NewEncoder(w).Encode(resp)
}

func DataResponse(w http.ResponseWriter, m string, d interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	resp := Message{
		Error:   false,
		Message: m,
		Data:    d,
	}

	json.NewEncoder(w).Encode(resp)
}

func ErrorResponse(w http.ResponseWriter, m string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	resp := Message{
		Error:   true,
		Message: m,
	}

	json.NewEncoder(w).Encode(resp)
}
