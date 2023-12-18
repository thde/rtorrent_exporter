package exporter

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mrobinsn/go-rtorrent/rtorrent"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "rtorrent"

var (
	rtorrentInfo       = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "info"), "rtorrent info.", []string{"name", "ip"}, nil)
	rtorrentUp         = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "up"), "Was the last scrape of rTorrent successful.", nil, nil)
	rtorrentDownloaded = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "downloaded_bytes"), "Total downloaded bytes.", nil, nil)
	rtorrentUploaded   = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "uploaded_bytes"), "Total uploaded bytes.", nil, nil)
	rtorrentTorrents   = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "torrents_total"), "Torrent count by view.", []string{"view"}, nil)
)

// Exporter returns a prometheus.Collector that gathers rTorrent metrics.
type Exporter struct {
	Namespace string
	Client    rtorrent.RTorrent
	Logger    log.Logger
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- rtorrentInfo
	ch <- rtorrentUp
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	up := e.scrape(ch)

	ch <- prometheus.MustNewConstMetric(rtorrentUp, prometheus.GaugeValue, up)
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) (up float64) {
	name, err := e.Client.Name()
	if err != nil {
		level.Error(e.Logger).Log("msg", "Can't scrape rTorrent", "err", err)
		return 1
	}

	ip, err := e.Client.IP()
	if err != nil {
		level.Error(e.Logger).Log("msg", "Can't scrape rTorrent", "err", err)
		return 1
	}

	ch <- prometheus.MustNewConstMetric(rtorrentInfo, prometheus.GaugeValue, 1, name, ip)

	downloaded, err := e.Client.DownTotal()
	if err != nil {
		level.Error(e.Logger).Log("msg", "Can't scrape rTorrent", "err", err)
		return 1

	}
	ch <- prometheus.MustNewConstMetric(rtorrentDownloaded, prometheus.CounterValue, float64(downloaded))

	uploaded, err := e.Client.UpTotal()
	if err != nil {
		level.Error(e.Logger).Log("msg", "Can't scrape rTorrent", "err", err)
		return 1
	}
	ch <- prometheus.MustNewConstMetric(rtorrentUploaded, prometheus.CounterValue, float64(uploaded))

	for name, view := range map[string]rtorrent.View{
		"main":    rtorrent.ViewMain,
		"seeding": rtorrent.ViewSeeding,
		"hashing": rtorrent.ViewHashing,
		"started": rtorrent.ViewStarted,
		"stopped": rtorrent.ViewStopped,
	} {
		torrents, err := e.Client.GetTorrents(view)
		if err != nil {
			level.Error(e.Logger).Log("msg", "Can't scrape rTorrent", "err", err)
			return 1
		}

		ch <- prometheus.MustNewConstMetric(rtorrentTorrents, prometheus.CounterValue, float64(len(torrents)), name)
	}

	return 0
}
