# Prometheus rTorrent Exporter

`rtorrent_exporter` exposes metrics from a rtorrent instance.

## Installation

Binaries can be downloaded from the [Github releases](https://github.com/thde/rtorrent_exporter/releases) page and need no
special installation.

A Docker container (`ghcr.io/thde/rtorrent_exporter:latest`) is also available through [Github's registry](https://github.com/thde/rtorrent_exporter/pkgs/container/rtorrent_exporter).

## Usage

```
usage: rtorrent_exporter [<flags>]

Flags:
  -h, --help                 Show context-sensitive help (also try --help-long and --help-man).
      --web.config.file=""   [EXPERIMENTAL] Path to configuration file that can enable TLS or authentication.
      --web.listen-address=":9135"
                             Address to listen on for web interface and telemetry.
      --web.telemetry-path="/metrics"
                             Path under which to expose metrics.
      --rtorrent.scrape-uri="http://localhost/RPC2"
                             URI on which to scrape rTorrent. Use http://user:pass@host.com to supply basic auth credentials.
      --rtorrent.ssl-verify  Flag that enables SSL certificate verification for the scrape URI.
      --rtorrent.timeout=5s  Timeout for trying to get stats from rtorrent.
      --log.level=info       Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt    Output format of log messages. One of: [logfmt, json]
      --version              Show application version.
```
