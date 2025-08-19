package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse sends a JSON response with the given status code and data
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// ErrorResponse sends a JSON error response
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

// SuccessResponse sends a JSON success response
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// CreatedResponse sends a JSON created response
func CreatedResponse(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// ValidationErrorResponse sends a JSON validation error response
func ValidationErrorResponse(w http.ResponseWriter, errors map[string]string) {
	JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
		"error":   true,
		"message": "Validation failed",
		"errors":  errors,
	})
}
