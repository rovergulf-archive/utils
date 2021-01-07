package httplib

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type HTTPInterceptor struct {
	Router Router
	Tracer opentracing.Tracer
	Logger *zap.SugaredLogger
}

func (i *HTTPInterceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	// define span name
	operation := strings.TrimPrefix(r.URL.Path, "/")
	spanName := fmt.Sprintf("%s:%s", r.Method, operation)

	query := r.URL.RawQuery
	if len(query) > 0 {
		spanName += "?" + query
	}

	if i.Tracer != nil {
		span := i.Tracer.StartSpan(spanName)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	ctx = context.WithValue(ctx, "request_path", operation)
	r = r.WithContext(ctx)

	i.Logger.Infof("Handling request [%s]", spanName)

	i.Router.ServeHTTP(w, r)
}
