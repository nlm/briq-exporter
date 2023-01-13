package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/nlm/briq-cli/briq"
	"github.com/nlm/briq-exporter/prombriq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	flagListen = flag.String("listen", ":9000", "interface to listen on")
)

func main() {
	flag.Parse()
	logger := log.New()

	// Setup briq client
	secretKey := os.Getenv("BRIQ_SECRET_KEY")
	if secretKey == "" {
		logger.Fatal("you must define BRIQ_SECRET_KEY in the environment")
	}
	client, err := briq.NewClient(os.Getenv("BRIQ_SECRET_KEY"))
	if err != nil {
		logger.WithError(err).Fatal("unable to initialize briq client")
	}

	// Setup collector
	collector := prombriq.NewCollector(client, prombriq.WithLogger(logger))
	prometheus.MustRegister(collector)

	// Listen and Serve
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*flagListen, nil)
}
