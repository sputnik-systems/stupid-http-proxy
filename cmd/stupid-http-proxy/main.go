package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	targetFlag := pflag.String("target", "http://127.0.0.1:8080", "Target http url")
	pflag.Parse()

	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	target, err := url.Parse(*targetFlag)
	if err != nil {
		panic(err)
	}
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.Host = r.In.Host
		},
	}
	http.ListenAndServe(":8080", logRequestMiddleware(proxy))
}

func logRequestMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method
		headers := r.Header
		logger.Info("access log",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Reflect("headers", headers),
		)
	}

	return http.HandlerFunc(fn)
}
