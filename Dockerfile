FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY rtorrent_exporter /usr/local/bin/rtorrent_exporter
USER guest
EXPOSE 9135
ENTRYPOINT [ "/usr/local/bin/rtorrent_exporter" ]
CMD [ "" ]
