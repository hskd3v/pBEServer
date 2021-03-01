package utils

import (
	"net/http"
)

// TResponse is the standard structure
type TResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"statusCode,omitempty"`
	StatusText string      `json:"statusText,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// NewResponse creates a bad request Success message
func NewResponse(aCode int, aMessage string, aData interface{}) *TResponse {
	_status := "Unknown"
	switch aCode / 100 {
	case 0, 1, 2, 3:
		_status = "success"
	case 4:
		_status = "error"
	case 5:
		_status = "fail"
	}

	return &TResponse{
		Status:     _status,
		StatusCode: aCode,
		StatusText: http.StatusText(aCode),
		Message:    aMessage,
		Data:       aData,
	}
}

// ToJSON d
func (oResponse TResponse) ToJSON(aResponse http.ResponseWriter) error {

	aResponse.Header().Add("Content-Type", "application/json")
	aResponse.WriteHeader(oResponse.StatusCode)

	_error := ToJSON(oResponse, aResponse)
	return _error
}
