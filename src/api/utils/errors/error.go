package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ApiError interface {
	GetStatus() int
	GetMessage() string
	GetErrors() string
}

type apiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Errors  string `json:"error,omitempty"`
}

func (e *apiError) GetStatus() int {
	return e.Status
}

func (e *apiError) GetMessage() string {
	return e.Message
}

func (e *apiError) GetErrors() string {
	return e.Errors
}

func NewApiError(status int, message string) ApiError {
	return &apiError{
		Status:  status,
		Message: message,
	}
}

func NewApiErrorFromBytes(body []byte) (ApiError, error) {
	var result apiError

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("invalid json for creating ApiError")
	}

	return &result, nil
}

func NewNotFoundApiError(message string) ApiError {
	return &apiError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func NewInternalServerError(message string) ApiError {
	return &apiError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

func NewBadRequestError(message string) ApiError {
	return &apiError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}
