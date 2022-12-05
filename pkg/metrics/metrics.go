package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

// Metrics configures the metrics handler
type Metrics struct {
	Enabled bool
	Port    string
}

// Handle ...
// HTTP handler for metrics
func (m *Metrics) Handle() {
	if m.Enabled == false {
		return
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Metrics listening on %v\n", m.Port)
	http.ListenAndServe(m.Port, nil)
}
