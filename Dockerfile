FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/

ADD build/gonotify /
ADD config/config.yml /config

EXPOSE 3333

CMD ["/gonotify"]
