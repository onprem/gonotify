FROM alpine:latest
RUN apk --no-cache add ca-certificates

ADD build/gonotify /
ADD config/config.yml /config

EXPOSE 3333

CMD ["/gonotify"]
