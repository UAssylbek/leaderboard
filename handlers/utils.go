package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse — структура для ошибок в JSON
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// SendError — утилита для отправки JSON-ошибки
func SendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  statusCode,
	})
}