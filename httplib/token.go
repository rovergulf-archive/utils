package httplib

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	PrefixBearer = "Bearer "
)

func ExtractTokenFromRequest(r *http.Request) (string, error) {
	cookie, _ := r.Cookie(CookieName)
	token := cookie.Value
	if len(token) > 0 {
		return token, nil
	}

	token = GetAuthorizationTokenFromRequestHeader(r, PrefixBearer)
	if len(token) > 0 {
		return token, nil
	}

	token = r.FormValue(CookieName)
	if len(token) == 0 || token == " " {
		return "", fmt.Errorf("not a token cookie, nor auth header, nor query parameter specified")
	}

	return token, nil
}

func GetAuthorizationTokenFromRequestHeader(r *http.Request, prefix string) string {
	if prefix == "" {
		prefix = PrefixBearer // ?? is it good enough for default?
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, prefix) {
		return ""
	}
	return authHeader[len(prefix):]
}
