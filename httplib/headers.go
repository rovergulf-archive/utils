package httplib

import (
	"net/http"
)

// GetDeviceTypeFromRequestHeaders
func GetDeviceTypeFromRequestHeaders(r *http.Request) string {
	return r.Header.Get("User-Agent")
}
