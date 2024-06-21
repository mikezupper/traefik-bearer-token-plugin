// Package traefikbearertokenplugin plugin.
package traefikbearertokenplugin

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	//Headers map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		//	Headers: make(map[string]string),
	}
}

// BearerTokenMiddleware plugin.
type BearerTokenMiddleware struct {
	next   http.Handler
	metric *prometheus.CounterVec
}

// New created BearerTokenMiddleware plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	prometheusHandler := promhttp.Handler()
	http.Handle("/metrics", prometheusHandler)

	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "traefik_bearer_token_requests_total",
			Help: "Total number of requests by endpoint and token.",
		},
		[]string{"endpoint", "token"},
	)
	prometheus.MustRegister(metric)

	return &BearerTokenMiddleware{
		next:   next,
		metric: metric,
	}, nil
}

func (bt *BearerTokenMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	token := ""
	authHeader := req.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}

	endpoint := req.URL.Path
	bt.metric.WithLabelValues(endpoint, token).Inc()

	bt.next.ServeHTTP(rw, req)
}
