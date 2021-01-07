package httplib

import (
	"github.com/rovergulf/utils"
	"net/http"
	"strings"
)

func GetRequestIPAddress(r *http.Request) string {
	ipAddress := r.RemoteAddr
	fwdAddress := r.Header.Get("X-Forwarded-For")
	if fwdAddress != "" {
		// Got X-Forwarded-For
		ipAddress = fwdAddress // if it's a single IP

		// if array – grab the first
		ips := strings.Split(fwdAddress, ", ")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	}
	return ipAddress
}

func GetRequestAuthorFootprint(r *http.Request) string {
	return utils.GenerateHashFromString(GetRequestIPAddress(r) + ":" + r.Header.Get("User-Agent"))
}
