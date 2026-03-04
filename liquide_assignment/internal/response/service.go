package response

import (
	"encoding/json"
	"liquide_assignment/internal/config"
	"net/http"
)

func HandleError(w http.ResponseWriter, err config.Errors) {
	response := ErrorHttpHandling(err)
	w.WriteHeader(err.StatusCode)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ErrorHttpHandling(err config.Errors) APIResponse {
	return APIResponse{
		Message:    err.Message,
		Status:     string(err.Status),
		Code:       string(err.Code),
		StatusCode: err.StatusCode,
	}
}

func SendResponse(w http.ResponseWriter, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func NewAPIResponse(message string, status config.Status, statusCode int) APIResponse {
	return APIResponse{Message: message, Status: string(status), StatusCode: statusCode}
}

type APIResponse struct {
	Message    string      `json:"message,omitempty"`
	Status     string      `json:"status,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Code       string      `json:"code,omitempty"`
	StatusCode int         `json:"-"`
}
