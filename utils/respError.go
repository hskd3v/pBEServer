package utils

import (
	"net/http"

	"github.com/harriklein/pBE/pBEServer/log"
)

// TResponseError is the structure with the Error
type TResponseError struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode,omitempty"`
	StatusText string `json:"statusText,omitempty"`
	Message    string `json:"message,omitempty"`
}

// NewResponseError creates a bad request error message
func NewResponseError(aCode int, aMessage string) *TResponseError {
	_status := "error"
	if aCode < http.StatusInternalServerError {
		_status = "fail"
	}

	log.Log.Errorf("%s (%d-%s): %s\n", _status, aCode, http.StatusText(aCode), aMessage)

	return &TResponseError{
		Status:     _status,
		StatusCode: aCode,
		StatusText: http.StatusText(aCode),
		Message:    aMessage,
	}
}

// ToJSON d
func (oResponseError TResponseError) ToJSON(aResponse http.ResponseWriter) error {

	aResponse.Header().Add("Content-Type", "application/json")
	aResponse.WriteHeader(oResponseError.StatusCode)

	_error := ToJSON(oResponseError, aResponse)
	return _error

}
