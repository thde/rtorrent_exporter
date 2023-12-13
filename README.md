# Prometheus rTorrent Exporter

`rtorrent_exporter` exposes metrics from a rtorrent instance.

## Installation

Binaries can be downloaded from the [Github releases](https://github.com/thde/rtorrent_exporter/releases) page and need no
special installation.

A Docker container (`ghcr.io/thde/rtorrent_exporter:latest`) is also available through [Github's registry](https://github.com/thde/rtorrent_exporter/pkgs/container/rtorrent_exporter).

## Usage

```
usage: rtorrent_exporter [<flags>]

Prometheus exporter for rTorrent.

Flags:
  -h, --[no-]help            Show context-sensitive help (also try --help-long and --help-man). ($RTORRENT_EXPORTER_HELP)
      --web.listen-address=:9135 ...  
                             Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.
                             ($RTORRENT_EXPORTER_WEB_LISTEN_ADDRESS)
      --web.config.file=""   Path to configuration file that can enable TLS or authentication. See:
                             https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
                             ($RTORRENT_EXPORTER_WEB_CONFIG_FILE)
      --web.telemetry-path="/metrics"  
                             Path under which to expose metrics. ($RTORRENT_EXPORTER_WEB_TELEMETRY_PATH)
      --rtorrent.scrape-uri="http://localhost/RPC2"  
                             URI on which to scrape rTorrent. Use http://user:pass@host.com to supply basic auth
                             credentials. ($RTORRENT_EXPORTER_RTORRENT_SCRAPE_URI)
      --[no-]rtorrent.ssl-verify  
                             Flag that enables SSL certificate verification for the scrape URI.
                             ($RTORRENT_EXPORTER_RTORRENT_SSL_VERIFY)
      --rtorrent.timeout=5s  Timeout for trying to get stats from rtorrent. ($RTORRENT_EXPORTER_RTORRENT_TIMEOUT)
      --log.level=info       Only log messages with the given severity or above. One of: [debug, info, warn, error]
                             ($RTORRENT_EXPORTER_LOG_LEVEL)
      --log.format=logfmt    Output format of log messages. One of: [logfmt, json] ($RTORRENT_EXPORTER_LOG_FORMAT)
      --[no-]version         Show application version. ($RTORRENT_EXPORTER_VERSION)
```
