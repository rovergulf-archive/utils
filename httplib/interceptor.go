package httplib

import (
	"context"
	"crypto/tls"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/rovergulf/utils/ipaddr"
	"github.com/rovergulf/utils/tracing"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Interceptor
type Interceptor struct {
	Router  *mux.Router
	Tracer  opentracing.Tracer
	Logger  *zap.SugaredLogger
	tlsConf *tls.Config
}

const (
	headersSep = ", "
)

var allowedHeaders = []string{
	ipaddr.XForwardedFor,
	ipaddr.CFConnectingIp,
	ipaddr.CFRealIp,
	"Accept",
	"Content-Type",
	"Content-Length",
	"Cookie",
	"Accept-Encoding",
	"Authorization",
	"X-CSRF-Token",
	"X-Requested-With",
}

var allowedMethods = []string{
	"OPTIONS",
	"GET",
	"PUT",
	"PATCH",
	"POST",
	"DELETE",
}

func NewInterceptor(lg *zap.SugaredLogger, j *tracing.Jaeger, tlsConf *tls.Config) *Interceptor {
	i := &Interceptor{
		Logger:  lg,
		tlsConf: tlsConf,
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
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, headersSep))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, headersSep))
	}

	// handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx = context.WithValue(ctx, "host", r.Host)
	ctx = context.WithValue(ctx, "path", r.URL.Path)
	ctx = context.WithValue(ctx, "remote_addr", r.RemoteAddr)
	ctx = context.WithValue(ctx, "x_forwarded_for", r.Header.Get(ipaddr.XForwardedFor))
	ctx = context.WithValue(ctx, "cf_connecting_ip", r.Header.Get(ipaddr.CFConnectingIp))

	if i.Tracer != nil {
		span := i.Tracer.StartSpan(strings.TrimPrefix(r.URL.Path, "/"))
		span.SetTag("host", r.Host)
		span.SetTag("method", r.Method)
		span.SetTag("path", r.URL.Path)
		span.SetTag("query", r.URL.RawQuery)
		span.SetTag("remote_addr", r.RemoteAddr)
		span.SetTag("x_forwarded_for", r.Header.Get(ipaddr.XForwardedFor))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	i.Logger.Infow("Handling request", "method", r.Method, "path", r.URL.Path, "query", r.URL.RawQuery)

	i.Router.ServeHTTP(w, r.WithContext(ctx))
}
