FROM golang:1.5.2-alpine
MAINTAINER Máximo Cuadros <mcuadros@gmail.com>

RUN apk -U add git && rm -rf /var/cache/apk/*

ADD . ${GOPATH}/src/github.com/mcuadros/gop.kg
WORKDIR ${GOPATH}/src/github.com/mcuadros/gop.kg
RUN go get -v ./...
RUN go install -v ./...

VOLUME /certificates
EXPOSE 443

CMD ["gop.kg", "server", "--addr=:443", "--cert=/certificates/gop.kg.cert.pem", "--key=/certificates/gop.kg.key.pem"]
