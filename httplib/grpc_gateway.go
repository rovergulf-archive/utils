package httplib

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type GatewayInterceptor struct {
	Mux    *runtime.ServeMux
	Tracer opentracing.Tracer
}

func (i *GatewayInterceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Set request headers for AJAX requests
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Cookie, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//if i.Tracer != nil {
	//	operation := strings.TrimPrefix(r.URL.Path, "/")
	//	span := i.Tracer.StartSpan(fmt.Sprintf("%s:%s?%s", r.Method, operation, r.URL.RawQuery))
	//	defer span.Finish()
	//
	//	ctx := opentracing.ContextWithSpan(r.Context(), span)
	//	r = r.WithContext(ctx)
	//}

	i.Mux.ServeHTTP(w, r)
}
