package httplib

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strings"
)

type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type HTTPInterceptor struct {
	Router Router
	Tracer opentracing.Tracer
}

func (i *HTTPInterceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Set request headers for AJAX requests
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Cookie, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if i.Tracer != nil {
		operation := strings.TrimPrefix(r.URL.Path, "/")
		spanName := fmt.Sprintf("%s:%s", r.Method, operation)

		query := r.URL.RawQuery
		if len(query) > 0 {
			spanName += "?" + query
		}

		span := i.Tracer.StartSpan(spanName)
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(r.Context(), span)
		r = r.WithContext(ctx)
	}

	i.Router.ServeHTTP(w, r)
}
