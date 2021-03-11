package httplib

import (
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/rovergulf/utils/tracing"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Interceptor
type Interceptor struct {
	Router *mux.Router
	Tracer opentracing.Tracer
	Logger *zap.SugaredLogger
}

func NewInterceptor(lg *zap.SugaredLogger, j *tracing.Jaeger) Interceptor {
	i := Interceptor{
		Logger: lg,
	}

	if j != nil {
		i.Tracer = j.Tracer
	}

	return i
}

// ServeHTTP
func (i *Interceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Set request headers for AJAX requests
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Cookie, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
	}

	// handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if i.Tracer != nil {
		span := i.Tracer.StartSpan(strings.TrimPrefix(r.URL.Path, "/"))
		span.SetTag("host", r.Host)
		span.SetTag("method", r.Method)
		span.SetTag("path", r.URL.Path)
		span.SetTag("query", r.URL.RawQuery)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	i.Logger.Infow("Handling request", "method", r.Method, "path", r.URL.Path, "query", r.URL.RawQuery)

	i.Router.ServeHTTP(w, r.WithContext(ctx))
}
