package httplib

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	contextRequestPath   = "request_path"
	contextRequestMethod = "request_method"
)

// Router
type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// HTTPInterceptor
type HTTPInterceptor struct {
	Router Router
	Tracer opentracing.Tracer
	Logger *zap.SugaredLogger
}

// ServeHTTP
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

	ctx = context.WithValue(ctx, contextRequestPath, operation)
	ctx = context.WithValue(ctx, contextRequestMethod, r.Method)
	r = r.WithContext(ctx)

	i.Logger.Infof("Handling request [%s][%s]", r.Method, operation)

	i.Router.ServeHTTP(w, r)
}

// ResponseJSON
func (i *HTTPInterceptor) ResponseJSON(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	i.Logger.Debugw("Successful HTTP Request",
		contextRequestMethod, ctx.Value(contextRequestMethod),
		contextRequestPath, ctx.Value(contextRequestPath))
	ResponseJSON(w, payload)
}

// ErrorResponseJSON
func (i *HTTPInterceptor) ErrorResponseJSON(ctx context.Context, w http.ResponseWriter, statusCode, internalCode int, err error) {
	i.Logger.Debugw("Failed HTTP Request",
		contextRequestMethod, ctx.Value(contextRequestMethod),
		contextRequestPath, ctx.Value(contextRequestPath))
	ErrorResponseJSON(w, statusCode, internalCode, err)
}
