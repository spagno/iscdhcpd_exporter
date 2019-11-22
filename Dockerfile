FROM golang as build

RUN "env"

ADD . /go/src/github.com/homecentr/iscdhcpd_exporter

RUN cd /go/src/github.com/homecentr/iscdhcpd_exporter \
    && make build