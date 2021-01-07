package httplib

import (
	"fmt"
	"time"
)

type ApiError struct {
	HttpStatus int         `json:"http_status, omitempty"`
	ErrorCode  int         `json:"code"`
	Message    interface{} `json:"message"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func (e ApiError) String() string {
	return fmt.Sprintf("%s", e.Message)
}

func NewApiError(code int, msg string) *ApiError {
	return &ApiError{
		ErrorCode: code,
		Message:   msg,
	}
}

type ResultAdditionalFields map[string]interface{}

type CreatedObjectId struct {
	Id interface{} `json:"id, omitempty"`
}

type ListResult struct {
	Results interface{} `json:"results, omitempty"`
	Count   int32       `json:"count"`
	HasPrev bool        `json:"has_prev"`
	HasNext bool        `json:"has_next"`
}

type SuccessfulRequestResult struct {
	Success   bool        `json:"success"`
	Timestamp interface{} `json:"timestamp, omitempty"`
	Message   *string     `json:"message, omitempty"`
	CreatedObjectId
}

func SuccessfulResult() SuccessfulRequestResult {
	msg := "Well done, mate!"
	return SuccessfulRequestResult{
		Success:   true,
		Timestamp: time.Now().Unix(),
		Message:   &msg,
	}
}

func SuccessfulResultMap() map[string]interface{} {
	return map[string]interface{}{
		"success":   true,
		"timestamp": time.Now().Unix(),
		"message":   "Well done, mate!",
	}
}
