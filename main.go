package traefik_bearer_token_plugin

import (
	"context"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
}

func CreateConfig() *Config {
	return &Config{}
}

type BearerTokenMiddleware struct {
	next   http.Handler
	metric *prometheus.CounterVec
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
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

func init() {
	prometheusHandler := promhttp.Handler()
	http.Handle("/metrics", prometheusHandler)
}
