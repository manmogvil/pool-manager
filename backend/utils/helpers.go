package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}

func RespondValidationError(w http.ResponseWriter, err *ValidationError) {
	RespondJSON(w, http.StatusBadRequest, map[string]string{
		"error":   "validation failed",
		"field":   err.Field,
		"message": err.Message,
	})
}

func ParseID(raw string) (int, error) {
	return strconv.Atoi(raw)
}

func ParseDate(raw string) (time.Time, error) {
	return time.Parse("2006-01-02", raw)
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
