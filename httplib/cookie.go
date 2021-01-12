package httplib

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var CookieName = "session"

func SetCookieName(name string) {
	CookieName = name
}

// GetRequestCookieStringValue
func GetRequestCookieStringValue(r *http.Request, cookieName string) (string, error) {
	prefix := cookieName + "="
	retrievedCookie, _ := r.Cookie(cookieName)
	if retrievedCookie != nil && retrievedCookie.Value != "" {
		return retrievedCookie.Value, nil
	}

	cookieVals, present := r.Header["Cookie"]
	if present && len(cookieVals) > 0 {
		cookieData := cookieVals[0]
		cookieElems := strings.Split(cookieData, "; ")
		for _, elem := range cookieElems {
			if strings.HasPrefix(elem, prefix) {
				elemVal := strings.TrimPrefix(elem, prefix)
				if len(elemVal) > 0 {
					return elemVal, nil
				}
			}
		}
	}

	return "", fmt.Errorf("%s", "Empty cookie or cannot be read")
}

// SetHttpCookieValue
func SetHttpCookieValue(w http.ResponseWriter, domain, cookieName, value string, expireTime time.Time) {
	cookie := http.Cookie{Name: cookieName, Value: value, Expires: expireTime, Domain: domain, HttpOnly: true}
	http.SetCookie(w, &cookie)
}
