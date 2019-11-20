package main

import (
	"net/http"

	"github.com/lucabrasi83/peppamon_versa/initializer"
	"github.com/lucabrasi83/peppamon_versa/logging"
	"github.com/lucabrasi83/peppamon_versa/versa_collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	collector = versa_collector.NewVersaAnalyticsExporter()
)

func init() {
	prometheus.MustRegister(collector)
}

func main() {

	initializer.Initialize()

	promHTTPSrv := http.Server{Addr: ":2112"}

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, errWelcomePage := w.Write([]byte(`<html>
             <head><title>Peppamon Versa Analytics Exporter</title></head>
             <body>
             <h1>Peppamon Versa Analytics Telemetry Exporter</h1>
             <p><a href="/metrics">Metrics</a></p>
             </body>
             </html>`))

		if errWelcomePage != nil {
			logging.PeppaMonLog(
				"error",
				"Failed to render Welcome page %v", errWelcomePage)
		}
	})

	logging.PeppaMonLog("info", "Starting Prometheus metrics web handler for Versa on port TCP 2112...")

	if err := promHTTPSrv.ListenAndServe(); err != http.ErrServerClosed {
		logging.PeppaMonLog(
			"fatal",
			"Failed to start Prometheus HTTP metrics handler %v", err)
	}
}
