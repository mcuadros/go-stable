FROM alpine:latest
MAINTAINER Máximo Cuadros <mcuadros@gmail.com>

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ADD cli/stable/stable /usr/local/bin/
CMD ["stable"]