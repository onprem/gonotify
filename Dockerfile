FROM alpine:latest
RUN apk --no-cache add ca-certificates

ADD build/gonotify /
ADD config/config.yml /config/config.yml

EXPOSE 3333

RUN mkdir -p /database

CMD ["/gonotify"]
