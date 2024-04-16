package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mrobinsn/go-rtorrent/rtorrent"
	"github.com/prometheus/client_golang/prometheus"
	collectorsversion "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	logflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"

	"github.com/thde/rtorrent_exporter/exporter"
)

const (
	namespace = "rtorrent"
	name      = namespace + "_exporter"
)

var (
	cli               = kingpin.New(name, "Prometheus exporter for rTorrent.").DefaultEnvars()
	webConfig         = webflag.AddFlags(cli, ":9135")
	metricsPath       = cli.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	rtorrentScrapeURI = cli.Flag("rtorrent.scrape-uri", "URI on which to scrape rTorrent. Use http://user:pass@host.com to supply basic auth credentials.").Default("http://localhost/RPC2").String()
	rtorrentSSLVerify = cli.Flag("rtorrent.ssl-verify", "Flag that enables SSL certificate verification for the scrape URI.").Default("true").Bool()
	rtorrentTimeout   = cli.Flag("rtorrent.timeout", "Timeout for trying to get stats from rtorrent.").Default("5s").Duration()
)

func run(logger log.Logger) error {
	level.Info(logger).Log("msg", "Starting rtorrent_exporter", "version", version.Info())
	level.Info(logger).Log("build_context", version.BuildContext())

	prometheus.MustRegister(collectorsversion.NewCollector(name))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
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

	srv := &http.Server{
		ReadTimeout:  *rtorrentTimeout,
		WriteTimeout: *rtorrentTimeout,
	}
	return web.ListenAndServe(srv, webConfig, logger)
}

func main() {
	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)

	logflag.AddFlags(cli, promlogConfig)
	cli.Version(version.Print(name))
	cli.HelpFlag.Short('h')
	if _, err := cli.Parse(os.Args[1:]); err != nil {
		level.Error(logger).Log("msg", "Error parsing CLI flags", "err", err)
		os.Exit(1)
	}

	if err := run(logger); err != nil {
		level.Error(logger).Log("msg", "Error starting exporter", "err", err)
		os.Exit(1)
	}
}
