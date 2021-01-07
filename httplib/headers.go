package httplib

import (
	"net/http"
)

func GetDeviceTypeFromRequestHeaders(r *http.Request) string {
	return r.Header.Get("User-Agent")
}
