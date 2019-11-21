package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// Channel to handle graceful shutdown of GRPC Server
	ch := make(chan os.Signal, 1)

	// Write in Channel in case of OS request to shut process
	signal.Notify(ch, os.Interrupt)

	promHTTPSrv := http.Server{Addr: ":2112"}

	// Start Prometheus HTTP handler
	go func() {
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
	}()

	// Block main function from exiting until ch receives value
	<-ch
	logging.PeppaMonLog("warning", "Shutting down Peppamon server...")

	// Stop Prom HTTP server
	ctxPromHTTP, ctxCancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer ctxCancel()

	errPromHTTPShut := promHTTPSrv.Shutdown(ctxPromHTTP)

	if errPromHTTPShut != nil {
		logging.PeppaMonLog(
			"warning",
			"Error while shutting down Prometheus HTTP Server %v", errPromHTTPShut)
	}
}
