package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	rxPktMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rxPkts",
			Help: "Received pakets",
		},
		[]string{"port", "location"},
	)
	txPktMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "txPkts",
			Help: "Transmitted pakets",
		},
		[]string{"port", "location"},
	)
	crcPktMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "crcPkts",
			Help: "Packets dropped by switch",
		},
		[]string{"port", "location"},
	)
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	remoteIp, ok := os.LookupEnv("REMOTE_IP")
	if !ok {
		panic("REMOTE_IP not set")
	}
	passwd, ok := os.LookupEnv("PASSWORD")
	if !ok {
		panic("PASSWORD not set")
	}
	location, ok := os.LookupEnv("LOCATION")
	if !ok {
		panic("LOCATION not set")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("PORT not set")
	}
	var sc scrapeClient
	client, err := NewScrapeClient()
	if err != nil {
		panic(err)
	}
	sc.client = client
	sc.remote = remoteIp
	sc.location = location
	sc.password = passwd
	sc.logger = logger

	logger.Info("build on", "go_version", runtime.Version())

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll),
			collectors.WithoutGoCollectorRuntimeMetrics(collectors.MetricsGC.Matcher),
		),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		rxPktMetric,
		txPktMetric,
		crcPktMetric,
	)

	go func() {
		for {
			err := sc.crunchMetrics()
			if err != nil {
				logger.Error(err.Error())
			}
			time.Sleep(10 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
