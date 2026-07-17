package httpapi

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(response http.ResponseWriter, status int, code, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-store")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(errorResponse{Code: code, Message: message})
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
