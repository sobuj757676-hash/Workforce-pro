package httpapi

import (
	"encoding/json"
	"net/http"
)

// APIError is the stable error response envelope returned by all API endpoints.
// It includes a correlation identifier for tracing and optional conflict details.
// No cross-Tenant detail leakage is permitted in any field.
// Requirements: 7.5–7.10
type APIError struct {
	Code          string          `json:"code"`
	Message       string          `json:"message"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	Conflict      *ConflictDetail `json:"conflict,omitempty"`
}

// ConflictDetail describes an optimistic concurrency conflict without
// leaking protected cross-Tenant information.
// Requirements: 7.5–7.7
type ConflictDetail struct {
	CurrentVersion int64    `json:"current_version"`
	ChangedFields  []string `json:"changed_fields,omitempty"`
	RequiresReview bool     `json:"requires_review"`
}

func writeError(response http.ResponseWriter, status int, code, message string) {
	writeErrorWithCorrelation(response, status, code, message, "")
}

func writeErrorWithCorrelation(response http.ResponseWriter, status int, code, message, correlationID string) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-store")
	response.WriteHeader(status)
	body := APIError{
		Code:          code,
		Message:       message,
		CorrelationID: correlationID,
	}
	_ = json.NewEncoder(response).Encode(body)
}

func writeConflictError(response http.ResponseWriter, correlationID string, detail ConflictDetail) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-store")
	response.WriteHeader(http.StatusConflict)
	body := APIError{
		Code:          "VERSION_CONFLICT",
		Message:       "The resource was modified. Review the current version before resubmission.",
		CorrelationID: correlationID,
		Conflict:      &detail,
	}
	_ = json.NewEncoder(response).Encode(body)
}

func writeUnauthenticated(response http.ResponseWriter) {
	response.Header().Set("WWW-Authenticate", "Bearer")
	writeError(response, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication is required.")
}

func writeAuthenticationUnavailable(response http.ResponseWriter) {
	writeError(response, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "The service is temporarily unavailable.")
}

func writeTenantScopedNotFound(response http.ResponseWriter) {
	writeError(response, http.StatusNotFound, "RESOURCE_NOT_FOUND", "The requested resource was not found.")
}

func writeForbidden(response http.ResponseWriter, correlationID string) {
	writeErrorWithCorrelation(response, http.StatusForbidden, "FORBIDDEN", "You do not have permission to perform this action.", correlationID)
}

func writeInvalidInput(response http.ResponseWriter, correlationID, detail string) {
	writeErrorWithCorrelation(response, http.StatusBadRequest, "INVALID_INPUT", detail, correlationID)
}

func writeIdempotencyConflict(response http.ResponseWriter, correlationID string) {
	writeErrorWithCorrelation(response, http.StatusConflict, "IDEMPOTENCY_KEY_REUSED", "The idempotency key was reused with a different request.", correlationID)
}
