package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"valid number", "123", 123, false},
		{"zero", "0", 0, false},
		{"negative", "-5", -5, false},
		{"not a number", "abc", 0, true},
		{"empty string", "", 0, true},
		{"float", "12.5", 0, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ParseID(testCase.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("ParseID() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if got != testCase.want {
				t.Errorf("ParseID() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid date", "2024-01-15", false},
		{"invalid format", "15/01/2024", true},
		{"invalid date", "2024-13-01", true},
		{"empty string", "", true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := ParseDate(testCase.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{Field: "name", Message: "is required"}
	got := err.Error()
	want := "name: is required"
	if got != want {
		t.Errorf("ValidationError.Error() = %v, want %v", got, want)
	}
}

func TestRespondJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"status": "ok"}

	RespondJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %v, want application/json", contentType)
	}

	var result map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("response status = %v, want ok", result["status"])
	}
}

func TestRespondError(t *testing.T) {
	rr := httptest.NewRecorder()

	RespondError(rr, http.StatusNotFound, "not found")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusNotFound)
	}

	var result map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["error"] != "not found" {
		t.Errorf("error = %v, want 'not found'", result["error"])
	}
}

func TestRespondValidationError(t *testing.T) {
	rr := httptest.NewRecorder()
	validationErr := &ValidationError{Field: "email", Message: "is required"}

	RespondValidationError(rr, validationErr)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}

	var result map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["field"] != "email" {
		t.Errorf("field = %v, want 'email'", result["field"])
	}
	if result["message"] != "is required" {
		t.Errorf("message = %v, want 'is required'", result["message"])
	}
}
