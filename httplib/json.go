package httplib

import (
	"encoding/json"
	"net/http"
)

func makeError(code int, err error) *ApiError {
	return NewApiError(code, err.Error())
}

// Sends error http response
func ErrorResponseJSON(w http.ResponseWriter, httpCode int, internalCode int, err error) {
	writeJSON(w, httpCode, makeError(internalCode, err))
}

// Sends OK JSON response
func ResponseJSON(w http.ResponseWriter, v interface{}) {
	writeJSON(w, http.StatusOK, v)
}

// Write JSON
func writeJSON(w http.ResponseWriter, httpCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response, err := json.Marshal(v)
	if err != nil {
		w.Write([]byte("Cannot marshal response: " + err.Error()))
		return
	}

	w.WriteHeader(httpCode)
	w.Write(response)
}
