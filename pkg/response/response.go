package response

import "net/http"

type ApiResponse struct {
	StatusCode int         `json:"status_code"`
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func SuccessResponse(message string, data interface{}) (int, ApiResponse) {
	return http.StatusOK, ApiResponse{
		StatusCode: http.StatusOK,
		Status:     true,
		Message:    message,
		Data:       data,
	}
}

func ErrorResponse(message string) (int, ApiResponse) {
	return http.StatusBadRequest, ApiResponse{
		StatusCode: http.StatusBadRequest,
		Status:     false,
		Message:    message,
		Data:       nil,
	}
}
