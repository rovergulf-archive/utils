package ipaddr

import (
	"github.com/rovergulf/utils"
	"net/http"
	"strings"
)

const (
	XForwardedFor  = "X-Forwarded-For"
	CFConnectingIp = "CF-Connecting-IP"
	CFRealIp       = "CF-Real-IP"
)

// HttpForwardedFor returns X-Forwarded-For Header value if not empty
func HttpForwardedFor(r *http.Request) string {
	ipAddress := r.RemoteAddr
	fwdAddress := r.Header.Get(XForwardedFor)
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

// cloudflare header names reference:
// https://support.cloudflare.com/hc/en-us/articles/200170786-Restoring-original-visitor-IPs

// HttpCloudflareConnectingIP returns CF-Connecting-IP Header value if not empty
func HttpCloudflareConnectingIP(r *http.Request) string {
	ipAddress := r.RemoteAddr
	fwdAddress := r.Header.Get(CFConnectingIp)
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

// HttpCloudflareRealIP returns CF-Real-IP Header value if not empty
func HttpCloudflareRealIP(r *http.Request) string {
	ipAddress := r.RemoteAddr
	fwdAddress := r.Header.Get(CFRealIp)
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

func httpRequestUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func HttpRequestFootprint(r *http.Request) string {
	return utils.GenerateHashFromString(GetRequestIPAddress(r) + ":" + httpRequestUserAgent(r))
}

func GetRequestIPAddress(r *http.Request) string {
	if cloudflareRealIp := HttpCloudflareRealIP(r); len(cloudflareRealIp) > 7 {
		return cloudflareRealIp
	}

	if cloudflareConnectingIp := HttpCloudflareConnectingIP(r); len(cloudflareConnectingIp) > 7 {
		return cloudflareConnectingIp
	}

	return HttpForwardedFor(r)
}
