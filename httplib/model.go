package httplib

import (
	"fmt"
	"time"
)

type ApiError struct {
	HttpStatus int         `json:"http_status,omitempty"`
	ErrorCode  int         `json:"code"`
	Message    interface{} `json:"message"`
	Timestamp  interface{} `json:"timestamp"`
}

func (e ApiError) Error() string {
	return e.String()
}

func (e ApiError) String() string {
	msg := fmt.Sprintf("%s", e.Message)

	if e.HttpStatus > 0 {
		msg = fmt.Sprintf("[HTTP_CODE:%d] %s", e.HttpStatus, msg)
	}

	if e.ErrorCode > 0 {
		msg = fmt.Sprintf("[INTERNAL_CODE:%d] %s", e.ErrorCode, msg)
	}

	if e.Timestamp != nil {
		msg = fmt.Sprintf("[%s] %s", time.Unix(e.Timestamp.(int64), 0), msg)
	}

	return fmt.Sprintf("%s", msg)
}

func NewApiError(code int, msg string) *ApiError {
	return &ApiError{
		ErrorCode: code,
		Message:   msg,
	}
}

type ResultAdditionalFields map[string]interface{}

type CreatedObjectId struct {
	Id   interface{} `json:"id,omitempty" yaml:"id,omitempty"`
	UUID interface{} `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ListResult struct {
	Results interface{} `json:"results,omitempty" yaml:"results,omitempty"`
	Count   int32       `json:"count" yaml:"count"`
	HasPrev bool        `json:"has_prev" yaml:"has_prev"`
	HasNext bool        `json:"has_next" yaml:"has_next"`
}

type Response struct {
	Success   bool        `json:"success" yaml:"success"`
	Timestamp interface{} `json:"timestamp,omitempty" yaml:"timestamp,omitempty" `
	Message   interface{} `json:"message,omitempty" yaml:"message,omitempty"`
	Data      interface{} `json:"data,omitempty" yaml:"data,omitempty"`
	CreatedObjectId
	RequestStatus *ApiError `json:"request_status" yaml:"request_status,omitempty"`
}

func SuccessfulResult() Response {
	return Response{
		Success:   true,
		Timestamp: time.Now().Unix(),
		Message:   "Well done, mate!",
	}
}

func SuccessfulResultMap() map[string]interface{} {
	return map[string]interface{}{
		"success":   true,
		"timestamp": time.Now().Unix(),
		"message":   "Well done, mate!",
	}
}
