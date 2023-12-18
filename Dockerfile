FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /usr/local/src/rtorrent_exporter
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o /usr/local/bin/rtorrent_exporter

FROM cgr.dev/chainguard/static:latest
COPY --from=builder /usr/local/bin/rtorrent_exporter /usr/local/bin/rtorrent_exporter
EXPOSE 9135
ENTRYPOINT [ "/usr/local/bin/rtorrent_exporter" ]
CMD [ "" ]
