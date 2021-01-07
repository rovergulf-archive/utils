package httplib

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

var SessionName = "session"

func SetSessionName(name string) {
	SessionName = name
}

const (
	PrefixBearer = "Bearer "
)

func ExtractSessionFromRequest(r *http.Request) (string, error) {
	session, _ := GetRequestCookieStringValue(r, SessionName)
	if len(session) > 0 {
		return session, nil
	}

	session = GetAuthorizationTokenFromRequestHeader(r, PrefixBearer)
	if len(session) > 0 {
		return session, nil
	}

	session = r.FormValue(SessionName)
	if len(session) == 0 || session == " " {
		return "", fmt.Errorf("not a session cookie, nor auth header, nor query parameter specified")
	}

	return session, nil
}

func GetSessionIdFromRequestCookie(r *http.Request) string {
	session, _ := GetRequestCookieStringValue(r, SessionName)

	return session
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

func HTTPSessionHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionId, _ := ExtractSessionFromRequest(r)
		ctx = context.WithValue(ctx, "sessionId", sessionId)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
