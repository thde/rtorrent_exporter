package main

import (
	"net/http"
	"os"

	"github.com/go-kit/log/level"
	"github.com/mrobinsn/go-rtorrent/rtorrent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/thde/rtorrent_exporter/exporter"
	"gopkg.in/alecthomas/kingpin.v2"

	logflag "github.com/prometheus/common/promlog/flag"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const namespace = "rtorrent"

var (
	webConfig         = webflag.AddFlags(kingpin.CommandLine)
	listenAddress     = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9135").String()
	metricsPath       = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	rtorrentScrapeURI = kingpin.Flag("rtorrent.scrape-uri", "URI on which to scrape rTorrent. Use http://user:pass@host.com to supply basic auth credentials.").Default("http://localhost/RPC2").String()
	rtorrentSSLVerify = kingpin.Flag("rtorrent.ssl-verify", "Flag that enables SSL certificate verification for the scrape URI.").Default("true").Bool()
	rtorrentTimeout   = kingpin.Flag("rtorrent.timeout", "Timeout for trying to get stats from rtorrent.").Default("5s").Duration()
)

func main() {
	promlogConfig := &promlog.Config{}
	logflag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("rtorrent_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting rtorrent_exporter", "version", version.Info())
	level.Info(logger).Log("build_context", version.BuildContext())

	prometheus.MustRegister(version.NewCollector("rtorrent_exporter"))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!doctype html>
		<html lang=en>
		<head>
		<meta charset=utf-8>
		<title>rtorrent exporter</title>
		</head>
		<body>
		<pre><a href='` + *metricsPath + `'>rtorrent exporter</a></pre>`))
	})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "client", "requests_total"),
			Help: "A counter for requests from the rtorrent client.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(counter)

	conn := rtorrent.New(*rtorrentScrapeURI, *rtorrentSSLVerify).WithHTTPClient(&http.Client{
		Timeout:   *rtorrentTimeout,
		Transport: promhttp.InstrumentRoundTripperCounter(counter, http.DefaultTransport),
	})
	e := exporter.Exporter{
		Namespace: namespace,
		Client:    *conn,
		Logger:    logger,
	}
	prometheus.MustRegister(&e)

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
