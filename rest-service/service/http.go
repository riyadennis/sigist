package service

import (
	"encoding/json"
	"net/http"
)

// HTTPError encapsulates http error
type HTTPError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Error     string `json:"error"`
}

// Success encapsulates success response
type Success struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HTTPResponse returns a populated HTTP response object
func HTTPResponse(w http.ResponseWriter, err error, status int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	var response interface{}

	response = &Success{
		Code:    status,
		Message: message,
	}

	if err != nil {
		response = &HTTPError{
			ErrorCode: status,
			Message:   message,
			Error:     err.Error(),
		}
	}

	errData, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = w.Write(errData)
	if err != nil {
		return err
	}
	return nil
}
