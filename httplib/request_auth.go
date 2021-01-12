package httplib

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	PrefixBearer = "Bearer "
	ContextToken = "access_token"
)

func ExtractTokenFromRequest(r *http.Request) (string, error) {
	token, _ := GetRequestCookieStringValue(r, CookieName)
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

func GetTokenIdFromRequestCookie(r *http.Request) string {
	token, _ := GetRequestCookieStringValue(r, CookieName)

	return token
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

func HTTPTokenHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenId, _ := ExtractTokenFromRequest(r)
		ctx = context.WithValue(ctx, ContextToken, tokenId)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
